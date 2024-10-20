package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin"
)

type playlistApi struct {
	app core.App
}

func (api *playlistApi) HandleGetPlaylists(c pyrin.Context) (any, error) {
	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	playlists, err := api.app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
	if err != nil {
		return nil, err
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

	return res, nil
}

func (api *playlistApi) HandlePostPlaylist(c pyrin.Context) (any, error) {
	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	body, err := Body[types.PostPlaylistBody](c)
	if err != nil {
		return nil, err
	}

	playlist, err := api.app.DB().CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
		Name:    body.Name,
		OwnerId: user.Id,
	})
	if err != nil {
		return nil, err
	}

	return types.PostPlaylist{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
	}, nil
}

func (api *playlistApi) HandlePostPlaylistFilter(c pyrin.Context) (any, error) {
	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	body, err := Body[types.PostPlaylistFilterBody](c)
	if err != nil {
		return nil, err
	}

	db, tx, err := api.app.DB().Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	playlist, err := db.CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
		Name:    body.Name,
		OwnerId: user.Id,
	})
	if err != nil {
		return nil, err
	}

	tracks, err := db.GetAllTracks(c.Request().Context(), body.Filter, body.Sort, false)
	if err != nil {
		if errors.Is(err, database.ErrInvalidFilter) {
			return nil, InvalidFilter(err)
		}

		if errors.Is(err, database.ErrInvalidSort) {
			return nil, InvalidSort(err)
		}

		return nil, err
	}

	ids := make([]string, 0, len(tracks))

	for _, track := range tracks {
		ids = append(ids, track.Id)
	}

	err = db.AddItemsToPlaylistRaw(c.Request().Context(), playlist.Id, ids)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return types.PostPlaylist{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
	}, nil
}

func (api *playlistApi) HandleGetPlaylistById(c pyrin.Context) (any, error) {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return nil, err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return nil, errors.New("No playlist")
	}

	tracks, err := api.app.DB().GetPlaylistTracks(c.Request().Context(), playlist.Id)
	if err != nil {
		return nil, err
	}

	res := types.GetPlaylistById{
		Playlist: types.Playlist{
			Id:   playlist.Id,
			Name: playlist.Name,
		},
		Items: make([]types.Track, len(tracks)),
	}

	for i, track := range tracks {
		res.Items[i] = ConvertDBTrack(c, track)
	}

	return nil, nil
}

func (api *playlistApi) HandlePostPlaylistItemsById(c pyrin.Context) (any, error) {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	body, err := Body[types.PostPlaylistItemsByIdBody](c)
	if err != nil {
		return nil, err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return nil, err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return nil, errors.New("No playlist")
	}

	err = api.app.DB().AddItemsToPlaylist(c.Request().Context(), playlist.Id, body.Tracks)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (api *playlistApi) HandleDeletePlaylistItemsById(c pyrin.Context) (any, error) {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	body, err := Body[types.DeletePlaylistItemsByIdBody](c)
	if err != nil {
		return nil, err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return nil, err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return nil, errors.New("No playlist")
	}

	err = api.app.DB().DeleteItemsFromPlaylist(c.Request().Context(), playlist.Id, body.TrackIds)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (api *playlistApi) HandlePostPlaylistsItemsMoveById(c pyrin.Context) (any, error) {
	playlistId := c.Param("id")

	user, err := User(api.app, c)
	if err != nil {
		return nil, err
	}

	body, err := Body[types.PostPlaylistsItemMoveByIdBody](c)
	if err != nil {
		return nil, err
	}

	playlist, err := api.app.DB().GetPlaylistById(c.Request().Context(), playlistId)
	if err != nil {
		return nil, err
	}

	if playlist.OwnerId != user.Id {
		// TODO(patrik): Fix error
		return nil, errors.New("No playlist")
	}

	err = api.app.DB().MovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId, body.ToIndex)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func InstallPlaylistHandlers(app core.App, group pyrin.Group) {
	api := playlistApi{app: app}

	group.Register(
		pyrin.ApiHandler{
			Name:        "GetPlaylists",
			Path:        "/playlists",
			Method:      http.MethodGet,
			DataType:    types.GetPlaylists{},
			HandlerFunc: api.HandleGetPlaylists,
		},

		pyrin.ApiHandler{
			Name:        "CreatePlaylist",
			Path:        "/playlists",
			Method:      http.MethodPost,
			DataType:    types.PostPlaylist{},
			BodyType:    types.PostPlaylistBody{},
			HandlerFunc: api.HandlePostPlaylist,
		},

		pyrin.ApiHandler{
			Name:        "CreatePlaylistFromFilter",
			Path:        "/playlists/filter",
			Method:      http.MethodPost,
			DataType:    types.PostPlaylist{},
			BodyType:    types.PostPlaylistFilterBody{},
			HandlerFunc: api.HandlePostPlaylistFilter,
		},

		pyrin.ApiHandler{
			Name:        "GetPlaylistById",
			Path:        "/playlists/:id",
			Method:      http.MethodGet,
			DataType:    types.GetPlaylistById{},
			HandlerFunc: api.HandleGetPlaylistById,
		},

		pyrin.ApiHandler{
			Name:        "AddItemsToPlaylist",
			Path:        "/playlists/:id/items",
			Method:      http.MethodPost,
			BodyType:    types.PostPlaylistItemsByIdBody{},
			HandlerFunc: api.HandlePostPlaylistItemsById,
		},

		pyrin.ApiHandler{
			Name:        "DeletePlaylistItems",
			Path:        "/playlists/:id/items",
			Method:      http.MethodDelete,
			BodyType:    types.DeletePlaylistItemsByIdBody{},
			HandlerFunc: api.HandleDeletePlaylistItemsById,
		},

		pyrin.ApiHandler{
			Name:        "MovePlaylistItem",
			Path:        "/playlists/:id/items/move",
			Method:      http.MethodPost,
			BodyType:    types.PostPlaylistsItemMoveByIdBody{},
			HandlerFunc: api.HandlePostPlaylistsItemsMoveById,
		},
	)
}
