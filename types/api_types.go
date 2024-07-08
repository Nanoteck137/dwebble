package types

import "github.com/faceair/jio"

type Artist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Album struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type Track struct {
	Id                string `json:"id"`
	Number            int    `json:"number"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	Duration          int    `json:"duration"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId           string `json:"albumId"`
	ArtistId          string `json:"artistId"`
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
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

type PostAuthSignupBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

var PostAuthSignupBodySchema = jio.Object().Keys(jio.K{
	"username":        jio.String().Min(4).Required(),
	"password":        jio.String().Min(8).Required(),
	"passwordConfirm": jio.String().Min(8).Required(),
})

type PostAuthSignup struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type PostAuthSigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var PostAuthSigninBodySchema = jio.Object().Keys(jio.K{
	"username": jio.String().Required(),
	"password": jio.String().Required(),
})

type PostAuthSignin struct {
	Token string `json:"token"`
}

type GetAuthMe struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type PostPlaylist struct {
	Playlist
}

type PostPlaylistBody struct {
	Name string `json:"name"`
}

var PostPlaylistBodySchema = jio.Object().Keys(jio.K{
	"name": jio.String().Required(),
})

type PostPlaylistItemsByIdBody struct {
	Tracks []string `json:"tracks"`
}

var PostPlaylistItemsByIdBodySchema = jio.Object().Keys(jio.K{
	"tracks": jio.Array().Items(jio.String()).Min(1).Required(),
})

type GetPlaylistById struct {
	Playlist

	Items []Track `json:"items"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type DeletePlaylistItemsByIdBody struct {
	TrackIndices []int `json:"trackIndices"`
}

var DeletePlaylistItemsByIdBodySchema = jio.Object().Keys(jio.K{
	"trackIndices": jio.Array().Items(jio.Number().Integer()).Min(1).Required(),
})

type PostPlaylistsItemMoveByIdBody struct {
	TrackId string `json:"trackId"`
	ToIndex int    `json:"toIndex"`
}

var PostPlaylistsItemMoveByIdBodySchema = jio.Object().Keys(jio.K{
	"trackId": jio.String().Required(),
	"toIndex": jio.Number().Integer().Required(),
})

type GetSystemInfo struct {
	Version string `json:"version"`
	IsSetup bool   `json:"isSetup"`
}

type PostSystemSetupBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

var PostSystemSetupBodySchema = jio.Object().Keys(jio.K{
	"username":        jio.String().Required(),
	"password":        jio.String().Required(),
	"passwordConfirm": jio.String().Required(),
})

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