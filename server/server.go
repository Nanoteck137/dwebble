package server

import (
	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

func New(db *database.Database, libraryDir string, workDir types.WorkDir) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		switch err := err.(type) {
		case *types.ApiError:
			c.JSON(err.Code, types.ApiResponse{
				Status: types.StatusError,
				Error:  err,
			})
		default:
			c.JSON(500, types.ApiResponse{
				Status: types.StatusError,
				Error: &types.ApiError{
					Code:    500,
					Message: err.Error(),
				},
			})
		}

	}

	e.Use(echolog.LoggerWithName("Dwebble"))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/tracks/mobile", workDir.MobileTracksDir())
	e.Static("/tracks/original", workDir.OriginalTracksDir())
	e.Static("/images", workDir.ImagesDir())

	h := handlers.New(db, libraryDir, workDir)
	apiGroup := e.Group("/api/v1")

	h.InstallArtistHandlers(apiGroup)
	h.InstallAlbumHandlers(apiGroup)
	h.InstallTrackHandlers(apiGroup)
	h.InstallSyncHandlers(apiGroup)
	h.InstallQueueHandlers(apiGroup)
	h.InstallTagHandlers(apiGroup)
	h.InstallAuthHandlers(apiGroup)

	return e
}
