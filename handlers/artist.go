package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandleGetArtists(c echo.Context) error {
	artists, err := h.db.GetAllArtists(c.Request().Context())
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

func (h *Handlers) HandleGetArtistById(c echo.Context) error {
	id := c.Param("id")
	artist, err := h.db.GetArtistById(c.Request().Context(), id)
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

func (h *Handlers) HandleGetArtistAlbumsById(c echo.Context) error {
	id := c.Param("id")

	artist, err := h.db.GetArtistById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	albums, err := h.db.GetAlbumsByArtist(c.Request().Context(), artist.Id)
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

func (h *Handlers) InstallArtistHandlers(group Group) {
	group.Register(
		Handler{
			Name:        "GetArtists",
			Method:      http.MethodGet,
			Path:        "/artists",
			DataType:    types.GetArtists{},
			BodyType:    nil,
			HandlerFunc: h.HandleGetArtists,
		},

		Handler{
			Name:        "GetArtistById",
			Path:        "/artists/:id",
			Method:      http.MethodGet,
			DataType:    types.GetArtistById{},
			BodyType:    nil,
			HandlerFunc: h.HandleGetArtistById,
		},

		Handler{
			Name:        "GetArtistAlbums",
			Path:        "/artists/:id/albums",
			Method:      http.MethodGet,
			DataType:    types.GetArtistAlbumsById{},
			BodyType:    nil,
			HandlerFunc: h.HandleGetArtistAlbumsById,
		},
	)
}
