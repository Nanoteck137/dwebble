package apis

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/kr/pretty"
	"github.com/nanoteck137/dwebble/core"
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
		validate.Field(&b.QuickPlaylist,
			validate.Required.When(b.QuickPlaylist != nil),
		),
	)
}

type AddToUserQuickPlaylistBody struct {
	Tracks []string
}

func (b *AddToUserQuickPlaylistBody) Transform() {
	// b.Tracks = *transform.DiscardEmptyStringEntries()
}

func InstallUserHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "UpdateUserSettings",
			Method:   http.MethodPatch,
			Path:     "/user/settings",
			DataType: nil,
			BodyType: nil,
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

					_, err := app.DB().GetPlaylistById(context.TODO(), id)
					if err != nil {
						// TODO(patrik): Handle error
						return nil, err
					}

					settings.QuickPlaylist = sql.NullString{
						String: id,
						Valid:  true,
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
			BodyType: nil,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[AddToUserQuickPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				pretty.Println(body)

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				if user.QuickPlaylist.Valid {
					err := app.DB().AddItemsToPlaylist(context.TODO(), user.QuickPlaylist.String, body.Tracks)
					if err != nil {
						return nil, err
					}
				}

				return nil, nil
			},
		},
	)
}
