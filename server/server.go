package server

import (
	"log"

	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/dwebble/assets"
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
	e.StaticFS("/images/default", assets.AssetsFS)
	e.Static("/images", workDir.ImagesDir())

	h := handlers.New(db, libraryDir, workDir)
	apiGroup := e.Group("/api/v1", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !handlers.IsSetup() {
				return types.NewApiError(400, "Server not setup")
			}

			return next(c)
		}
	})

	h.InstallArtistHandlers(apiGroup)
	h.InstallAlbumHandlers(apiGroup)
	h.InstallTrackHandlers(apiGroup)
	h.InstallSyncHandlers(apiGroup)
	h.InstallQueueHandlers(apiGroup)
	h.InstallTagHandlers(apiGroup)
	h.InstallAuthHandlers(apiGroup)
	h.InstallPlaylistHandlers(apiGroup)

	apiGroup = e.Group("/api/v1")
	h.InstallSystemHandlers(apiGroup)

	err := handlers.InitializeConfig(db)
	if err != nil {
		log.Fatal(err)
	}

	return e
}
