// THIS FILE IS GENERATED BY PYRIN GOGEN CODE GENERATOR
package api

type Images struct {
	Original string `json:"original"`
	Small string `json:"small"`
	Medium string `json:"medium"`
	Large string `json:"large"`
}

type Artist struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Picture Images `json:"picture"`
}

type GetArtists struct {
	Artists []Artist `json:"artists"`
}

type GetArtistById Artist

type Album struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CoverArt Images `json:"coverArt"`
	ArtistId string `json:"artistId"`
	ArtistName string `json:"artistName"`
	Year *int `json:"year,omitempty"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}

type GetAlbumById Album

type Track struct {
	Id string `json:"id"`
	Name string `json:"name"`
	AlbumId string `json:"albumId"`
	ArtistId string `json:"artistId"`
	Number *int `json:"number,omitempty"`
	Duration *int `json:"duration,omitempty"`
	Year *int `json:"year,omitempty"`
	OriginalMediaUrl string `json:"originalMediaUrl"`
	MobileMediaUrl string `json:"mobileMediaUrl"`
	CoverArt Images `json:"coverArt"`
	AlbumName string `json:"albumName"`
	ArtistName string `json:"artistName"`
	Created int `json:"created"`
	Updated int `json:"updated"`
	Tags []string `json:"tags"`
}

type GetAlbumTracksById struct {
	Tracks []Track `json:"tracks"`
}

type PatchAlbumBody struct {
	Name *string `json:"name,omitempty"`
	ArtistId *string `json:"artistId,omitempty"`
	ArtistName *string `json:"artistName,omitempty"`
	Year *int `json:"year,omitempty"`
}

type PostAlbumImport struct {
	AlbumId string `json:"albumId"`
}

type PostAlbumImportBody struct {
	Name string `json:"name"`
	Artist string `json:"artist"`
}

type GetTracks struct {
	Tracks []Track `json:"tracks"`
}

type GetTrackById Track

type PatchTrackBody struct {
	Name *string `json:"name,omitempty"`
	ArtistId *string `json:"artistId,omitempty"`
	ArtistName *string `json:"artistName,omitempty"`
	Year *int `json:"year,omitempty"`
	Number *int `json:"number,omitempty"`
	Tags *[]string `json:"tags,omitempty"`
}

type PostAuthSignup struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

type PostAuthSignupBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type PostAuthSignin struct {
	Token string `json:"token"`
}

type PostAuthSigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetAuthMe struct {
	Id string `json:"id"`
	Username string `json:"username"`
	IsOwner bool `json:"isOwner"`
}

type Playlist struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type PostPlaylist Playlist

type PostPlaylistBody struct {
	Name string `json:"name"`
}

type PostPlaylistFilterBody struct {
	Name string `json:"name"`
	Filter string `json:"filter"`
	Sort string `json:"sort"`
}

type GetPlaylistById struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Items []Track `json:"items"`
}

type PostPlaylistItemsByIdBody struct {
	Tracks []string `json:"tracks"`
}

type DeletePlaylistItemsByIdBody struct {
	TrackIds []string `json:"trackIds"`
}

type PostPlaylistsItemMoveByIdBody struct {
	TrackId string `json:"trackId"`
	ToIndex int `json:"toIndex"`
}

type GetSystemInfo struct {
	Version string `json:"version"`
	IsSetup bool `json:"isSetup"`
}

