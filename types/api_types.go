package types

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

type PostAuthSigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
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

type PostPlaylistBody struct {
	Name string `json:"name"`
}

type PostPlaylistFilterBody struct {
	Name   string `json:"name"`
	Filter string `json:"filter"`
	Sort   string `json:"sort"`
}

type PostPlaylistItemsByIdBody struct {
	Tracks []string `json:"tracks"`
}

type GetPlaylistById struct {
	Playlist

	Items []Track `json:"items"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type DeletePlaylistItemsByIdBody struct {
	TrackIds []string `json:"trackIds"`
}

type PostPlaylistsItemMoveByIdBody struct {
	TrackId string `json:"trackId"`
	ToIndex int    `json:"toIndex"`
}

type GetSystemInfo struct {
	Version string `json:"version"`
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
