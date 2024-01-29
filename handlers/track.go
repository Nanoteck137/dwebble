package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/nanoteck137/dwebble/types"
)

// HandleGetTracks godoc
//
//	@Summary		Get all tracks
//	@Description	Get all tracks
//	@Tags			tracks
//	@Produce		json
//	@Success		200	{object}	types.ApiResponse[types.ApiGetTracksData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/tracks [get]
func (api *ApiConfig) HandleGetTracks(c *fiber.Ctx) error {
	tracks, err := api.queries.GetAllTracks(c.UserContext())
	if err != nil {
		return err
	}

	result := types.ApiGetTracksData{
		Tracks: make([]types.ApiGetTracksDataTrackItem, len(tracks)),
	}

	for i, track := range tracks {
		result.Tracks[i] = types.ApiGetTracksDataTrackItem{
			ApiTrack: types.ApiTrack{
				Id:                track.ID,
				Number:            track.TrackNumber,
				Name:              track.Name,
				CoverArt:          ConvertURL(c, "/images/"+track.CoverArt),
				BestQualityFile:   ConvertURL(c, "/tracks/"+track.BestQualityFile),
				MobileQualityFile: ConvertURL(c, "/tracks/"+track.MobileQualityFile),
				AlbumId:           track.AlbumID,
				ArtistId:          track.ArtistID,
			},
			AlbumName:  track.AlbumName,
			ArtistName: track.ArtistName,
		}
	}

	return c.JSON(types.NewApiResponse(result))
}

// HandleGetTrackById godoc
//
//	@Summary		Get track by id
//	@Description	Get track by id
//	@Tags			tracks
//	@Produce		json
//	@Param			id	path		string	true	"Track Id"
//	@Success		200	{object}	types.ApiResponse[types.ApiGetTrackByIdData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		404	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/tracks/{id} [get]
func (api *ApiConfig) HandleGetTrackById(c *fiber.Ctx) error {
	id := c.Params("id")

	track, err := api.queries.GetTrack(c.UserContext(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return types.ApiNotFoundError(fmt.Sprintf("No track with id: '%s'", id))
		} else {
			return err
		}
	}

	return c.JSON(types.NewApiResponse(types.ApiGetTrackByIdData{
		Id:                track.ID,
		Number:            track.TrackNumber,
		Name:              track.Name,
		CoverArt:          ConvertURL(c, "/images/"+track.CoverArt),
		BestQualityFile:   ConvertURL(c, "/tracks/"+track.BestQualityFile),
		MobileQualityFile: ConvertURL(c, "/tracks/"+track.MobileQualityFile),
		AlbumId:           track.AlbumID,
		ArtistId:          track.ArtistID,
	}))
}

func InstallTrackHandlers(router *fiber.App, apiConfig *ApiConfig) {
	router.Get("/tracks", apiConfig.HandleGetTracks)
	router.Get("/tracks/:id", apiConfig.HandleGetTrackById)
}
