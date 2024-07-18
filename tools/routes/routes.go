package routes

import (
	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
)

type Route struct {
	Name   string
	Path   string
	Method string
	Data   any
	Body   types.Body
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

func (r *RouteGroup) AddRoute(name, path, method string, data any, body types.Body) {
	r.Routes = append(r.Routes, Route{
		Name:   name,
		Path:   path,
		Method: method,
		Data:   data,
		Body:   body,
	})
}

func (r *RouteGroup) Register(handlers ...apis.Handler) {
	for _, h := range handlers {
		r.AddRoute(h.Name, r.Prefix+h.Path, h.Method, h.DataType, h.BodyType)
	}
}

func ServerRoutes(app core.App) []Route {
	g := NewRouteGroup("/api/v1")
	apis.InstallHandlers(app, g)

	return g.Routes
}
