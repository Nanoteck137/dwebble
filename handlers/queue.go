package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandlePostQueue(c echo.Context) error {
	tracks, err := h.db.GetAllTracks(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.PostQueue{
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

func (h *Handlers) InstallQueueHandlers(group *echo.Group) {
	group.POST("/queue", h.HandlePostQueue)
}
