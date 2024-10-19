package apis

import (
	"net/http"

	"github.com/nanoteck137/pyrin"
)

const (
	ErrTypeArtistNotFound pyrin.ErrorType = "ARTIST_NOT_FOUND"
	ErrTypeAlbumNotFound  pyrin.ErrorType = "ALBUM_NOT_FOUND"
	ErrTypeTrackNotFound  pyrin.ErrorType = "TRACK_NOT_FOUND"
	ErrTypeInvalidFilter  pyrin.ErrorType = "INVALID_FILTER"
	ErrTypeInvalidSort    pyrin.ErrorType = "INVALID_SORT"
)

func ArtistNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeArtistNotFound,
		Message: "Artist not found",
	}
}

func AlbumNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeAlbumNotFound,
		Message: "Album not found",
	}
}

func TrackNotFound() *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusNotFound,
		Type:    ErrTypeTrackNotFound,
		Message: "Track not found",
	}
}

func InvalidFilter(err error) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidFilter,
		Message: err.Error(),
	}
}

func InvalidSort(err error) *pyrin.Error {
	return &pyrin.Error{
		Code:    http.StatusBadRequest,
		Type:    ErrTypeInvalidSort,
		Message: err.Error(),
	}
}
