package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (api *ApiConfig) HandleGetArtists(c echo.Context) error {
	artists, err := api.db.GetAllArtists(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.GetArtists{
		Artists: make([]types.GetArtistsItem, len(artists)),
	}

	for i, artist := range artists {
		pictureName := "default_artist.png"
		if artist.Picture.Valid {
			pictureName = artist.Picture.String
		}

		res.Artists[i] = types.GetArtistsItem{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: ConvertURL(c, "/images/"+pictureName),
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *ApiConfig) HandleGetArtistById(c echo.Context) error {
	id := c.Param("id")
	artist, err := api.db.GetArtistById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetArtistById{
		Id:      artist.Id,
		Name:    artist.Name,
		Picture: artist.Picture.String,
	}))
}

func (api *ApiConfig) HandleGetArtistAlbumsById(c echo.Context) error {
	id := c.Param("id")

	artist, err := api.db.GetArtistById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	albums, err := api.db.GetAlbumsByArtist(c.Request().Context(), artist.Id);
	if err != nil {
		return err
	}

	res := types.GetArtistAlbumsById{
		Albums: make([]types.GetArtistAlbumsByIdItem, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.GetArtistAlbumsByIdItem{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: album.CoverArt.String,
			ArtistId: album.ArtistId,
		}
	}


	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func InstallArtistHandlers(group *echo.Group, apiConfig *ApiConfig) {
	group.GET("/artists", apiConfig.HandleGetArtists)
	group.GET("/artists/:id", apiConfig.HandleGetArtistById)
	group.GET("/artists/:id/albums", apiConfig.HandleGetArtistAlbumsById)
}
