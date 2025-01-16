import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  getArtists(options?: ExtraOptions) {
    return this.request("/api/v1/artists", "GET", api.GetArtists, z.any(), undefined, options)
  }
  
  searchArtists(options?: ExtraOptions) {
    return this.request("/api/v1/artists/search", "GET", api.GetArtists, z.any(), undefined, options)
  }
  
  getArtistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}`, "GET", api.GetArtistById, z.any(), undefined, options)
  }
  
  getArtistAlbums(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}/albums`, "GET", api.GetArtistAlbumsById, z.any(), undefined, options)
  }
  
  editArtist(id: string, body: api.EditArtistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  createArtist(body: api.CreateArtistBody, options?: ExtraOptions) {
    return this.request("/api/v1/artists", "POST", api.CreateArtist, z.any(), body, options)
  }
  
  mergeArtists(id: string, body: api.MergeArtistsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}/merge`, "POST", z.undefined(), z.any(), body, options)
  }
  
  deleteArtist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getAlbums(options?: ExtraOptions) {
    return this.request("/api/v1/albums", "GET", api.GetAlbums, z.any(), undefined, options)
  }
  
  searchAlbums(options?: ExtraOptions) {
    return this.request("/api/v1/albums/search", "GET", api.GetAlbums, z.any(), undefined, options)
  }
  
  getAlbumById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "GET", api.GetAlbumById, z.any(), undefined, options)
  }
  
  getAlbumTracks(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}/tracks`, "GET", api.GetAlbumTracks, z.any(), undefined, options)
  }
  
  getAlbumTracksForPlay(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}/tracks/play`, "GET", api.GetAlbumTracksForPlay, z.any(), undefined, options)
  }
  
  editAlbum(id: string, body: api.EditAlbumBody, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  deleteAlbum(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  createAlbum(body: api.CreateAlbumBody, options?: ExtraOptions) {
    return this.request("/api/v1/albums", "POST", api.CreateAlbum, z.any(), body, options)
  }
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getDetailedTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks/detailed", "GET", api.GetDetailedTracks, z.any(), undefined, options)
  }
  
  searchTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks/search", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, z.any(), undefined, options)
  }
  
  getDetailedTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}/detailed`, "GET", api.GetDetailedTrackById, z.any(), undefined, options)
  }
  
  removeTrack(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editTrack(id: string, body: api.EditTrackBody, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  deleteTrack(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getDefaultQueue(playerId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/default/${playerId}`, "GET", api.Queue, z.any(), undefined, options)
  }
  
  clearQueue(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/${id}/clear`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  addToQueueFromAlbum(id: string, albumId: string, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/${id}/add/album/${albumId}`, "POST", z.undefined(), z.any(), undefined, options)
  }
  
  addToQueueFromPlaylist(id: string, playlistId: string, body: api.AddToQueue, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/${id}/add/playlist/${playlistId}`, "POST", z.undefined(), z.any(), body, options)
  }
  
  getQueueItems(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/${id}/items`, "GET", api.GetQueueItems, z.any(), undefined, options)
  }
  
  updateQueue(id: string, body: api.UpdateQueueBody, options?: ExtraOptions) {
    return this.request(`/api/v1/queue/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  signup(body: api.SignupBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signup", "POST", api.Signup, z.any(), body, options)
  }
  
  signin(body: api.SigninBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signin", "POST", api.Signin, z.any(), body, options)
  }
  
  changePassword(body: api.ChangePasswordBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/password", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
  }
  
  getPlaylists(options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "GET", api.GetPlaylists, z.any(), undefined, options)
  }
  
  createPlaylist(body: api.CreatePlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  createPlaylistFromFilter(body: api.PostPlaylistFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists/filter", "POST", api.CreatePlaylist, z.any(), body, options)
  }
  
  getPlaylistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}`, "GET", api.GetPlaylistById, z.any(), undefined, options)
  }
  
  getPlaylistItems(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "GET", api.GetPlaylistItems, z.any(), undefined, options)
  }
  
  addItemToPlaylist(id: string, body: api.AddItemToPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  removePlaylistItem(id: string, body: api.RemovePlaylistItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "DELETE", z.undefined(), z.any(), body, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  systemExport(options?: ExtraOptions) {
    return this.request("/api/v1/system/export", "POST", api.Export, z.any(), undefined, options)
  }
  
  process(options?: ExtraOptions) {
    return this.request("/api/v1/system/process", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  refillSearch(options?: ExtraOptions) {
    return this.request("/api/v1/system/search", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  getTaglists(options?: ExtraOptions) {
    return this.request("/api/v1/taglists", "GET", api.GetTaglists, z.any(), undefined, options)
  }
  
  getTaglistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/taglists/${id}`, "GET", api.GetTaglistById, z.any(), undefined, options)
  }
  
  getTaglistTracks(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/taglists/${id}/tracks`, "GET", api.GetTaglistTracks, z.any(), undefined, options)
  }
  
  createTaglist(body: api.CreateTaglistBody, options?: ExtraOptions) {
    return this.request("/api/v1/taglists", "POST", api.CreateTaglist, z.any(), body, options)
  }
  
  deleteTaglist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/taglists/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  updateTaglist(id: string, body: api.UpdateTaglistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/taglists/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  addToUserQuickPlaylist(body: api.TrackId, options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "POST", z.undefined(), z.any(), body, options)
  }
  
  removeItemFromUserQuickPlaylist(body: api.TrackId, options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "DELETE", z.undefined(), z.any(), body, options)
  }
  
  getUserQuickPlaylistItemIds(options?: ExtraOptions) {
    return this.request("/api/v1/user/quickplaylist", "GET", api.GetUserQuickPlaylistItemIds, z.any(), undefined, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  getAllApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "GET", api.GetAllApiTokens, z.any(), undefined, options)
  }
  
  removeApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  changeArtistPicture(id: string, formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData(`/api/v1/artists/${id}/picture`, "POST", z.undefined(), z.undefined(), formData, options)
  }
  
  changeAlbumCover(id: string, formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData(`/api/v1/albums/${id}/cover`, "POST", z.undefined(), z.undefined(), formData, options)
  }
  
  uploadTrack(formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData("/api/v1/tracks", "POST", z.undefined(), z.undefined(), formData, options)
  }
}
