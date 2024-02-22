package server

import (
	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/types"
)

func New(db *database.Database, workDir types.WorkDir) *echo.Echo {
	e := echo.New()

	apiConfig := handlers.New(db, workDir)

	_ = apiConfig

	apiGroup := e.Group("/api/v1")

	handlers.InstallArtistHandlers(apiGroup, apiConfig)
	// handlers.InstallAlbumHandlers(router, &apiConfig)
	// handlers.InstallTrackHandlers(router, &apiConfig)
	// handlers.InstallSyncHandlers(router, &apiConfig)
	//
	// router.Post("/queue/album", apiConfig.HandlerCreateQueueFromAlbum)

	return e
}
