package apis

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/tools/utils"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/api"
)

type artistApi struct {
	app core.App
}

func (api *artistApi) HandleGetArtists(c echo.Context) error {
	artists, err := api.app.DB().GetAllArtists(c.Request().Context())
	if err != nil {
		return err
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

	return c.JSON(200, SuccessResponse(res))
}

func (api *artistApi) HandleGetArtistById(c echo.Context) error {
	id := c.Param("id")
	artist, err := api.app.DB().GetArtistById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ArtistNotFound()
		}
		return err
	}

	return c.JSON(200, SuccessResponse(types.GetArtistById{
		Artist: types.Artist{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: utils.ConvertArtistPictureURL(c, artist.Picture),
		},
	}))
}

func (api *artistApi) HandleGetArtistAlbumsById(c echo.Context) error {
	id := c.Param("id")

	artist, err := api.app.DB().GetArtistById(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, database.ErrItemNotFound) {
			return ArtistNotFound()
		}
		return err
	}

	albums, err := api.app.DB().GetAlbumsByArtist(c.Request().Context(), artist.Id)
	if err != nil {
		return err
	}

	res := types.GetArtistAlbumsById{
		Albums: make([]types.Album, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.Album{
			Id:         album.Id,
			Name:       album.Name,
			CoverArt:   utils.ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId:   album.ArtistId,
			ArtistName: album.ArtistName,
		}
	}

	return c.JSON(200, SuccessResponse(res))
}

func InstallArtistHandlers(app core.App, group Group) {
	a := artistApi{app: app}

	group.Register(
		Handler{
			Name:        "GetArtists",
			Method:      http.MethodGet,
			Path:        "/artists",
			DataType:    types.GetArtists{},
			BodyType:    nil,
			HandlerFunc: a.HandleGetArtists,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetArtistById",
			Method:      http.MethodGet,
			Path:        "/artists/:id",
			DataType:    types.GetArtistById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: a.HandleGetArtistById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetArtistAlbums",
			Method:      http.MethodGet,
			Path:        "/artists/:id/albums",
			DataType:    types.GetArtistAlbumsById{},
			BodyType:    nil,
			Errors:      []api.ErrorType{ErrTypeArtistNotFound},
			HandlerFunc: a.HandleGetArtistAlbumsById,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
