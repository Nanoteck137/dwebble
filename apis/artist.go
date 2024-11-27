package apis

import (
	"context"
	"database/sql"
	"encoding/json"
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
	"github.com/nanoteck137/pyrin/tools/validate"
)

var _ pyrin.Body = (*EditArtistBody)(nil)

type EditArtistBody struct {
	Name *string `json:"name"`
}

func (b EditArtistBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func InstallArtistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetArtists",
			Method:   http.MethodGet,
			Path:     "/artists",
			DataType: types.GetArtists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				artists, err := app.DB().GetAllArtists(c.Request().Context())
				if err != nil {
					return nil, err
				}

				res := types.GetArtists{
					Artists: make([]types.Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = types.Artist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "SearchArtists",
			Method:   http.MethodGet,
			Path:     "/artists/search",
			DataType: types.GetArtists{},
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

				res := types.GetArtists{
					Artists: make([]types.Artist, len(artists)),
				}

				for i, artist := range artists {
					res.Artists[i] = types.Artist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetArtistById",
			Method:   http.MethodGet,
			Path:     "/artists/:id",
			DataType: types.GetArtistById{},
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

				return types.GetArtistById{
					Artist: types.Artist{
						Id:      artist.Id,
						Name:    artist.Name,
						Picture: utils.ConvertArtistPicture(c, artist.Id, artist.Picture),
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetArtistAlbums",
			Method:   http.MethodGet,
			Path:     "/artists/:id/albums",
			DataType: types.GetArtistAlbumsById{},
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

				res := types.GetArtistAlbumsById{
					Albums: make([]types.Album, len(albums)),
				}

				for i, album := range albums {
					// TODO(patrik): Fix not all fields filled in
					res.Albums[i] = types.Album{
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

				// TODO(patrik): Change this to pyrin version
				var body EditArtistBody
				d := json.NewDecoder(c.Request().Body)
				err := d.Decode(&body)
				if err != nil {
					return nil, err
				}

				var name string
				if body.Name != nil {
					name = strings.TrimSpace(*body.Name)
				}

				ctx := context.TODO()

				artist, err := app.DB().GetArtistById(ctx, id)
				if err != nil {
					return nil, ArtistNotFound()
				}

				err = app.DB().UpdateArtist(ctx, artist.Id, database.ArtistChanges{
					Name: types.Change[string]{
						Value: name,
						Changed: name != "" && artist.Name != name,
					},
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
