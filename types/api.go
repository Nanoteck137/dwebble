package types

import (
	"net/http"
)

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

type ApiArtistListItem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type ApiArtistList struct {
	Artists []ApiArtistListItem `json:"artists"`
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
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
}

type ApiGetArtistAlbumsData struct {
	Albums []ApiAlbum `json:"albums"`
}

type ApiGetAlbums struct {
	Albums []ApiAlbum `json:"albums"`
}

type ApiGetAlbumTracksData struct {
	Tracks []ApiTrack `json:"tracks"`
}
