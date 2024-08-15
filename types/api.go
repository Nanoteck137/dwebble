package types

import (
	"net/http"
)

var (
	ErrNoArtist   = NewApiError(http.StatusNotFound, "Artist not found")
	ErrNoAlbum    = NewApiError(http.StatusNotFound, "Album not found")
	ErrNoTrack    = NewApiError(http.StatusNotFound, "Track not found")
	ErrNoPlaylist = NewApiError(http.StatusNotFound, "Playlist not found")
	ErrNoTag      = NewApiError(http.StatusNotFound, "Tag not found")
	ErrNoGenre    = NewApiError(http.StatusNotFound, "Genre not found")

	ErrInvalidAuthHeader = NewApiError(http.StatusUnauthorized, "Invalid Auth Header")
	ErrInvalidToken      = NewApiError(http.StatusUnauthorized, "Invalid Token")
	ErrIncorrectCreds    = NewApiError(http.StatusUnauthorized, "Incorrect credentials")

	ErrEmptyBody = NewApiError(http.StatusBadRequest, "Expected body not to be empty")
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (err *ApiError) Error() string {
	return err.Message
}

func NewApiError(code int, message string, errors ...any) *ApiError {
	var e any
	if len(errors) > 0 {
		e = errors[0]
	}

	return &ApiError{
		Code:    code,
		Message: message,
		Errors:  e,
	}
}

type ApiResponse struct {
	Status string    `json:"status"`
	Data   any       `json:"data,omitempty"`
	Error  *ApiError `json:"error,omitempty"`
}

func NewApiSuccessResponse(data any) ApiResponse {
	return ApiResponse{
		Status: StatusSuccess,
		Data:   data,
	}
}
