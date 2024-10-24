package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

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
						Picture: utils.ConvertArtistPictureURL(c, artist.Picture),
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
						Picture: utils.ConvertArtistPictureURL(c, artist.Picture),
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
	)
}
