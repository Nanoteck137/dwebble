package apis

import (
	"context"
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
}

type GetPlaylistById struct {
	Playlist

	Items []Track `json:"items"`
}

// TODO(patrik): Validation and transform
type AddItemToPlaylistBody struct {
	TrackId string `json:"trackId"`
}

// TODO(patrik): Validation and transform
type RemovePlaylistItemBody struct {
	TrackId string `json:"trackId"`
}

// TODO(patrik): Validation and transform
type MovePlaylistItemBody struct {
	TrackId string `json:"trackId"`
	ToIndex int    `json:"toIndex"`
}

func InstallPlaylistHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:       "GetPlaylists",
			Path:       "/playlists",
			Method:     http.MethodGet,
			ReturnType: GetPlaylists{},
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
			Name:       "CreatePlaylist",
			Path:       "/playlists",
			Method:     http.MethodPost,
			ReturnType: CreatePlaylist{},
			BodyType:   CreatePlaylistBody{},
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
			Name:       "CreatePlaylistFromFilter",
			Path:       "/playlists/filter",
			Method:     http.MethodPost,
			ReturnType: CreatePlaylist{},
			BodyType:   PostPlaylistFilterBody{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[PostPlaylistFilterBody](c)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()

				db, tx, err := app.DB().Begin()
				if err != nil {
					return nil, err
				}
				defer tx.Rollback()

				playlist, err := db.CreatePlaylist(ctx, database.CreatePlaylistParams{
					Name:    body.Name,
					OwnerId: user.Id,
				})
				if err != nil {
					return nil, err
				}

				tracks, err := db.GetAllTracks(ctx, body.Filter, "")
				if err != nil {
					if errors.Is(err, database.ErrInvalidFilter) {
						return nil, InvalidFilter(err)
					}

					return nil, err
				}

				for _, track := range tracks {
					err = db.AddItemToPlaylist(ctx, playlist.Id, track.Id)
					if err != nil {
						return nil, err
					}
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
			Name:       "GetPlaylistById",
			Path:       "/playlists/:id",
			Method:     http.MethodGet,
			ReturnType: GetPlaylistById{},
			Errors:     []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
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
			Name:     "AddItemToPlaylist",
			Path:     "/playlists/:id/items",
			Method:   http.MethodPost,
			BodyType: AddItemToPlaylistBody{},
			Errors:   []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[AddItemToPlaylistBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = app.DB().AddItemToPlaylist(c.Request().Context(), playlist.Id, body.TrackId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},

		pyrin.ApiHandler{
			Name:     "RemovePlaylistItem",
			Path:     "/playlists/:id/items",
			Method:   http.MethodDelete,
			BodyType: RemovePlaylistItemBody{},
			Errors:   []pyrin.ErrorType{ErrTypePlaylistNotFound},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				playlistId := c.Param("id")

				user, err := User(app, c)
				if err != nil {
					return nil, err
				}

				body, err := pyrin.Body[RemovePlaylistItemBody](c)
				if err != nil {
					return nil, err
				}

				playlist, err := app.DB().GetPlaylistById(c.Request().Context(), playlistId)
				if err != nil {
					if errors.Is(err, database.ErrItemNotFound) {
						return nil, PlaylistNotFound()
					}

					return nil, err
				}

				if playlist.OwnerId != user.Id {
					return nil, PlaylistNotFound()
				}

				err = app.DB().RemovePlaylistItem(c.Request().Context(), playlist.Id, body.TrackId)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
