package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (api *ApiConfig) HandleGetTracks(c echo.Context) error {
	tracks, err := api.db.GetAllTracks(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.GetTracks{
		Tracks: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Tracks[i] = types.Track{
			Id:                track.Id,
			Number:            track.Number,
			Name:              track.Name,
			CoverArt:          ConvertTrackCoverURL(c, track.CoverArt),
			Duration:          track.Duration,
			BestQualityFile:   ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
			MobileQualityFile: ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
			AlbumId:           track.AlbumId,
			ArtistId:          track.ArtistId,
			AlbumName:         track.AlbumName,
			ArtistName:        track.ArtistName,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *ApiConfig) HandleGetTrackById(c echo.Context) error {
	id := c.Param("id")
	track, err := api.db.GetTrackById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetTrackById{
		Id:                track.Id,
		Number:            track.Number,
		Name:              track.Name,
		CoverArt:          ConvertTrackCoverURL(c, track.CoverArt),
		Duration:          track.Duration,
		BestQualityFile:   ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
		MobileQualityFile: ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
		AlbumId:           track.AlbumId,
		ArtistId:          track.ArtistId,
	}))
}

func InstallTrackHandlers(group *echo.Group, apiConfig *ApiConfig) {
	group.GET("/tracks", apiConfig.HandleGetTracks)
	group.GET("/tracks/:id", apiConfig.HandleGetTrackById)
}
