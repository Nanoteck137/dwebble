package server

import (
	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/dwebble/assets"
	"github.com/nanoteck137/dwebble/core"
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
	app core.App

	Prefix string
	Group  *echo.Group
}

func (g *EchoGroup) Register(handlers ...handlers.Handler) {
	for _, h := range handlers {
		if h.NewHandlerFunc == nil {
			log.Error("Need fixing", "method", h.Method, "name", h.Name, "path", g.Prefix+h.Path)
			continue
		}

		log.Debug("Registering", "method", h.Method, "name", h.Name, "path", g.Prefix+h.Path)
		fixedHandler := func(c echo.Context) error { return h.NewHandlerFunc(g.app, c) }
		g.Group.Add(h.Method, h.Path, fixedHandler, h.Middlewares...)
	}
}

func NewEchoGroup(app core.App, e *echo.Echo, prefix string, m ...echo.MiddlewareFunc) *EchoGroup {
	g := e.Group(prefix, m...)

	return &EchoGroup{
		app:    app,
		Prefix: prefix,
		Group:  g,
	}
}

func New(app core.App) *echo.Echo {
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

	e.Static("/tracks/mobile", app.WorkDir().MobileTracksDir())
	e.Static("/tracks/original", app.WorkDir().OriginalTracksDir())
	e.StaticFS("/images/default", assets.AssetsFS)
	e.Static("/images", app.WorkDir().ImagesDir())

	h := handlers.New(app.DB(), app.Config().LibraryDir, app.WorkDir())

	isSetupMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !handlers.IsSetup() {
				return types.NewApiError(400, "Server not setup")
			}

			return next(c)
		}
	}

	g := NewEchoGroup(app, e, "/api/v1", isSetupMiddleware)
	h.InstallArtistHandlers(g)
	h.InstallAlbumHandlers(g)
	h.InstallTrackHandlers(g)
	h.InstallSyncHandlers(g)
	h.InstallQueueHandlers(g)
	h.InstallTagHandlers(g)
	h.InstallAuthHandlers(g)
	h.InstallPlaylistHandlers(g)

	g = NewEchoGroup(app, e, "/api/v1")
	h.InstallSystemHandlers(g)

	err := handlers.InitializeConfig(app.DB())
	if err != nil {
		// TODO(patrik): Remove?
		log.Fatal("Failed to initialize config", "err", err)
	}

	app.DB().Invalidate()

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
