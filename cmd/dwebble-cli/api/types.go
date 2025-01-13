// THIS FILE IS GENERATED BY PYRIN GOGEN CODE GENERATOR
package api

type Name struct {
	Default string `json:"default"`
	Other *string `json:"other,omitempty"`
}

type Images struct {
	Original string `json:"original"`
	Small string `json:"small"`
	Medium string `json:"medium"`
	Large string `json:"large"`
}

type Artist struct {
	Id string `json:"id"`
	Name Name `json:"name"`
	Picture Images `json:"picture"`
	Tags []string `json:"tags"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

type GetArtists struct {
	Artists []Artist `json:"artists"`
}

type GetArtistById Artist

type ArtistInfo struct {
	Id string `json:"id"`
	Name Name `json:"name"`
}

type Album struct {
	Id string `json:"id"`
	Name Name `json:"name"`
	Year *int `json:"year,omitempty"`
	CoverArt Images `json:"coverArt"`
	ArtistId string `json:"artistId"`
	ArtistName Name `json:"artistName"`
	Tags []string `json:"tags"`
	FeaturingArtists []ArtistInfo `json:"featuringArtists"`
	AllArtists []ArtistInfo `json:"allArtists"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type EditArtistBody struct {
	Name *string `json:"name,omitempty"`
	OtherName *string `json:"otherName,omitempty"`
	Tags *[]string `json:"tags,omitempty"`
}

type CreateArtist struct {
	Id string `json:"id"`
}

type CreateArtistBody struct {
	Name string `json:"name"`
	OtherName string `json:"otherName"`
}

type MergeArtistsBody struct {
	Artists []string `json:"artists"`
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}

type GetAlbumById Album

type Track struct {
	Id string `json:"id"`
	Name Name `json:"name"`
	Duration int `json:"duration"`
	Number *int `json:"number,omitempty"`
	Year *int `json:"year,omitempty"`
	OriginalMediaUrl string `json:"originalMediaUrl"`
	MobileMediaUrl string `json:"mobileMediaUrl"`
	CoverArt Images `json:"coverArt"`
	AlbumId string `json:"albumId"`
	ArtistId string `json:"artistId"`
	AlbumName Name `json:"albumName"`
	ArtistName Name `json:"artistName"`
	Tags []string `json:"tags"`
	FeaturingArtists []ArtistInfo `json:"featuringArtists"`
	AllArtists []ArtistInfo `json:"allArtists"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

type GetAlbumTracks struct {
	Tracks []Track `json:"tracks"`
}

type EditAlbumBody struct {
	Name *string `json:"name,omitempty"`
	OtherName *string `json:"otherName,omitempty"`
	ArtistId *string `json:"artistId,omitempty"`
	Year *int `json:"year,omitempty"`
	Tags *[]string `json:"tags,omitempty"`
	FeaturingArtists *[]string `json:"featuringArtists,omitempty"`
}

type CreateAlbum struct {
	AlbumId string `json:"albumId"`
}

type CreateAlbumBody struct {
	Name string `json:"name"`
	OtherName string `json:"otherName"`
	ArtistId string `json:"artistId"`
	Year int `json:"year"`
	Tags []string `json:"tags"`
	FeaturingArtists []string `json:"featuringArtists"`
}

type Page struct {
	Page int `json:"page"`
	PerPage int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type GetTracks struct {
	Page Page `json:"page"`
	Tracks []Track `json:"tracks"`
}

type GetTrackById Track

type EditTrackBody struct {
	Name *string `json:"name,omitempty"`
	OtherName *string `json:"otherName,omitempty"`
	ArtistId *string `json:"artistId,omitempty"`
	Year *int `json:"year,omitempty"`
	Number *int `json:"number,omitempty"`
	Tags *[]string `json:"tags,omitempty"`
	FeaturingArtists *[]string `json:"featuringArtists,omitempty"`
}

type UploadTrackBody struct {
	Name string `json:"name"`
	OtherName string `json:"otherName"`
	Number int `json:"number"`
	Year int `json:"year"`
	AlbumId string `json:"albumId"`
	ArtistId string `json:"artistId"`
	Tags []string `json:"tags"`
	FeaturingArtists []string `json:"featuringArtists"`
}

type Queue struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type MusicTrackArtist struct {
	ArtistId string `json:"artistId"`
	ArtistName string `json:"artistName"`
}

type MusicTrackAlbum struct {
	AlbumId string `json:"albumId"`
	AlbumName string `json:"albumName"`
}

type MusicTrack struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Artists []MusicTrackArtist `json:"artists"`
	Album MusicTrackAlbum `json:"album"`
	CoverArt Images `json:"coverArt"`
	MediaUrl string `json:"mediaUrl"`
}

type QueueItem struct {
	Id string `json:"id"`
	QueueId string `json:"queueId"`
	OrderNumber int `json:"orderNumber"`
	TrackId string `json:"trackId"`
	Track MusicTrack `json:"track"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

type GetQueueItems struct {
	Index int `json:"index"`
	Items []QueueItem `json:"items"`
}

type UpdateQueueBody struct {
	Name *string `json:"name,omitempty"`
	ItemIndex *int `json:"itemIndex,omitempty"`
}

type Signup struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

type SignupBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type Signin struct {
	Token string `json:"token"`
}

type SigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePasswordBody struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword string `json:"newPassword"`
	NewPasswordConfirm string `json:"newPasswordConfirm"`
}

type GetMe struct {
	Id string `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
	DisplayName string `json:"displayName"`
	QuickPlaylist *string `json:"quickPlaylist,omitempty"`
}

type Playlist struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type CreatePlaylist Playlist

type CreatePlaylistBody struct {
	Name string `json:"name"`
}

type PostPlaylistFilterBody struct {
	Name string `json:"name"`
	Filter string `json:"filter"`
}

type GetPlaylistById Playlist

type GetPlaylistItems struct {
	Items []Track `json:"items"`
}

type AddItemToPlaylistBody struct {
	TrackId string `json:"trackId"`
}

type RemovePlaylistItemBody struct {
	TrackId string `json:"trackId"`
}

type GetSystemInfo struct {
	Version string `json:"version"`
}

type ExportArtist struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Picture string `json:"picture"`
}

type ExportAlbum struct {
	Id string `json:"id"`
	Name string `json:"name"`
	ArtistId string `json:"artistId"`
	CoverArt string `json:"coverArt"`
	Year int `json:"year"`
}

type ExportTrack struct {
	Id string `json:"id"`
	Name string `json:"name"`
	AlbumId string `json:"albumId"`
	ArtistId string `json:"artistId"`
	Duration int `json:"duration"`
	Number int `json:"number"`
	Year int `json:"year"`
	OriginalFilename string `json:"originalFilename"`
	MobileFilename string `json:"mobileFilename"`
	Created int `json:"created"`
	Tags []string `json:"tags"`
}

type Export struct {
	Artists []ExportArtist `json:"artists"`
	Albums []ExportAlbum `json:"albums"`
	Tracks []ExportTrack `json:"tracks"`
}

type Taglist struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Filter string `json:"filter"`
	Created int `json:"created"`
	Updated int `json:"updated"`
}

type GetTaglists struct {
	Taglists []Taglist `json:"taglists"`
}

type GetTaglistById Taglist

type GetTaglistTracks struct {
	Page Page `json:"page"`
	Tracks []Track `json:"tracks"`
}

type CreateTaglist struct {
	Id string `json:"id"`
}

type CreateTaglistBody struct {
	Name string `json:"name"`
	Filter string `json:"filter"`
}

type UpdateTaglistBody struct {
	Name *string `json:"name,omitempty"`
	Filter *string `json:"filter,omitempty"`
}

type UpdateUserSettingsBody struct {
	DisplayName *string `json:"displayName,omitempty"`
	QuickPlaylist *string `json:"quickPlaylist,omitempty"`
}

type TrackId struct {
	TrackId string `json:"trackId"`
}

type GetUserQuickPlaylistItemIds struct {
	TrackIds []string `json:"trackIds"`
}

type CreateApiToken struct {
	Token string `json:"token"`
}

type CreateApiTokenBody struct {
	Name string `json:"name"`
}

type ApiToken struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type GetAllApiTokens struct {
	Tokens []ApiToken `json:"tokens"`
}

