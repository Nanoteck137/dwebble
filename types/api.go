package types

import (
	"net/http"
)

var (
	ErrNoArtist = NewApiError(http.StatusNotFound, "Artist not found")
	ErrNoAlbum  = NewApiError(http.StatusNotFound, "Album not found")
	ErrNoTrack  = NewApiError(http.StatusNotFound, "Track not found")
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

type Res struct {
	Status string    `json:"status"`
	Data   any       `json:"data,omitempty"`
	Error  *ApiError `json:"error,omitempty"`
}

func NewSuccessRes(data any) Res {
	return Res{
		Status: StatusSuccess,
		Data:   data,
	}
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

type GetArtistsItem struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GetArtists struct {
	Artists []GetArtistsItem `json:"artists"`
}

type GetArtistById struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type GetArtistAlbumsByIdItem struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type GetArtistAlbumsById struct {
	Albums []GetArtistAlbumsByIdItem `json:"albums"`
}

type GetAlbumsItem struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type GetAlbums struct {
	Albums []GetAlbumsItem `json:"albums"`
}

type GetAlbumById struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type GetAlbumTracksByIdItem struct {
	Id                string `json:"id"`
	Number            int    `json:"number"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId           string `json:"albumId"`
	ArtistId          string `json:"artistId"`
}

type GetAlbumTracksById struct {
	Tracks []GetAlbumTracksByIdItem `json:"tracks"`
}

type GetTracksItem struct {
	Id                string `json:"id"`
	Number            int    `json:"number"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId           string `json:"albumId"`
	ArtistId          string `json:"artistId"`
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
}

type GetTracks struct {
	Tracks []GetTracksItem `json:"tracks"`
}

type GetTrackById struct {
	Id                string `json:"id"`
	Number            int    `json:"number"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId           string `json:"albumId"`
	ArtistId          string `json:"artistId"`
}

type GetSync struct {
	IsSyncing bool `json:"isSyncing"`
}

// NOTE(patrik): Old

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
