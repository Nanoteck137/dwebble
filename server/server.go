package server

import (
	"net/http"

	"github.com/MadAppGang/httplog/echolog"
	"github.com/kr/pretty"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/database"
	"github.com/nanoteck137/dwebble/handlers"
	"github.com/nanoteck137/dwebble/log"
	"github.com/nanoteck137/dwebble/types"
)

type Route struct {
	Name   string
	Path   string
	Method string
	Data   any
	Body   any
}

type RouteManager struct {
	Routes []Route
}

func NewRouteManager() *RouteManager {
	return &RouteManager{
		Routes: []Route{},
	}
}

func (r *RouteManager) AddRoute(name, path, method string, data, body any) {
	r.Routes = append(r.Routes, Route{
		Name:   name,
		Path:   path,
		Method: method,
		Data:   data,
		Body:   body,
	})
}

type EchoGroup struct {
	Prefix string
	Group  *echo.Group

	routeManager *RouteManager
}

func (g *EchoGroup) GET(name string, path string, f echo.HandlerFunc, data, body any) {
	log.Debug("Registering GET", "name", name, "path", g.Prefix+path)
	g.Group.GET(path, f)

	g.routeManager.AddRoute(name, g.Prefix+path, http.MethodGet, data, body)
}

func NewGroup(e *echo.Echo, r *RouteManager, prefix string) *EchoGroup {
	g := e.Group(prefix)

	return &EchoGroup{
		Prefix:       prefix,
		Group:        g,
		routeManager: r,
	}
}

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

	routeManager := NewRouteManager()

	g := NewGroup(e, routeManager, "/api/v1")
	h.InstallArtistHandlers(g)
	h.InstallAlbumHandlers(apiGroup)
	h.InstallTrackHandlers(apiGroup)
	h.InstallSyncHandlers(apiGroup)
	h.InstallQueueHandlers(apiGroup)
	h.InstallTagHandlers(apiGroup)
	h.InstallAuthHandlers(apiGroup)
	h.InstallPlaylistHandlers(apiGroup)

	apiGroup = e.Group("/api/v1")
	h.InstallSystemHandlers(apiGroup)

	pretty.Println(routeManager)

	err := handlers.InitializeConfig(db)
	if err != nil {
		// TODO(patrik): Remove?
		log.Fatal("Failed to initialize config", "err", err)
	}

	db.Invalidate()

	return e
}
