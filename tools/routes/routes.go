package routes

import (
	"github.com/nanoteck137/dwebble/apis"
	"github.com/nanoteck137/dwebble/core"
	"github.com/nanoteck137/dwebble/types"
	"github.com/nanoteck137/pyrin/api"
)

type Route struct {
	Name       string
	Path       string
	Method     string
	ErrorTypes []api.ErrorType
	Data       any
	Body       types.Body
	IsMultiForm bool
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

func (r *RouteGroup) AddRoute(name, path, method string, errorTypes []api.ErrorType, data any, body types.Body, isMultiForm bool) {
	r.Routes = append(r.Routes, Route{
		Name:        name,
		Path:        path,
		Method:      method,
		ErrorTypes:  errorTypes,
		Data:        data,
		Body:        body,
		IsMultiForm: isMultiForm,
	})
}

func (r *RouteGroup) Register(handlers ...apis.Handler) {
	for _, h := range handlers {
		r.AddRoute(h.Name, r.Prefix+h.Path, h.Method, h.Errors, h.DataType, h.BodyType, h.IsMultiForm)
	}
}

func ServerRoutes(app core.App) []Route {
	g := NewRouteGroup("/api/v1")
	apis.InstallHandlers(app, g)

	return g.Routes
}
