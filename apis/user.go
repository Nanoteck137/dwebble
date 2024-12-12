package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type UpdateUserSettingsBody struct {
	DisplayName   *string `json:"displayName,omitempty"`
	QuickPlaylist *string `json:"quickPlaylist,omitempty"`
}

func (b *UpdateUserSettingsBody) Transform() {
	b.DisplayName = transform.StringPtr(b.DisplayName)
	b.QuickPlaylist = transform.StringPtr(b.QuickPlaylist)
}

func (b UpdateUserSettingsBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.DisplayName,
			validate.Required.When(b.DisplayName != nil),
		),
		validate.Field(&b.QuickPlaylist), // validate.Required.When(b.QuickPlaylist != nil),

	)
}

type TrackId struct {
	TrackId string `json:"trackId"`
}

func (b *TrackId) Transform() {
	// b.Tracks = *transform.DiscardEmptyStringEntries()
}

type GetUserQuickPlaylistItemIds struct {
	TrackIds []string `json:"trackIds"`
}

type CreateApiToken struct {
	Token string `json:"token"`
}

type CreateApiTokenBody struct {
	Name string `json:"name"`
}

func (b *CreateApiTokenBody) Transform() {
	b.Name = transform.String(b.Name)
}

func (b CreateApiTokenBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

type ApiToken struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetAllApiTokens struct {
	Tokens []ApiToken `json:"tokens"` 
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "UpdateUserSettings",
			Method:   http.MethodPatch,
			Path:     "/user/settings",
			BodyType: UpdateUserSettingsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[UpdateUserSettingsBody](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				pretty.Println(user)
				pretty.Println(body)

				settings := user.ToUserSettings()

				if body.DisplayName != nil {
					settings.DisplayName = sql.NullString{
						String: *body.DisplayName,
						Valid:  true,
					}
				}

				if body.QuickPlaylist != nil {
					id := *body.QuickPlaylist

					if id != "" {
						_, err := app.DB().GetPlaylistById(context.TODO(), id)
						if err != nil {
							// TODO(patrik): Handle error
							return nil, err
						}
					}

					settings.QuickPlaylist = sql.NullString{
						String: id,
						Valid:  id != "",
					}
				}

				err = app.DB().UpdateUserSettings(context.TODO(), settings)
				if err != nil {
					// TODO(patrik): Handle error
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddToUserQuickPlaylist",
			Method:   http.MethodPost,
			Path:     "/user/quickplaylist",
			DataType: nil,
			BodyType: TrackId{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[TrackId](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				if user.QuickPlaylist.Valid {
					err := app.DB().AddItemToPlaylist(ctx, user.QuickPlaylist.String, body.TrackId)
					if err != nil {
						// TODO(patrik): Handle error
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "RemoveItemFromUserQuickPlaylist",
			Method:   http.MethodDelete,
			Path:     "/user/quickplaylist",
			DataType: nil,
			BodyType: TrackId{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[TrackId](c)
				if err != nil {
					return nil, err
				}

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				if user.QuickPlaylist.Valid {
					err := app.DB().RemovePlaylistItem(ctx, user.QuickPlaylist.String, body.TrackId)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetUserQuickPlaylistItemIds",
			Method:   http.MethodGet,
			Path:     "/user/quickplaylist",
			DataType: GetUserQuickPlaylistItemIds{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				if user.QuickPlaylist.Valid {
					items, err := app.DB().GetPlaylistItems(ctx, user.QuickPlaylist.String)
					if err != nil {
						return nil, err
					}

					res := GetUserQuickPlaylistItemIds{
						TrackIds: make([]string, len(items)),
					}

					for i, item := range items {
						res.TrackIds[i] = item.TrackId
					}

					return res, nil
				}

				// TODO(patrik): Better error
				return nil, errors.New("No Quick Playlist set")
			},
		},

		pyrin.ApiHandler{
			Name:     "CreateApiToken",
			Method:   http.MethodPost,
			Path:     "/user/apitoken",
			DataType: CreateApiToken{},
			BodyType: CreateApiTokenBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreateApiTokenBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				token, err := app.DB().CreateApiToken(ctx, database.CreateApiTokenParams{
					UserId: user.Id,
					Name:   body.Name,
				})
				if err != nil {
					return nil, err
				}

				return CreateApiToken{
					Token: token.Id,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetAllApiTokens",
			Method:   http.MethodGet,
			Path:     "/user/apitoken",
			DataType: GetAllApiTokens{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				tokens, err := app.DB().GetAllApiTokensForUser(ctx, user.Id)
				if err != nil {
					return nil, err
				}

				res := GetAllApiTokens{
					Tokens: make([]ApiToken, len(tokens)),
				}

				for i, token := range tokens {
					res.Tokens[i] = ApiToken{
						Id:   token.Id,
						Name: token.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RemoveApiToken",
			Method: http.MethodDelete,
			Path:   "/user/apitoken/:id",
			Errors: []pyrin.ErrorType{ErrTypeApiTokenNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				tokenId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				token, err := app.DB().GetApiTokenById(ctx, tokenId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ApiTokenNotFound()
					}

					return nil, err
				}

				if token.UserId != user.Id {
					return nil, ApiTokenNotFound()
				}

				err = app.DB().RemoveApiToken(ctx, tokenId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
