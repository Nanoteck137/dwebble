package handlers

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/types"
)

var syncing = false

func (api *ApiConfig) HandleGetSync(c echo.Context) error {
	return c.JSON(200, types.NewApiResponse(types.ApiGetSyncData{
		Syncing: syncing,
	}))
}

func (api *ApiConfig) HandlePostSync(c echo.Context) error {
	go func() {
		syncing = true
		defer func() { syncing = false }()

		dir := "/Volumes/media/music"
		fsys := os.DirFS(dir)
		lib, err := library.ReadFromFS(fsys)
		if err != nil {
			log.Printf("Failed to sync: %v", err)
			return
		}

		err = lib.Sync(api.workDir, dir, api.db)
		if err != nil {
			log.Printf("Failed to sync: %v", err)
			return
		}
	}()

	return nil
}

func InstallSyncHandlers(group *echo.Group, apiConfig *ApiConfig) {
	group.GET("/sync", apiConfig.HandleGetSync)
	group.POST("/sync", apiConfig.HandlePostSync)
}