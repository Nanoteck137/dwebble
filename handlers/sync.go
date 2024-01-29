package handlers

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/types"
)

var syncing = false

// HandleGetSync godoc
//
//	@Summary		Get sync status
//	@Description	Get sync status
//	@Tags			sync
//	@Produce		json
//	@Success		200	{object}	types.ApiResponse[types.ApiGetSyncData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/sync [get]
func (api *ApiConfig) HandleGetSync(c *fiber.Ctx) error {
	return c.JSON(types.NewApiResponse(types.ApiGetSyncData{
		Syncing: syncing,
	}))
}

// HandlePostSync godoc
//
//	@Summary		Sync Library
//	@Description	Sync Library
//	@Tags			sync
//	@Produce		json
//	@Success		200	{object}	types.ApiResponse[types.ApiPostSyncData]
//	@Failure		400	{object}	types.ApiError
//	@Failure		500	{object}	types.ApiError
//	@Router			/sync [post]
func (api *ApiConfig) HandlePostSync(c *fiber.Ctx) error {
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

func InstallSyncHandlers(router *fiber.App, apiConfig *ApiConfig) {
	router.Get("/sync", apiConfig.HandleGetSync)
	router.Post("/sync", apiConfig.HandlePostSync)
}
