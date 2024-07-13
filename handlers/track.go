package handlers

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandleGetTracks(c echo.Context) error {
	f := c.QueryParam("filter")
	s := c.QueryParam("sort")

	tracks, err := h.db.GetAllTracks(c.Request().Context(), f, s)
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
			Tags:              strings.Split(track.Tags, ","),
			Genres:            strings.Split(track.Genres, ","),
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (h *Handlers) HandleGetTrackById(c echo.Context) error {
	id := c.Param("id")
	track, err := h.db.GetTrackById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetTrackById{
		Track: types.Track{
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
			Tags:              strings.Split(track.Tags, ","),
			Genres:            strings.Split(track.Genres, ","),
		},
	}))
}

func (h *Handlers) InstallTrackHandlers(group Group) {
	group.GET(
		"GetTracks", "/tracks", 
		h.HandleGetTracks, 
		types.GetTracks{}, nil,
	)

	group.GET(
		"GetTrackById", "/tracks/:id", 
		h.HandleGetTrackById,
		types.GetTrackById{}, nil,
	)
}
