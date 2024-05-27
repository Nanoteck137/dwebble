package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/types"
)

var syncing = false

func (h *Handlers) HandleGetSync(c echo.Context) error {
	return c.JSON(200, types.NewApiSuccessResponse(types.GetSync{
		IsSyncing: syncing,
	}))
}

func (h *Handlers) HandlePostSync(c echo.Context) error {
	go func() {
		syncing = true
		defer func() { syncing = false }()

		lib, err := library.ReadFromDir(h.libraryDir)
		if err != nil {
			log.Printf("Failed to sync: %v", err)
			return
		}

		err = lib.Sync(h.workDir, h.db)
		if err != nil {
			log.Printf("Failed to sync: %v", err)
			return
		}
	}()

	return c.NoContent(http.StatusNoContent);
}

func (h *Handlers) InstallSyncHandlers(group *echo.Group) {
	group.GET("/sync", h.HandleGetSync)
	group.POST("/sync", h.HandlePostSync)
}
