package handlers

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandleGetAlbums(c echo.Context) error {
	f := c.QueryParam("filter")

	albums, err := h.db.GetAllAlbums(c.Request().Context(), f)
	if err != nil {
		return err
	}

	res := types.GetAlbums{
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

func (h *Handlers) HandleGetAlbumById(c echo.Context) error {
	id := c.Param("id")
	album, err := h.db.GetAlbumById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetAlbumById{
		Album: types.Album{
			Id:       album.Id,
			Name:     album.Name,
			CoverArt: ConvertAlbumCoverURL(c, album.CoverArt),
			ArtistId: album.ArtistId,
		},
	}))
}

func (h *Handlers) HandleGetAlbumTracksById(c echo.Context) error {
	id := c.Param("id")

	album, err := h.db.GetAlbumById(c.Request().Context(), id)
	if err != nil {
		return err
	}

	tracks, err := h.db.GetTracksByAlbum(c.Request().Context(), album.Id)
	if err != nil {
		return err
	}

	res := types.GetAlbumTracksById{
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

func (h *Handlers) InstallAlbumHandlers(group Group) {
	group.GET(
		"GetAlbums", "/albums", 
		types.GetAlbums{}, nil,
		h.HandleGetAlbums, 
	)

	group.GET(
		"GetAlbumById", "/albums/:id", 
		types.GetAlbumById{}, nil,
		h.HandleGetAlbumById, 
	)

	group.GET(
		"GetAlbumTracks", "/albums/:id/tracks", 
		types.GetAlbumTracksById{}, nil,
		h.HandleGetAlbumTracksById,
	)
}
