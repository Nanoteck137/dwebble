package apis

import (
	"errors"
	"net/http"

	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/pyrin"
)

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// TODO(patrik): Add all fields
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type CreatePlaylist struct {
	Playlist
}

// TODO(patrik): Validation and transform
type CreatePlaylistBody struct {
	Name string `json:"name"`
}

// TODO(patrik): Validation and transform
type PostPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
}

type GetPlaylistById struct {
	Playlist

	Items []Track `json:"items"`
}

// TODO(patrik): Validation and transform
type AddItemsToPlaylistBody struct {
	Tracks []string `json:"tracks"`
}

// TODO(patrik): Validation and transform
type DeletePlaylistItemsBody struct {
	TrackIds []string `json:"trackIds"`
}

// TODO(patrik): Validation and transform
type MovePlaylistItemBody struct {
	TrackId string `json:"trackId"`
	ToIndex int    `json:"toIndex"`
}

func InstallPlaylistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:     "GetPlaylists",
			Path:     "/playlists",
			Method:   http.MethodGet,
			DataType: GetPlaylists{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlists, err := app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
				if err != nil {
					return nil, err
				}

				res := GetPlaylists{
					Playlists: make([]Playlist, len(playlists)),
				}

				for i, playlist := range playlists {
					res.Playlists[i] = Playlist{
						Id:   playlist.Id,
						Name: playlist.Name,
					}
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "CreatePlaylist",
			Path:     "/playlists",
			Method:   http.MethodPost,
			DataType: CreatePlaylist{},
			BodyType: CreatePlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[CreatePlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().CreatePlaylist(c.Request().Context(), database.CreatePlaylistParams{
					Name:    body.Name,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				return CreatePlaylist{
					Playlist: Playlist{
						Id:   playlist.Id,
						Name: playlist.Name,
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "CreatePlaylistFromFilter",
			Path:     "/playlists/filter",
			Method:   http.MethodPost,
			DataType: CreatePlaylist{},
			BodyType: PostPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[PostPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				db, tx, err := app.DB().Begin()
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

				tracks, err := db.GetAllTracks(c.Request().Context(), body.Filter, body.Sort)
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

				return CreatePlaylist{
					Playlist: Playlist{
						Id:   playlist.Id,
						Name: playlist.Name,
					},
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "GetPlaylistById",
			Path:     "/playlists/:id",
			Method:   http.MethodGet,
			DataType: GetPlaylistById{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					return nil, err
				}

				if playlist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No playlist")
				}

				tracks, err := app.DB().GetPlaylistTracks(c.Request().Context(), playlist.Id)
				if err != nil {
					return nil, err
				}

				res := GetPlaylistById{
					Playlist: Playlist{
						Id:   playlist.Id,
						Name: playlist.Name,
					},
					Items: make([]Track, len(tracks)),
				}

				for i, track := range tracks {
					res.Items[i] = ConvertDBTrack(c, track)
				}

				return res, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "AddItemsToPlaylist",
			Path:     "/playlists/:id/items",
			Method:   http.MethodPost,
			BodyType: AddItemsToPlaylistBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddItemsToPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					return nil, err
				}

				if playlist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No playlist")
				}

				err = app.DB().AddItemsToPlaylist(c.Request().Context(), playlist.Id, body.Tracks)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "DeletePlaylistItems",
			Path:     "/playlists/:id/items",
			Method:   http.MethodDelete,
			BodyType: DeletePlaylistItemsBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[DeletePlaylistItemsBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					return nil, err
				}

				if playlist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No playlist")
				}

				err = app.DB().DeleteItemsFromPlaylist(c.Request().Context(), playlist.Id, body.TrackIds)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "MovePlaylistItem",
			Path:     "/playlists/:id/items/move",
			Method:   http.MethodPost,
			BodyType: MovePlaylistItemBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[MovePlaylistItemBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					return nil, err
				}

				if playlist.OwnerId != user.Id {
					// TODO(patrik): Fix error
					return nil, errors.New("No playlist")
				}

				err = app.DB().MovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId, body.ToIndex)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
