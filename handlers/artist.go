package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
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
			Picture: ConvertArtistPictureURL(c, artist.Picture),
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *artistApi) HandleGetArtistById(c echo.Context) error {
	id := c.Param("id")
	artist, err := api.app.DB().GetArtistById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetArtistById{
		Artist: types.Artist{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: ConvertArtistPictureURL(c, artist.Picture),
		},
	}))
}

func (api *artistApi) HandleGetArtistAlbumsById(c echo.Context) error {
	id := c.Param("id")

	artist, err := api.app.DB().GetArtistById(c.Request().Context(), id)
	if err != nil {
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
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId: album.ArtistId,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func InstallArtistHandlers(app core.App, group Group) {
	api := artistApi{app: app}

	group.Register(
		Handler{
			Name:        "GetArtists",
			Method:      http.MethodGet,
			Path:        "/artists",
			DataType:    types.GetArtists{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetArtists,
		},

		Handler{
			Name:        "GetArtistById",
			Path:        "/artists/:id",
			Method:      http.MethodGet,
			DataType:    types.GetArtistById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetArtistById,
		},

		Handler{
			Name:        "GetArtistAlbums",
			Path:        "/artists/:id/albums",
			Method:      http.MethodGet,
			DataType:    types.GetArtistAlbumsById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetArtistAlbumsById,
		},
	)
}
