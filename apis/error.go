package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin/api"
)

const (
	ErrTypeArtistNotFound api.ErrorType = "ARTIST_NOT_FOUND"
	ErrTypeAlbumNotFound  api.ErrorType = "ALBUM_NOT_FOUND"
	ErrTypeTrackNotFound  api.ErrorType = "TRACK_NOT_FOUND"
	ErrTypeRouteNotFound  api.ErrorType = "ROUTE_NOT_FOUND"
	ErrTypeInvalidFilter  api.ErrorType = "INVALID_FILTER"
	ErrTypeSortFilter     api.ErrorType = "INVALID_SORT"
)

func ArtistNotFound() *api.Error {
	return &api.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeArtistNotFound,
		Message: "Artist not found",
	}
}

func AlbumNotFound() *api.Error {
	return &api.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeAlbumNotFound,
		Message: "Album not found",
	}
}

func TrackNotFound() *api.Error {
	return &api.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeTrackNotFound,
		Message: "Track not found",
	}
}

func InvalidFilter(err error) *api.Error {
	return &api.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidFilter,
		Message: err.Error(),
	}
}

func InvalidSort(err error) *api.Error {
	return &api.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeSortFilter,
		Message: err.Error(),
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
	return api.SuccessResponse(data)
}
