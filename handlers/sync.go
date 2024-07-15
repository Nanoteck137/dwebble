package handlers

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/log"
	"github.com/nanoteck137/dwebble/types"
)

var syncing atomic.Bool

func (h *Handlers) HandleGetSync(c echo.Context) error {
	return c.JSON(200, types.NewApiSuccessResponse(types.GetSync{
		IsSyncing: syncing.Load(),
	}))
}

func (h *Handlers) HandlePostSync(c echo.Context) error {
	go func() {
		syncing.Store(true)
		defer syncing.Store(false)

		start := time.Now()

		log.Info("Sync: Reading library")
		lib, err := library.ReadFromDir(h.libraryDir)
		if err != nil {
			log.Error("Failed to sync", "err", err)
			return
		}
		log.Info("Sync: Done Reading library", "time", time.Since(start))

		start = time.Now()

		log.Info("Sync: Started Syncing library")
		err = lib.Sync(h.workDir, h.db)
		if err != nil {
			log.Error("Failed to sync", "err", err)
			return
		}
		log.Info("Sync: Done Syncing library", "time", time.Since(start))

		h.db.Invalidate()
	}()

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func (h *Handlers) InstallSyncHandlers(group Group) {
	group.Register(
		Handler{
			Name:        "GetSyncStatus",
			Path:        "/sync",
			Method:      http.MethodGet,
			DataType:    types.GetSync{},
			BodyType:    nil,
			HandlerFunc: h.HandleGetSync,
		},

		Handler{
			Name:        "RunSync",
			Path:        "/sync",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: h.HandlePostSync,
		},
	)
}
