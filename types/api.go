package types

import (
	"net/http"
)

var (
	ErrNoArtist = NewResError(http.StatusNotFound, "No artist found")
)

const (
	StatusSuccess = "success"
)

type ResError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (err *ResError) Error() string {
	return err.Message
}

func NewResError(code int, message string, errors ...any) *ResError {
	var e any 
	if len(errors) > 0 {
		e = errors[0]
	}

	return &ResError{
		Code:    code,
		Message: message,
		Errors:  e,
	}
}

type Res struct {
	Status string    `json:"status"`
	Data   any       `json:"data,omitempty"`
	Error  *ResError `json:"error,omitempty"`
}

type ApiError struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (err ApiError) Error() string {
	return err.Message
}

func NewApiError(status int, message string, data ...any) ApiError {
	var d any = nil

	if data != nil && len(data) > 0 && data[0] != nil {
		d = data[0]
	}

	return ApiError{
		Status:  status,
		Message: message,
		Data:    d,
	}
}

func ApiBadRequestError(message string, data ...any) ApiError {
	return NewApiError(http.StatusBadRequest, message, data...)
}

func ApiNotFoundError(message string, data ...any) ApiError {
	return NewApiError(http.StatusNotFound, message, data...)
}

type ApiResponse[T any] struct {
	Status int `json:"status" example:"200"`
	Data   T   `json:"data"`
}

func NewApiResponse[T any](data T) ApiResponse[T] {
	return ApiResponse[T]{
		Status: 200,
		Data:   data,
	}
}

type ApiArtist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type ApiAlbum struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type ApiTrack struct {
	Id                string `json:"id"`
	Number            int32  `json:"number"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId           string `json:"albumId"`
	ArtistId          string `json:"artistId"`
}

// NOTE(patrik): Artist Handlers

type ApiGetArtistsData struct {
	Artists []ApiArtist `json:"artists"`
}

type ApiGetArtistByIdData ApiArtist

type ApiGetArtistAlbumsByIdData struct {
	Albums []ApiAlbum `json:"albums"`
}

// NOTE(patrik): Album Handlers

type ApiGetAlbumsData struct {
	Albums []ApiAlbum `json:"albums"`
}

type ApiGetAlbumByIdData ApiAlbum

type ApiGetAlbumTracksByIdData struct {
	Tracks []ApiTrack `json:"tracks"`
}

// NOTE(patrik): Album Handlers
// TODO(patrik): Use these in the track handlers
type ApiGetTracksDataTrackItem struct {
	ApiTrack
	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`
}
type ApiGetTracksData struct {
	Tracks []ApiGetTracksDataTrackItem `json:"tracks"`
}

type ApiGetTrackByIdData ApiTrack

type ApiGetSyncData struct {
	Syncing bool `json:"syncing"`
}

type ApiPostSyncData struct{}
