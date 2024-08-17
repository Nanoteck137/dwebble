package apis

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

type playlistApi struct {
	app core.App
}

func (api *playlistApi) HandleGetPlaylists(c echo.Context) error {
	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	playlists, err := api.app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
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

	return c.JSON(200, SuccessResponse(res))
}

func (api *playlistApi) HandlePostPlaylist(c echo.Context) error {
	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistBody](c)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
		Name:    body.Name,
		OwnerId: user.Id,
	})
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(types.PostPlaylist{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
	}))
}

func (api *playlistApi) HandleGetPlaylistById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return errors.New("No playlist")
	}

	items, err := api.app.DB().GetPlaylistItems(c.Request().Context(), playlist.Id)
	if err != nil {
		return err
	}

	tracks := []types.Track{}
	for _, item := range items {
		// TODO(patrik): Optimize
		track, err := api.app.DB().GetTrackById(c.Request().Context(), item.TrackId)
		if err != nil {
			return err
		}

		t := ConvertDBTrack(c, track) 
		t.Number = item.ItemIndex

		tracks = append(tracks, t)
	}

	return c.JSON(200, SuccessResponse(types.GetPlaylistById{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
		Items: tracks,
	}))
}

func (api *playlistApi) HandlePostPlaylistItemsById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistItemsByIdBody](c)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return errors.New("No playlist")
	}

	err = api.app.DB().AddItemsToPlaylist(c.Request().Context(), playlist.Id, body.Tracks)
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(nil))
}

func (api *playlistApi) HandleDeletePlaylistItemsById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	body, err := Body[types.DeletePlaylistItemsByIdBody](c)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return errors.New("No playlist")
	}

	err = api.app.DB().DeleteItemsFromPlaylist(c.Request().Context(), playlist.Id, body.TrackIndices)
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(nil))
}

func (api *playlistApi) HandlePostPlaylistsItemsMoveById(c echo.Context) error {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	body, err := Body[types.PostPlaylistsItemMoveByIdBody](c)
	if err != nil {
		return err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return errors.New("No playlist")
	}

	err = api.app.DB().MovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId, body.ToIndex)
	if err != nil {
		return err
	}

	return c.JSON(200, SuccessResponse(nil))
}

func InstallPlaylistHandlers(app core.App, group Group) {
	api := playlistApi{app: app}

	group.Register(
		Handler{
			Name:        "GetPlaylists",
			Path:        "/playlists",
			Method:      http.MethodGet,
			DataType:    types.GetPlaylists{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetPlaylists,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "CreatePlaylist",
			Path:        "/playlists",
			Method:      http.MethodPost,
			DataType:    types.PostPlaylist{},
			BodyType:    types.PostPlaylistBody{},
			HandlerFunc: api.HandlePostPlaylist,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "GetPlaylistById",
			Path:        "/playlists/:id",
			Method:      http.MethodGet,
			DataType:    types.GetPlaylistById{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetPlaylistById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "AddItemsToPlaylist",
			Path:        "/playlists/:id/items",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    types.PostPlaylistItemsByIdBody{},
			HandlerFunc: api.HandlePostPlaylistItemsById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "DeletePlaylistItems",
			Path:        "/playlists/:id/items",
			Method:      http.MethodDelete,
			DataType:    nil,
			BodyType:    types.DeletePlaylistItemsByIdBody{},
			HandlerFunc: api.HandleDeletePlaylistItemsById,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "MovePlaylistItem",
			Path:        "/playlists/:id/items/move",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    types.PostPlaylistsItemMoveByIdBody{},
			HandlerFunc: api.HandlePostPlaylistsItemsMoveById,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
