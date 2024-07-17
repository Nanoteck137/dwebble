package apis

import (
	"net/http"

	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/config"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

type systemApi struct {
	app core.App
}

func (api *systemApi) HandleGetSystemInfo(c echo.Context) error {
	return c.JSON(200, types.NewApiSuccessResponse(types.GetSystemInfo{
		Version: config.Version,
		IsSetup: api.app.IsSetup(),
	}))
}

func (api *systemApi) HandlePostSystemSetup(c echo.Context) error {
	if api.app.IsSetup() {
		return types.NewApiError(400, "System already setup")
	}

	body, err := handlers.Body[types.PostSystemSetupBody](c, types.PostSystemSetupBodySchema)
	if err != nil {
		return err
	}

	db, tx, err := api.app.DB().Begin()
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

	api.app.UpdateDBConfig(&conf)

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (api *systemApi) HandlePostSystemExport(c echo.Context) error {
	user, err := handlers.User(api.app.DB(), c)
	if err != nil {
		return err
	}

	if user.Id != api.app.DBConfig().OwnerId {
		return types.NewApiError(403, "Only the owner can export")
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

	return c.JSON(200, types.NewApiSuccessResponse(res))
}

func (api *systemApi) HandlePostSystemImport(c echo.Context) error {
	user, err := handlers.User(api.app.DB(), c)
	if err != nil {
		return err
	}

	if user.Id != api.app.DBConfig().OwnerId {
		return types.NewApiError(403, "Only the owner can import")
	}

	return types.NewApiError(400, "Import is not supported right now")
}

func InstallSystemHandlers(app core.App, group handlers.Group) {
	api := systemApi{app: app}

	requireSetup := RequireSetup(app)

	group.Register(
		handlers.Handler{
			Name:        "GetSystemInfo",
			Path:        "/system/info",
			Method:      http.MethodGet,
			DataType:    types.GetSystemInfo{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetSystemInfo,
		},

		handlers.Handler{
			Name:        "RunSystemSetup",
			Path:        "/system/setup",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    types.PostSystemSetupBody{},
			HandlerFunc: api.HandlePostSystemSetup,
		},

		handlers.Handler{
			Name:        "SystemExport",
			Path:        "/system/export",
			Method:      http.MethodPost,
			DataType:    types.PostSystemExport{},
			BodyType:    nil,
			HandlerFunc: api.HandlePostSystemExport,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		handlers.Handler{
			Name:        "SystemImport",
			Path:        "/system/import",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: api.HandlePostSystemImport,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},
	)
}
