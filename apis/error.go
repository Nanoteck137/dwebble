package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin/api"
)

const (
	ErrTypeTrackNotFound api.ErrorType = "TRACK_NOT_FOUND"
	ErrTypeRouteNotFound api.ErrorType = "ROUTE_NOT_FOUND"
)

func TrackNotFound() *api.Error {
	return &api.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeTrackNotFound,
		Message: "Track not found",
	}
}

func RouteNotFound() *api.Error {
	return &api.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeRouteNotFound,
		Message: "Route not found",
	}
}

func SuccessResponse(data any) api.Response {
	return api.Response{
		Success: true,
		Data:    data,
	}
}
