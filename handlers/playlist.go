package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

func (h *Handlers) HandleGetPlaylists(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return err
	}

	playlists, err := h.db.GetPlaylistsByUser(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	res := types.GetPlaylists{
		Playlists: make([]types.Playlist, len(playlists)),
	}

	for i, playlist := range playlists {
		res.Playlists[i] = types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		}
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

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
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
	}))
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

	tracks := []types.Track{}
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

func (h *Handlers) HandleDeletePlaylistItemsById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := h.User(c)
	if err != nil {
		return err
	}

	body, err := Body[types.DeletePlaylistItemsByIdBody](c, types.DeletePlaylistItemsByIdBodySchema)
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

	err = h.db.DeleteItemsFromPlaylist(c.Request().Context(), playlist.Id, body.TrackIndices)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (h *Handlers) HandlePostPlaylistsItemsMoveById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := h.User(c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistsItemMoveByIdBody](c, types.PostPlaylistsItemMoveByIdBodySchema)
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

	err = h.db.MovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId, body.ToIndex)
	if err != nil {
		return err
	}

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (h *Handlers) InstallPlaylistHandlers(group *echo.Group) {
	group.GET("/playlists", h.HandleGetPlaylists)
	group.POST("/playlists", h.HandlePostPlaylist)
	group.GET("/playlists/:id", h.HandleGetPlaylistById)
	group.POST("/playlists/:id/items", h.HandlePostPlaylistItemsById)
	group.DELETE("/playlists/:id/items", h.HandleDeletePlaylistItemsById)
	group.POST("/playlists/:id/items/move", h.HandlePostPlaylistsItemsMoveById)
}
