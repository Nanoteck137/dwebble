package apis

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
)

type systemApi struct {
	app core.App
}

func (api *systemApi) HandleGetSystemInfo(c echo.Context) error {
	return c.JSON(200, SuccessResponse(types.GetSystemInfo{
		Version: dwebble.Version,
		IsSetup: api.app.IsSetup(),
	}))
}

func (api *systemApi) HandlePostSystemExport(c echo.Context) error {
	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	if user.Id != api.app.DBConfig().OwnerId {
		// TODO(patrik): Fix error
		return errors.New("Only the owner can export")
	}

	users, err := api.app.DB().GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	res := types.PostSystemExport{
		Users: []types.ExportUser{},
	}

	for _, user := range users {
		playlists, err := api.app.DB().GetPlaylistsByUser(c.Request().Context(), user.Id)
		if err != nil {
			return err
		}

		var exportedPlaylists []types.ExportPlaylist

		for _, playlist := range playlists {
			items, err := api.app.DB().GetPlaylistItems(c.Request().Context(), playlist.Id)
			if err != nil {
				return err
			}

			playlistTracks := make([]types.ExportTrack, 0, len(items))

			for _, item := range items {
				track, err := api.app.DB().GetTrackById(c.Request().Context(), item.TrackId)
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

	return c.JSON(200, SuccessResponse(res))
}

func (api *systemApi) HandlePostSystemImport(c echo.Context) error {
	user, err := User(api.app, c)
	if err != nil {
		return err
	}

	if user.Id != api.app.DBConfig().OwnerId {
		// TODO(patrik): Fix error
		return errors.New("Only the owner can export")
	}

	// TODO(patrik): Fix error
	return errors.New("Import is not supported right now")
}

func InstallSystemHandlers(app core.App, group Group) {
	api := systemApi{app: app}

	group.Register(
		Handler{
			Name:        "GetSystemInfo",
			Path:        "/system/info",
			Method:      http.MethodGet,
			DataType:    types.GetSystemInfo{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetSystemInfo,
		},

		Handler{
			Name:        "SystemExport",
			Path:        "/system/export",
			Method:      http.MethodPost,
			DataType:    types.PostSystemExport{},
			BodyType:    nil,
			HandlerFunc: api.HandlePostSystemExport,
			Middlewares: []echo.MiddlewareFunc{},
		},

		Handler{
			Name:        "SystemImport",
			Path:        "/system/import",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: api.HandlePostSystemImport,
			Middlewares: []echo.MiddlewareFunc{},
		},
	)
}
