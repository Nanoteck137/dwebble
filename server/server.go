package server

import (
	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

func New(db *database.Database, libraryDir string, workDir types.WorkDir) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		c.JSON(500, map[string]any{
			"message": err.Error(),
		})
	}

	e.Use(echolog.LoggerWithName("Dwebble"))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Static("/tracks/mobile", workDir.MobileTracksDir())
	e.Static("/tracks/original", workDir.OriginalTracksDir())
	e.Static("/images", workDir.ImagesDir())

	apiConfig := handlers.New(db, libraryDir, workDir)

	_ = apiConfig

	apiGroup := e.Group("/api/v1")

	handlers.InstallArtistHandlers(apiGroup, apiConfig)
	handlers.InstallAlbumHandlers(apiGroup, apiConfig)
	handlers.InstallTrackHandlers(apiGroup, apiConfig)
	handlers.InstallSyncHandlers(apiGroup, apiConfig)
	// handlers.InstallAlbumHandlers(router, &apiConfig)
	// handlers.InstallTrackHandlers(router, &apiConfig)
	// handlers.InstallSyncHandlers(router, &apiConfig)
	//
	// router.Post("/queue/album", apiConfig.HandlerCreateQueueFromAlbum)

	return e
}
