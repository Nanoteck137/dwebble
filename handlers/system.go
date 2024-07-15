package handlers

import (
	"context"
	"net/http"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/types"
)

var cfg *database.Config

// TODO(patrik): Should this be here?
func IsSetup() bool {
	return cfg != nil
}

func (h *Handlers) HandleGetSystemInfo(c echo.Context) error {
	return c.JSON(200, types.NewApiSuccessResponse(types.GetSystemInfo{
		Version: config.Version,
		IsSetup: IsSetup(),
	}))
}

func (h *Handlers) HandlePostSystemSetup(c echo.Context) error {
	if IsSetup() {
		return types.NewApiError(400, "System already setup")
	}

	body, err := Body[types.PostSystemSetupBody](c, types.PostSystemSetupBodySchema)
	if err != nil {
		return err
	}

	db, tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	user, err := db.CreateUser(c.Request().Context(), body.Username, body.Password)
	if err != nil {
		return err
	}

	conf, err := db.CreateConfig(c.Request().Context(), user.Id)
	if err != nil {
		return err
	}

	pretty.Println(body)

	err = tx.Commit()
	if err != nil {
		return err
	}

	cfg = &conf

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (h *Handlers) HandlePostSystemExport(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return err
	}

	if user.Id != cfg.OwnerId {
		return types.NewApiError(403, "Only the owner can export")
	}

	users, err := h.db.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.PostSystemExport{
		Users: []types.ExportUser{},
	}

	for _, user := range users {
		playlists, err := h.db.GetPlaylistsByUser(c.Request().Context(), user.Id)
		if err != nil {
			return err
		}

		var exportedPlaylists []types.ExportPlaylist

		for _, playlist := range playlists {
			items, err := h.db.GetPlaylistItems(c.Request().Context(), playlist.Id)
			if err != nil {
				return err
			}

			playlistTracks := make([]types.ExportTrack, 0, len(items))

			for _, item := range items {
				track, err := h.db.GetTrackById(c.Request().Context(), item.TrackId)
				if err != nil {
					return err
				}

				playlistTracks = append(playlistTracks, types.ExportTrack{
					Name:   track.Name,
					Album:  track.AlbumName,
					Artist: track.ArtistName,
				})
			}

			exportedPlaylists = append(exportedPlaylists, types.ExportPlaylist{
				Name:   playlist.Name,
				Tracks: playlistTracks,
			})
		}

		res.Users = append(res.Users, types.ExportUser{
			Username:  user.Username,
			Playlists: exportedPlaylists,
		})
	}

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (h *Handlers) HandlePostSystemImport(c echo.Context) error {
	user, err := h.User(c)
	if err != nil {
		return err
	}

	if user.Id != cfg.OwnerId {
		return types.NewApiError(403, "Only the owner can import")
	}

	return types.NewApiError(400, "Import is not supported right now")
}

func (h *Handlers) InstallSystemHandlers(group Group) {
	group.Register(
		Handler {
			Name: "GetSystemInfo", 
			Path: "/system/info", 
			Method: http.MethodGet,
			DataType: types.GetSystemInfo{}, 
			BodyType: nil,
			HandlerFunc: h.HandleGetSystemInfo, 
		},

		Handler {
			Name: "RunSystemSetup", 
			Path: "/system/setup", 
			Method: http.MethodPost,
			DataType: nil, 
			BodyType: types.PostSystemSetupBody{},
			HandlerFunc: h.HandlePostSystemSetup, 
		},

		Handler {
			Name: "SystemExport", 
			Path: "/system/export", 
			Method: http.MethodPost,
			DataType: types.PostSystemExport{}, 
			BodyType: nil,
			HandlerFunc: h.HandlePostSystemExport, 
		},

		Handler{
			Name: "SystemImport", 
			Path: "/system/import", 
			Method: http.MethodPost,
			DataType: nil,
			BodyType: nil,
			HandlerFunc: h.HandlePostSystemImport, 
		},
	)
}

func InitializeConfig(db *database.Database) error {
	conf, err := db.GetConfig(context.Background())
	if err != nil {
		return err
	}

	cfg = conf

	return nil
}
