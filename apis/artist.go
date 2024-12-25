package apis

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type ArtistInfo struct {
	Id   string `json:"id"`
	Name Name   `json:"name"`
}

func ConvertDBExtraArtists(extras database.ExtraArtists) []ArtistInfo {
	res := []ArtistInfo{}
	for _, extraArtist := range extras {
		res = append(res, ArtistInfo{
			Id: extraArtist.Id,
			Name: Name{
				Default: extraArtist.Name,
				Other:   extraArtist.OtherName,
			},
		})
	}

	return res
}

type Artist struct {
	Id   string `json:"id"`
	Name Name   `json:"name"`

	Picture types.Images `json:"picture"`

	Tags []string `json:"tags"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

func ConvertDBArtist(c pyrin.Context, artist database.Artist) Artist {
	return Artist{
		Id: artist.Id,
		Name: Name{
			Default: artist.Name,
			Other:   ConvertSqlNullString(artist.OtherName),
		},
		Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
		Tags:    utils.SplitString(artist.Tags.String),
		Created: artist.Created,
		Updated: artist.Updated,
	}
}

type GetArtists struct {
	Artists []Artist `json:"artists"`
}

type GetArtistById struct {
	Artist
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type EditArtistBody struct {
	Name *string   `json:"name"`
	Tags *[]string `json:"tags,omitempty"`
}

func (b *EditArtistBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
	b.Tags = DiscardEntriesStringArrayPtr(b.Tags)
}

func (b EditArtistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

type CreateArtist struct {
	Id string `json:"id"`
}

type CreateArtistBody struct {
	Name      string `json:"name"`
	OtherName string `json:"otherName"`
}

func (b *CreateArtistBody) Transform() {
	b.Name = transform.String(b.Name)
	b.OtherName = transform.String(b.OtherName)
}

func (b CreateArtistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required),
	)
}

func InstallArtistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetArtists",
			Method:       http.MethodGet,
			Path:         "/artists",
			ResponseType: GetArtists{},
			Errors:       []pyrin.ErrorType{ErrTypeInvalidFilter, ErrTypeInvalidSort},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				filter := q.Get("filter")
				sort := q.Get("sort")

				artists, err := app.DB().GetAllArtists(c.Request().Context(), filter, sort)
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					if errors.Is(err, database.ErrInvalidSort) {
						return nil, InvalidSort(err)
					}

					return nil, err
				}

				res := GetArtists{
					Artists: make([]Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = ConvertDBArtist(c, artist)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "SearchArtists",
			Method:       http.MethodGet,
			Path:         "/artists/search",
			ResponseType: GetArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))
				artists, err := app.DB().SearchArtists(query)
				if err != nil {
					return nil, err
				}

				res := GetArtists{
					Artists: make([]Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = ConvertDBArtist(c, artist)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetArtistById",
			Method:       http.MethodGet,
			Path:         "/artists/:id",
			ResponseType: GetArtistById{},
			Errors:       []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				artist, err := app.DB().GetArtistById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}
					return nil, err
				}

				return GetArtistById{
					Artist: ConvertDBArtist(c, artist),
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "GetArtistAlbums",
			Method:       http.MethodGet,
			Path:         "/artists/:id/albums",
			ResponseType: GetArtistAlbumsById{},
			Errors:       []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				artist, err := app.DB().GetArtistById(c.Request().Context(), id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}

					return nil, err
				}

				albums, err := app.DB().GetAlbumsByArtist(c.Request().Context(), artist.Id)
				if err != nil {
					return nil, err
				}

				res := GetArtistAlbumsById{
					Albums: make([]Album, len(albums)),
				}

				for i, album := range albums {
					res.Albums[i] = ConvertDBAlbum(c, album)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "EditArtist",
			Method:   http.MethodPatch,
			Path:     "/artists/:id",
			BodyType: EditArtistBody{},
			Errors:   []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				body, err := pyrin.Body[EditArtistBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				artist, err := app.DB().GetArtistById(ctx, id)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, ArtistNotFound()
					}

					return nil, err
				}
				changes := database.ArtistChanges{}

				if body.Name != nil {
					changes.Name = types.Change[string]{
						Value:   *body.Name,
						Changed: artist.Name != *body.Name,
					}
				}

				err = app.DB().UpdateArtist(ctx, artist.Id, changes)
				if err != nil {
					return nil, err
				}

				// TODO(patrik): Use transaction
				if body.Tags != nil {
					tags := *body.Tags

					err = app.DB().RemoveAllTagsFromArtist(ctx, artist.Id)
					if err != nil {
						return nil, err
					}

					for _, tag := range tags {
						slug := utils.Slug(tag)

						err := app.DB().CreateTag(ctx, slug)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}

						err = app.DB().AddTagToArtist(ctx, slug, artist.Id)
						if err != nil && !errors.Is(err, database.ErrItemAlreadyExists) {
							return nil, err
						}
					}
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:         "CreateArtist",
			Method:       http.MethodPost,
			Path:         "/artists",
			ResponseType: CreateArtist{},
			BodyType:     CreateArtistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				body, err := pyrin.Body[CreateArtistBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				artist, err := app.DB().CreateArtist(ctx, database.CreateArtistParams{
					Name: body.Name,
					OtherName: sql.NullString{
						String: body.OtherName,
						Valid:  body.OtherName != "",
					},
				})
				if err != nil {
					return nil, err
				}

				return CreateArtist{
					Id: artist.Id,
				}, nil
			},
		},

		// TODO(patrik): Fix
		// pyrin.ApiHandler{
		// 	Name:        "ChangeArtistPicture",
		// 	Method:      http.MethodPatch,
		// 	Path:        "/artists/:id/picture",
		// 	RequireForm: true,
		// 	Errors:      []pyrin.ErrorType{ErrTypeArtistNotFound},
		// 	HandlerFunc: func(c pyrin.Context) (any, error) {
		// 		id := c.Param("id")
		//
		// 		ctx := context.TODO()
		//
		// 		err := c.Request().ParseMultipartForm(defaultMemory)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		form := c.Request().MultipartForm
		//
		// 		artist, err := app.DB().GetArtistById(ctx, id)
		// 		if err != nil {
		// 			return nil, ArtistNotFound()
		// 		}
		//
		// 		formFile := form.File["picture"][0]
		//
		// 		// TODO(patrik): Check the file ext for png jpeg jpg
		//
		// 		f, err := formFile.Open()
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		artistDir := app.WorkDir().Artist(artist.Id)
		//
		// 		name := "picture-original" + path.Ext(formFile.Filename)
		// 		dstName := path.Join(artistDir, name)
		// 		dst, err := os.OpenFile(dstName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		defer dst.Close()
		//
		// 		_, err = io.Copy(dst, f)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		dst.Close()
		//
		// 		i := path.Join(artistDir, "picture-128.png")
		// 		err = utils.CreateResizedImage(dstName, i, 128)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		i = path.Join(artistDir, "picture-256.png")
		// 		err = utils.CreateResizedImage(dstName, i, 256)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		i = path.Join(artistDir, "picture-512.png")
		// 		err = utils.CreateResizedImage(dstName, i, 512)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		err = app.DB().UpdateArtist(ctx, artist.Id, database.ArtistChanges{
		// 			Picture: types.Change[sql.NullString]{
		// 				Value: sql.NullString{
		// 					String: name,
		// 					Valid:  true,
		// 				},
		// 				Changed: true,
		// 			},
		// 		})
		// 		if err != nil {
		// 			return nil, err
		// 		}
		//
		// 		return nil, nil
		// 	},
		// },
	)
}
