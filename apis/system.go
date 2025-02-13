package apis

import (
	"context"
	"net/http"

	"github.com/nanoteck137/dwebble"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/pyrin"
)

type GetSystemInfo struct {
	Version string `json:"version"`
}

func InstallSystemHandlers(app core.App, group pyrin.Group) {
	group.Register(
		pyrin.ApiHandler{
			Name:         "GetSystemInfo",
			Path:         "/system/info",
			Method:       http.MethodGet,
			ResponseType: GetSystemInfo{},
			HandlerFunc: func(c pyrin.Context) (any, error) {
				return GetSystemInfo{
					Version: dwebble.Version,
				}, nil
			},
		},

		pyrin.ApiHandler{
			Name:   "RefillSearch",
			Path:   "/system/search",
			Method: http.MethodPost,
			HandlerFunc: func(c pyrin.Context) (any, error) {
				_, err := User(app, c, RequireAdmin)
				if err != nil {
					return nil, err
				}

				ctx := context.TODO()
				err = app.DB().RefillSearchTables(ctx)
				if err != nil {
					return nil, err
				}

				return nil, nil
			},
		},
	)
}
