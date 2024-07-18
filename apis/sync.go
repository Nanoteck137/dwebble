package apis

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/core/log"
	"github.com/nanoteck137/dwebble/library"
	"github.com/nanoteck137/dwebble/types"
)

type syncApi struct {
	app     core.App
	syncing atomic.Bool
}

func (api *syncApi) HandleGetSync(c echo.Context) error {
	return c.JSON(200, types.NewApiSuccessResponse(types.GetSync{
		IsSyncing: api.syncing.Load(),
	}))
}

func (api *syncApi) HandlePostSync(c echo.Context) error {
	go func() {
		api.syncing.Store(true)
		defer api.syncing.Store(false)

		start := time.Now()

		log.Info("Sync: Reading library")
		lib, err := library.ReadFromDir(api.app.Config().LibraryDir)
		if err != nil {
			log.Error("Failed to sync", "err", err)
			return
		}
		log.Info("Sync: Done Reading library", "time", time.Since(start))

		start = time.Now()

		log.Info("Sync: Started Syncing library")
		err = lib.Sync(api.app.Config().WorkDir(), api.app.DB())
		if err != nil {
			log.Error("Failed to sync", "err", err)
			return
		}
		log.Info("Sync: Done Syncing library", "time", time.Since(start))

		api.app.DB().Invalidate()
	}()

	return c.JSON(200, types.NewApiSuccessResponse(nil))
}

func InstallSyncHandlers(app core.App, group Group) {
	api := syncApi{
		app:     app,
		syncing: atomic.Bool{},
	}

	requireSetup := RequireSetup(app)

	group.Register(
		Handler{
			Name:        "GetSyncStatus",
			Path:        "/sync",
			Method:      http.MethodGet,
			DataType:    types.GetSync{},
			BodyType:    nil,
			HandlerFunc: api.HandleGetSync,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},

		Handler{
			Name:        "RunSync",
			Path:        "/sync",
			Method:      http.MethodPost,
			DataType:    nil,
			BodyType:    nil,
			HandlerFunc: api.HandlePostSync,
			Middlewares: []echo.MiddlewareFunc{requireSetup},
		},
	)
}
