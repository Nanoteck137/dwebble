package server

import (
	"github.com/MadAppGang/httplog/echolog"
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

type RouteGroup struct {
	Prefix string
	Routes []Route
}

func NewRouteGroup(prefix string) *RouteGroup {
	return &RouteGroup{
		Prefix: prefix,
		Routes: []Route{},
	}
}

func (r *RouteGroup) AddRoute(name, path, method string, data, body any) {
	r.Routes = append(r.Routes, Route{
		Name:   name,
		Path:   path,
		Method: method,
		Data:   data,
		Body:   body,
	})
}

func (r *RouteGroup) Register(handlers ...handlers.Handler) {
	for _, h := range handlers {
		r.AddRoute(h.Name, r.Prefix+h.Path, h.Method, h.DataType, h.BodyType)
	}
}

type EchoGroup struct {
	Prefix string
	Group  *echo.Group
}

func (g *EchoGroup) Register(handlers ...handlers.Handler) {
	for _, h := range handlers {
		log.Debug("Registering", "method", h.Method, "name", h.Name, "path", g.Prefix+h.Path)
		g.Group.Add(h.Method, h.Path, h.HandlerFunc, h.Middlewares...)
	}
}

func NewEchoGroup(e *echo.Echo, prefix string, m ...echo.MiddlewareFunc) *EchoGroup {
	g := e.Group(prefix, m...)

	return &EchoGroup{
		Prefix: prefix,
		Group:  g,
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

	isSetupMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !handlers.IsSetup() {
				return types.NewApiError(400, "Server not setup")
			}

			return next(c)
		}
	}

	g := NewEchoGroup(e, "/api/v1", isSetupMiddleware)
	h.InstallArtistHandlers(g)
	h.InstallAlbumHandlers(g)
	h.InstallTrackHandlers(g)
	h.InstallSyncHandlers(g)
	h.InstallQueueHandlers(g)
	h.InstallTagHandlers(g)
	h.InstallAuthHandlers(g)
	h.InstallPlaylistHandlers(g)

	g = NewEchoGroup(e, "/api/v1")
	h.InstallSystemHandlers(g)

	err := handlers.InitializeConfig(db)
	if err != nil {
		// TODO(patrik): Remove?
		log.Fatal("Failed to initialize config", "err", err)
	}

	db.Invalidate()

	return e
}

func ServerRoutes() []Route {
	var h handlers.Handlers

	g := NewRouteGroup("/api/v1")
	h.InstallArtistHandlers(g)
	h.InstallAlbumHandlers(g)
	h.InstallTrackHandlers(g)
	h.InstallSyncHandlers(g)
	h.InstallQueueHandlers(g)
	h.InstallTagHandlers(g)
	h.InstallAuthHandlers(g)
	h.InstallPlaylistHandlers(g)

	h.InstallSystemHandlers(g)

	return g.Routes
}
