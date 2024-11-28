package types

import (
	"github.com/faceair/jio"
	"github.com/nanoteck137/pyrin"
	"github.com/nanoteck137/pyrin/tools/validate"
)

type Body interface {
	Schema() jio.Schema
}

type Page struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type Images struct {
	Original string `json:"original"`
	Small    string `json:"small"`
	Medium   string `json:"medium"`
	Large    string `json:"large"`
}

type Artist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture Images `json:"picture"`
	Created int64  `json:"created"`
	Updated int64  `json:"updated"`
}

type Album struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CoverArt   Images `json:"coverArt"`
	ArtistId   string `json:"artistId"`
	ArtistName string `json:"artistName"`
	Year       *int64 `json:"year"`
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
}

type Track struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	AlbumId  string `json:"albumId"`
	ArtistId string `json:"artistId"`

	Number   *int64 `json:"number"`
	Duration *int64 `json:"duration"`
	Year     *int64 `json:"year"`

	OriginalMediaUrl string `json:"originalMediaUrl"`
	MobileMediaUrl   string `json:"mobileMediaUrl"`
	CoverArt         Images `json:"coverArt"`

	AlbumName  string `json:"albumName"`
	ArtistName string `json:"artistName"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`

	Tags []string `json:"tags"`
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetArtists struct {
	Artists []Artist `json:"artists"`
}

type GetArtistById struct {
	Artist
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}

type GetAlbumById struct {
	Album
}

type GetAlbumTracksById struct {
	Tracks []Track `json:"tracks"`
}

type GetTracks struct {
	Page   Page    `json:"page"`
	Tracks []Track `json:"tracks"`
}

type GetTrackById struct {
	Track
}

type GetSync struct {
	IsSyncing bool `json:"isSyncing"`
}

type PostQueue struct {
	Tracks []Track `json:"tracks"`
}

type GetTags struct {
	Tags []Tag `json:"tags"`
}

var _ pyrin.Body = (*PostAuthSignupBody)(nil)

type PostAuthSignupBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (b PostAuthSignupBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b PostAuthSignupBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"username":        jio.String().Min(4).Required(),
		"password":        jio.String().Min(8).Required(),
		"passwordConfirm": jio.String().Min(8).Required(),
	})
}

type PostAuthSignup struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

var _ pyrin.Body = (*PostAuthSigninBody)(nil)

type PostAuthSigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (b PostAuthSigninBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b PostAuthSigninBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"username": jio.String().Required(),
		"password": jio.String().Required(),
	})
}

type PostAuthSignin struct {
	Token string `json:"token"`
}

type GetAuthMe struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	IsOwner  bool   `json:"isOwner"`
}

type PostPlaylist struct {
	Playlist
}

var _ pyrin.Body = (*PostPlaylistBody)(nil)

type PostPlaylistBody struct {
	Name string `json:"name"`
}

func (b PostPlaylistBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b PostPlaylistBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"name": jio.String().Required(),
	})
}

var _ pyrin.Body = (*PostPlaylistFilterBody)(nil)

type PostPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
}

func (b PostPlaylistFilterBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b PostPlaylistFilterBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"name":   jio.String().Required().Min(1),
		"filter": jio.String(),
		"sort":   jio.String(),
	})
}

var _ pyrin.Body = (*PostPlaylistItemsByIdBody)(nil)

type PostPlaylistItemsByIdBody struct {
	Tracks []string `json:"tracks"`
}

func (b PostPlaylistItemsByIdBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b PostPlaylistItemsByIdBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"tracks": jio.Array().Items(jio.String()).Min(1).Required(),
	})
}

type GetPlaylistById struct {
	Playlist

	Items []Track `json:"items"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

var _ pyrin.Body = (*DeletePlaylistItemsByIdBody)(nil)

type DeletePlaylistItemsByIdBody struct {
	TrackIds []string `json:"trackIds"`
}

func (b DeletePlaylistItemsByIdBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b DeletePlaylistItemsByIdBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"trackIds": jio.Array().Items(jio.String().Min(1)).Min(1).Required(),
	})
}

var _ pyrin.Body = (*PostPlaylistsItemMoveByIdBody)(nil)

type PostPlaylistsItemMoveByIdBody struct {
	TrackId string `json:"trackId"`
	ToIndex int    `json:"toIndex"`
}

func (b PostPlaylistsItemMoveByIdBody) Validate(validator validate.Validator) error {
	panic("unimplemented")
}

func (b PostPlaylistsItemMoveByIdBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"trackId": jio.String().Required(),
		"toIndex": jio.Number().Integer().Required(),
	})
}

type GetSystemInfo struct {
	Version string `json:"version"`
}

type PostSystemSetupBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

func (b PostSystemSetupBody) Schema() jio.Schema {
	return jio.Object().Keys(jio.K{
		"username":        jio.String().Required(),
		"password":        jio.String().Required(),
		"passwordConfirm": jio.String().Required(),
	})
}

type ExportTrack struct {
	Name   string `json:"name"`
	Album  string `json:"album"`
	Artist string `json:"artist"`
}

type ExportPlaylist struct {
	Name   string        `json:"name"`
	Tracks []ExportTrack `json:"tracks"`
}

type ExportUser struct {
	Username  string           `json:"username"`
	Playlists []ExportPlaylist `json:"playlists"`
}

type PostSystemExport struct {
	Users []ExportUser `json:"users"`
}
