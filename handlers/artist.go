package handlers

import (
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (api *ApiConfig) HandleGetArtists(c echo.Context) error {
	artists, err := api.db.GetAllArtists(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.ApiGetArtistsData{
		Artists: make([]types.ApiArtist, len(artists)),
	}

	for i, artist := range artists {
		res.Artists[i] = types.ApiArtist{
			Id:      artist.Id,
			Name:    artist.Name,
			Picture: artist.Picture,
		}
	}

	return c.JSON(200, types.NewApiResponse(res))
}

func (api *ApiConfig) HandleGetArtistById(c echo.Context) error {
	id := c.Param("id")
	artist, err := api.db.GetArtistById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiResponse(types.ApiGetArtistByIdData{
		Id:      artist.Id,
		Name:    artist.Name,
		Picture: artist.Picture,
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

	pretty.Println(albums)

	res := types.ApiGetArtistAlbumsByIdData{
		Albums: make([]types.ApiAlbum, len(albums)),
	}

	for i, album := range albums {
		res.Albums[i] = types.ApiAlbum{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: album.CoverArt,
			ArtistId: album.ArtistId,
		}
	}


	return c.JSON(200, types.NewApiResponse(res))
}

// func (api *ApiConfig) HandleGetArtistById(c *fiber.Ctx) error {
// 	id := c.Params("id")
//
// 	artist, err := api.queries.GetArtist(c.UserContext(), id)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return types.ApiNotFoundError(fmt.Sprintf("No artist with id: '%s'", id))
// 		} else {
// 			return err
// 		}
// 	}
//
// 	return c.JSON(types.NewApiResponse(types.ApiGetArtistByIdData{
// 		Id:      artist.ID,
// 		Name:    artist.Name,
// 		Picture: ConvertURL(c, "/images/"+artist.Picture),
// 	}))
// }
//
// func (api *ApiConfig) HandleGetArtistAlbumsById(c *fiber.Ctx) error {
// 	id := c.Params("id")
//
// 	var albums []database.Album
//
// 	nameQuery := c.Query("name")
// 	if nameQuery != "" {
// 		a, err := api.queries.GetAlbumsByArtistAndName(c.UserContext(), database.GetAlbumsByArtistAndNameParams{
// 			ArtistID: id,
// 			Name:     nameQuery,
// 		})
//
// 		if err != nil {
// 			return err
// 		}
//
// 		albums = a
// 	} else {
// 		a, err := api.queries.GetAlbumsByArtist(c.UserContext(), id)
// 		if err != nil {
// 			return err
// 		}
//
// 		albums = a
// 	}
//
// 	result := types.ApiGetArtistAlbumsByIdData{
// 		Albums: make([]types.ApiAlbum, len(albums)),
// 	}
//
// 	for i, album := range albums {
// 		result.Albums[i] = types.ApiAlbum{
// 			Id:       album.ID,
// 			Name:     album.Name,
// 			CoverArt: ConvertURL(c, "/images/"+album.CoverArt),
// 			ArtistId: album.ArtistID,
// 		}
// 	}
//
// 	return c.JSON(types.NewApiResponse(result))
// }

func InstallArtistHandlers(group *echo.Group, apiConfig *ApiConfig) {
	group.GET("/artists", apiConfig.HandleGetArtists)
	group.GET("/artists/:id", apiConfig.HandleGetArtistById)
	group.GET("/artists/:id/albums", apiConfig.HandleGetArtistAlbumsById)
}
