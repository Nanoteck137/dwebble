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

// type GetArtistsItem struct {
// 	Id      string `json:"id"`
// 	Name    string `json:"name"`
// 	Picture string `json:"picture"`
// }
//
// type GetArtists struct {
// 	Artists []GetArtistsItem `json:"artists"`
// }
//
// type GetArtistById struct {
// 	Id      string `json:"id"`
// 	Name    string `json:"name"`
// 	Picture string `json:"picture"`
// }
//
// type GetArtistAlbumsByIdItem struct {
// 	Id       string `json:"id"`
// 	Name     string `json:"name"`
// 	CoverArt string `json:"coverArt"`
// 	ArtistId string `json:"artistId"`
// }
//
// type GetArtistAlbumsById struct {
// 	Albums []GetArtistAlbumsByIdItem `json:"albums"`
// }
//
// type GetAlbumsItem struct {
// 	Id       string `json:"id"`
// 	Name     string `json:"name"`
// 	CoverArt string `json:"coverArt"`
// 	ArtistId string `json:"artistId"`
// }
//
// type GetAlbums struct {
// 	Albums []GetAlbumsItem `json:"albums"`
// }
//
// type GetAlbumById struct {
// 	Id       string `json:"id"`
// 	Name     string `json:"name"`
// 	CoverArt string `json:"coverArt"`
// 	ArtistId string `json:"artistId"`
// }
//
// type GetAlbumTracksByIdItem struct {
// 	Id                string `json:"id"`
// 	Number            int    `json:"number"`
// 	Name              string `json:"name"`
// 	CoverArt          string `json:"coverArt"`
// 	Duration          int    `json:"duration"`
// 	BestQualityFile   string `json:"bestQualityFile"`
// 	MobileQualityFile string `json:"mobileQualityFile"`
// 	AlbumId           string `json:"albumId"`
// 	ArtistId          string `json:"artistId"`
// 	AlbumName         string `json:"albumName"`
// 	ArtistName        string `json:"artistName"`
// }
//
// type GetAlbumTracksById struct {
// 	Tracks []GetAlbumTracksByIdItem `json:"tracks"`
// }
//
// type GetTracksItem struct {
// 	Id                string `json:"id"`
// 	Number            int    `json:"number"`
// 	Name              string `json:"name"`
// 	CoverArt          string `json:"coverArt"`
// 	Duration          int    `json:"duration"`
// 	BestQualityFile   string `json:"bestQualityFile"`
// 	MobileQualityFile string `json:"mobileQualityFile"`
// 	AlbumId           string `json:"albumId"`
// 	ArtistId          string `json:"artistId"`
// 	AlbumName         string `json:"albumName"`
// 	ArtistName        string `json:"artistName"`
// }
//
// type GetTracks struct {
// 	Tracks []GetTracksItem `json:"tracks"`
// }
//
// type GetTrackById struct {
// 	Id                string `json:"id"`
// 	Number            int    `json:"number"`
// 	Name              string `json:"name"`
// 	CoverArt          string `json:"coverArt"`
// 	Duration          int    `json:"duration"`
// 	BestQualityFile   string `json:"bestQualityFile"`
// 	MobileQualityFile string `json:"mobileQualityFile"`
// 	AlbumId           string `json:"albumId"`
// 	ArtistId          string `json:"artistId"`
// }
//
// type GetSync struct {
// 	IsSyncing bool `json:"isSyncing"`
// }
