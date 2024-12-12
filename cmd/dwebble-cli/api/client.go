// DO NOT EDIT THIS: This file was generated by the Pyrin Golang Generator
package api

func (c *Client) GetArtists(options Options) (*GetArtists, error) {
	path := "/api/v1/artists"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetArtists](data)
}

func (c *Client) SearchArtists(options Options) (*GetArtists, error) {
	path := "/api/v1/artists/search"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetArtists](data)
}

func (c *Client) GetArtistById(id string, options Options) (*GetArtistById, error) {
	path := Sprintf("/api/v1/artists/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetArtistById](data)
}

func (c *Client) GetArtistAlbums(id string, options Options) (*GetArtistAlbumsById, error) {
	path := Sprintf("/api/v1/artists/%v/albums", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetArtistAlbumsById](data)
}

func (c *Client) EditArtist(id string, body EditArtistBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/artists/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) ChangeArtistPicture(id string, body Reader, options Options) (*any, error) {
	path := Sprintf("/api/v1/artists/%v/picture", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: nil,
	}
	return RequestForm[any](data, options.Boundary, body)
}

func (c *Client) GetAlbums(options Options) (*GetAlbums, error) {
	path := "/api/v1/albums"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetAlbums](data)
}

func (c *Client) SearchAlbums(options Options) (*GetAlbums, error) {
	path := "/api/v1/albums/search"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetAlbums](data)
}

func (c *Client) GetAlbumById(id string, options Options) (*GetAlbumById, error) {
	path := Sprintf("/api/v1/albums/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetAlbumById](data)
}

func (c *Client) GetAlbumTracks(id string, options Options) (*GetAlbumTracks, error) {
	path := Sprintf("/api/v1/albums/%v/tracks", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetAlbumTracks](data)
}

func (c *Client) EditAlbum(id string, body EditAlbumBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/albums/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) DeleteAlbum(id string, options Options) (*any, error) {
	path := Sprintf("/api/v1/albums/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "DELETE",
		Token: c.token,
		Body: nil,
	}
	return Request[any](data)
}

func (c *Client) CreateAlbum(body CreateAlbumBody, options Options) (*CreateAlbum, error) {
	path := "/api/v1/albums"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[CreateAlbum](data)
}

func (c *Client) ChangeAlbumCover(id string, body Reader, options Options) (*any, error) {
	path := Sprintf("/api/v1/albums/%v/cover", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: nil,
	}
	return RequestForm[any](data, options.Boundary, body)
}

func (c *Client) UploadTracks(id string, body Reader, options Options) (*any, error) {
	path := Sprintf("/api/v1/albums/%v/upload", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return RequestForm[any](data, options.Boundary, body)
}

func (c *Client) GetTracks(options Options) (*GetTracks, error) {
	path := "/api/v1/tracks"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetTracks](data)
}

func (c *Client) SearchTracks(options Options) (*GetTracks, error) {
	path := "/api/v1/tracks/search"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetTracks](data)
}

func (c *Client) GetTrackById(id string, options Options) (*GetTrackById, error) {
	path := Sprintf("/api/v1/tracks/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetTrackById](data)
}

func (c *Client) RemoveTrack(id string, options Options) (*any, error) {
	path := Sprintf("/api/v1/tracks/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "DELETE",
		Token: c.token,
		Body: nil,
	}
	return Request[any](data)
}

func (c *Client) EditTrack(id string, body EditTrackBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/tracks/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) Signup(body SignupBody, options Options) (*Signup, error) {
	path := "/api/v1/auth/signup"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[Signup](data)
}

func (c *Client) Signin(body SigninBody, options Options) (*Signin, error) {
	path := "/api/v1/auth/signin"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[Signin](data)
}

func (c *Client) ChangePassword(body ChangePasswordBody, options Options) (*any, error) {
	path := "/api/v1/auth/password"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) GetMe(options Options) (*GetMe, error) {
	path := "/api/v1/auth/me"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetMe](data)
}

func (c *Client) GetPlaylists(options Options) (*GetPlaylists, error) {
	path := "/api/v1/playlists"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetPlaylists](data)
}

func (c *Client) CreatePlaylist(body CreatePlaylistBody, options Options) (*CreatePlaylist, error) {
	path := "/api/v1/playlists"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[CreatePlaylist](data)
}

func (c *Client) CreatePlaylistFromFilter(body PostPlaylistFilterBody, options Options) (*CreatePlaylist, error) {
	path := "/api/v1/playlists/filter"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[CreatePlaylist](data)
}

func (c *Client) GetPlaylistById(id string, options Options) (*GetPlaylistById, error) {
	path := Sprintf("/api/v1/playlists/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetPlaylistById](data)
}

func (c *Client) AddItemsToPlaylist(id string, body AddItemsToPlaylistBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/playlists/%v/items", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) DeletePlaylistItems(id string, body DeletePlaylistItemsBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/playlists/%v/items", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "DELETE",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) MovePlaylistItem(id string, body MovePlaylistItemBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/playlists/%v/items/move", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) GetSystemInfo(options Options) (*GetSystemInfo, error) {
	path := "/api/v1/system/info"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetSystemInfo](data)
}

func (c *Client) SystemExport(options Options) (*Export, error) {
	path := "/api/v1/system/export"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: nil,
	}
	return Request[Export](data)
}

func (c *Client) Process(options Options) (*any, error) {
	path := "/api/v1/system/process"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: nil,
	}
	return Request[any](data)
}

func (c *Client) GetTaglists(options Options) (*GetTaglists, error) {
	path := "/api/v1/taglists"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetTaglists](data)
}

func (c *Client) GetTaglistById(id string, options Options) (*GetTaglistById, error) {
	path := Sprintf("/api/v1/taglists/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetTaglistById](data)
}

func (c *Client) GetTaglistTracks(id string, options Options) (*GetTaglistTracks, error) {
	path := Sprintf("/api/v1/taglists/%v/tracks", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetTaglistTracks](data)
}

func (c *Client) CreateTaglist(body CreateTaglistBody, options Options) (*CreateTaglist, error) {
	path := "/api/v1/taglists"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[CreateTaglist](data)
}

func (c *Client) UpdateTaglist(id string, body UpdateTaglistBody, options Options) (*any, error) {
	path := Sprintf("/api/v1/taglists/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) UpdateUserSettings(body UpdateUserSettingsBody, options Options) (*any, error) {
	path := "/api/v1/user/settings"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "PATCH",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) AddToUserQuickPlaylist(body TrackIds, options Options) (*any, error) {
	path := "/api/v1/user/quickplaylist"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "POST",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) RemoveItemFromUserQuickPlaylist(body TrackIds, options Options) (*any, error) {
	path := "/api/v1/user/quickplaylist"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "DELETE",
		Token: c.token,
		Body: body,
	}
	return Request[any](data)
}

func (c *Client) GetUserQuickPlaylistItemIds(options Options) (*GetUserQuickPlaylistItemIds, error) {
	path := "/api/v1/user/quickplaylist"
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: "GET",
		Token: c.token,
		Body: nil,
	}
	return Request[GetUserQuickPlaylistItemIds](data)
}
