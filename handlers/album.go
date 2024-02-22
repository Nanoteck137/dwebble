package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// func (api *ApiConfig) HandleGetAlbums(c *fiber.Ctx) error {
// 	albums, err := api.queries.GetAllAlbums(c.UserContext())
// 	if err != nil {
// 		return err
// 	}
//
// 	result := types.ApiGetAlbumsData{
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
//
// type CreateAlbumBody struct {
// 	Name   string `json:"name" form:"name" validate:"required"`
// 	Artist string `json:"artist" form:"artist" validate:"required"`
// }
//
// // HandleGetAlbumById godoc
// //
// //	@Summary		Get album by id
// //	@Description	Get album by id
// //	@Tags			albums
// //	@Produce		json
// //	@Param			id	path		string	true	"Album Id"
// //	@Success		200	{object}	types.ApiResponse[types.ApiGetAlbumByIdData]
// //	@Failure		400	{object}	types.ApiError
// //	@Failure		404	{object}	types.ApiError
// //	@Failure		500	{object}	types.ApiError
// //	@Router			/albums/{id} [get]
// func (api *ApiConfig) HandleGetAlbumById(c *fiber.Ctx) error {
// 	id := c.Params("id")
//
// 	album, err := api.queries.GetAlbum(c.UserContext(), id)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return types.ApiNotFoundError(fmt.Sprintf("No album with id: '%s'", id))
// 		} else {
// 			return err
// 		}
// 	}
//
// 	return c.JSON(types.NewApiResponse(types.ApiGetAlbumByIdData{
// 		Id:       album.ID,
// 		Name:     album.Name,
// 		CoverArt: album.CoverArt,
// 		ArtistId: album.ArtistID,
// 	}))
// }
//
// // HandleGetAlbumTracksById godoc
// //
// //	@Summary		Get all tracks from album
// //	@Description	Get all tracks from album
// //	@Tags			albums
// //	@Produce		json
// //	@Param			id	path		string	true	"Artist Id"
// //	@Success		200	{object}	types.ApiResponse[types.ApiGetAlbumTracksByIdData]
// //	@Failure		400	{object}	types.ApiError
// //	@Failure		500	{object}	types.ApiError
// //	@Router			/albums/{id}/tracks [get]
// func (api *ApiConfig) HandleGetAlbumTracksById(c *fiber.Ctx) error {
// 	id := c.Params("id")
//
// 	tracks, err := api.queries.GetTracksByAlbum(c.UserContext(), id)
// 	if err != nil {
// 		return err
// 	}
//
// 	result := types.ApiGetAlbumTracksByIdData{
// 		Tracks: make([]types.ApiTrack, len(tracks)),
// 	}
//
// 	for i, track := range tracks {
// 		result.Tracks[i] = types.ApiTrack{
// 			Id:                track.ID,
// 			Number:            track.TrackNumber,
// 			Name:              track.Name,
// 			CoverArt:          ConvertURL(c, "/images/"+track.CoverArt),
// 			BestQualityFile:   ConvertURL(c, "/tracks/"+track.BestQualityFile),
// 			MobileQualityFile: ConvertURL(c, "/tracks/"+track.MobileQualityFile),
// 			AlbumId:           track.AlbumID,
// 			ArtistId:          track.ArtistID,
// 			// TODO(patrik): Fix
// 			// AlbumName:         track.AlbumName,
// 			// ArtistName:        track.ArtistName,
// 		}
// 	}
//
// 	return c.JSON(types.NewApiResponse(result))
// }

func InstallAlbumHandlers(router *fiber.App, apiConfig *ApiConfig) {
	// router.Get("/albums", apiConfig.HandleGetAlbums)
	// router.Get("/albums/:id", apiConfig.HandleGetAlbumById)
	// router.Get("/albums/:id/tracks", apiConfig.HandleGetAlbumTracksById)
}
