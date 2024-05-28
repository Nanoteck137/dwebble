package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandlePostPlaylist(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistBody](c, types.PostPlaylistBodySchema)
	if err != nil {
		return err
	}

	playlist, err := h.db.CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
		Name:    body.Name,
		OwnerId: user.Id,
	})
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.PostPlaylist{
		Id:   playlist.Id,
		Name: playlist.Name,
	}))
}

func (h *Handlers) HandlePostPlaylistItemsById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := h.User(c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistItemsByIdBody](c, types.PostPlaylistItemsByIdBodySchema)
	if err != nil {
		return err
	}

	playlist, err := h.db.GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		return types.ErrNoPlaylist
	}

	err = h.db.AddItemsToPlaylist(c.Request().Context(), playlist.Id, body.Tracks)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (h *Handlers) HandleGetPlaylistById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := h.User(c)
	if err != nil {
		return err
	}

	playlist, err := h.db.GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		return types.ErrNoPlaylist
	}

	items, err := h.db.GetPlaylistItems(c.Request().Context(), playlist.Id)
	if err != nil {
		return err
	}

	var tracks []types.Track
	for _, item := range items {
		// TODO(patrik): Optimize
		track, err := h.db.GetTrackById(c.Request().Context(), item.TrackId)
		if err != nil {
			return err
		}

		tracks = append(tracks, types.Track{
			Id:                track.Id,
			Number:            item.ItemIndex,
			Name:              track.Name,
			CoverArt:          ConvertTrackCoverURL(c, track.CoverArt),
			Duration:          track.Duration,
			BestQualityFile:   ConvertURL(c, "/tracks/original/"+track.BestQualityFile),
			MobileQualityFile: ConvertURL(c, "/tracks/mobile/"+track.MobileQualityFile),
			AlbumId:           track.AlbumId,
			ArtistId:          track.ArtistId,
			AlbumName:         track.AlbumName,
			ArtistName:        track.ArtistName,
		})
	}

	return c.JSON(200, types.NewApiSuccessResponse(types.GetPlaylistById{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
		Items: tracks,
	}))
}

func (h *Handlers) InstallPlaylistHandlers(group *echo.Group) {
	group.POST("/playlist", h.HandlePostPlaylist)
	group.POST("/playlist/:id/items", h.HandlePostPlaylistItemsById)
	group.GET("/playlist/:id", h.HandleGetPlaylistById)
}
