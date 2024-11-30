package apis

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/transform"
	"github.com/nanoteck137/validate"
)

type Artist struct {
	Id      string       `json:"id"`
	Name    string       `json:"name"`
	Picture types.Images `json:"picture"`
	Created int64        `json:"created"`
	Updated int64        `json:"updated"`
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
	Name *string `json:"name"`
}

func (b *EditArtistBody) Transform() {
	b.Name = transform.StringPtr(b.Name)
}

func (b EditArtistBody) Validate() error {
	return validate.ValidateStruct(&b,
		validate.Field(&b.Name, validate.Required.When(b.Name != nil)),
	)
}

func InstallArtistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetArtists",
			Method:   http.MethodGet,
			Path:     "/artists",
			DataType: GetArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				artists, err := app.DB().GetAllArtists(c.Request().Context())
				if err != nil {
					return nil, err
				}

				res := GetArtists{
					Artists: make([]Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = Artist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
						Created: artist.Created,
						Updated: artist.Updated,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SearchArtists",
			Method:   http.MethodGet,
			Path:     "/artists/search",
			DataType: GetArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				q := c.Request().URL.Query()

				query := strings.TrimSpace(q.Get("query"))

				var err error
				var artists []database.Artist

				if query != "" {
					artists, err = app.DB().SearchArtists(query)
					if err != nil {
						return nil, err
					}
				}

				res := GetArtists{
					Artists: make([]Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = Artist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
						Created: artist.Created,
						Updated: artist.Updated,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetArtistById",
			Method:   http.MethodGet,
			Path:     "/artists/:id",
			DataType: GetArtistById{},
			Errors:   []pyrin.ErrorType{ErrTypeArtistNotFound},
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
					Artist: Artist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
						Created: artist.Created,
						Updated: artist.Updated,
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetArtistAlbums",
			Method:   http.MethodGet,
			Path:     "/artists/:id/albums",
			DataType: GetArtistAlbumsById{},
			Errors:   []pyrin.ErrorType{ErrTypeArtistNotFound},
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
					// TODO(patrik): Fix not all fields filled in
					res.Albums[i] = Album{
						Id:         album.Id,
						Name:       album.Name,
						CoverArt:   utils.ConvertAlbumCoverURL(c, album.Id, album.CoverArt),
						ArtistId:   album.ArtistId,
						ArtistName: album.ArtistName,
					}
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

				// TODO(patrik): Create helper function for this
				var name types.Change[string]
				if body.Name != nil {
					name = types.Change[string]{
						Value:   *body.Name,
						Changed: artist.Name != *body.Name,
					}
				}

				err = app.DB().UpdateArtist(ctx, artist.Id, database.ArtistChanges{
					Name: name,
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:        "ChangeArtistPicture",
			Method:      http.MethodPatch,
			Path:        "/artists/:id/picture",
			RequireForm: true,
			Errors:      []pyrin.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				id := c.Param("id")

				ctx := context.TODO()

				err := c.Request().ParseMultipartForm(defaultMemory)
				if err != nil {
					return nil, err
				}

				form := c.Request().MultipartForm

				artist, err := app.DB().GetArtistById(ctx, id)
				if err != nil {
					return nil, ArtistNotFound()
				}

				formFile := form.File["picture"][0]

				// TODO(patrik): Check the file ext for png jpeg jpg

				f, err := formFile.Open()
				if err != nil {
					return nil, err
				}

				artistDir := app.WorkDir().Artist(artist.Id)

				// TODO(patrik): This should not be here
				err = os.Mkdir(artistDir, 0755)
				if err != nil && !os.IsExist(err) {
					return nil, err
				}

				name := "picture-original" + path.Ext(formFile.Filename)
				dstName := path.Join(artistDir, name)
				dst, err := os.OpenFile(dstName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
				if err != nil {
					return nil, err
				}
				defer dst.Close()

				_, err = io.Copy(dst, f)
				if err != nil {
					return nil, err
				}

				dst.Close()

				i := path.Join(artistDir, "picture-128.png")
				err = utils.CreateResizedImage(dstName, i, 128)
				if err != nil {
					return nil, err
				}

				i = path.Join(artistDir, "picture-256.png")
				err = utils.CreateResizedImage(dstName, i, 256)
				if err != nil {
					return nil, err
				}

				i = path.Join(artistDir, "picture-512.png")
				err = utils.CreateResizedImage(dstName, i, 512)
				if err != nil {
					return nil, err
				}

				err = app.DB().UpdateArtist(ctx, artist.Id, database.ArtistChanges{
					Picture: types.Change[sql.NullString]{
						Value: sql.NullString{
							String: name,
							Valid:  true,
						},
						Changed: true,
					},
				})
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
