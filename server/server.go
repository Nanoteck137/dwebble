package server

import (
	"github.com/MadAppGang/httplog/echolog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nanoteck137/dwebble/apis"
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
		log.Debug("Registering", "method", h.Method, "name", h.Name, "path", g.Prefix+h.Path)
		g.Group.Add(h.Method, h.Path, h.HandlerFunc, h.Middlewares...)
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

	g := NewEchoGroup(app, e, "/api/v1")
	apis.InstallArtistHandlers(app, g)
	apis.InstallAlbumHandlers(app, g)
	apis.InstallTrackHandlers(app, g)
	h.InstallSyncHandlers(g)
	h.InstallQueueHandlers(g)
	h.InstallTagHandlers(g)
	h.InstallAuthHandlers(g)
	h.InstallPlaylistHandlers(g)

	g = NewEchoGroup(app, e, "/api/v1")
	apis.InstallSystemHandlers(app, g)

	app.DB().Invalidate()

	return e
}

func ServerRoutes(app core.App) []Route {
	var h handlers.Handlers

	g := NewRouteGroup("/api/v1")
	apis.InstallArtistHandlers(app, g)
	apis.InstallAlbumHandlers(app, g)
	apis.InstallTrackHandlers(app, g)
	h.InstallSyncHandlers(g)
	h.InstallQueueHandlers(g)
	h.InstallTagHandlers(g)
	h.InstallAuthHandlers(g)
	h.InstallPlaylistHandlers(g)

	apis.InstallSystemHandlers(app, g)

	return g.Routes
}
