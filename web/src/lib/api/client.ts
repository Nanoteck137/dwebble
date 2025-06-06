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
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  searchTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks/search", "GET", api.GetTracks, z.any(), undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, z.any(), undefined, options)
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
  
  deletePlaylist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
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
  
  clearPlaylist(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items/all`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  refillSearch(options?: ExtraOptions) {
    return this.request("/api/v1/system/search", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  getSyncStatus(options?: ExtraOptions) {
    return this.request("/api/v1/system/library", "GET", z.undefined(), z.any(), undefined, options)
  }
  
  syncLibrary(options?: ExtraOptions) {
    return this.request("/api/v1/system/library", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  cleanupLibrary(options?: ExtraOptions) {
    return this.request("/api/v1/system/library/cleanup", "POST", z.undefined(), z.any(), undefined, options)
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
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getMediaFromPlaylist(playlistId: string, body: api.GetMediaFromPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/playlist/${playlistId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromTaglist(taglistId: string, body: api.GetMediaFromTaglistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/taglist/${taglistId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromFilter(body: api.GetMediaFromFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/media/filter", "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromArtist(artistId: string, body: api.GetMediaFromArtistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/artist/${artistId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromAlbum(albumId: string, body: api.GetMediaFromAlbumBody, options?: ExtraOptions) {
    return this.request(`/api/v1/media/album/${albumId}`, "POST", api.GetMedia, z.any(), body, options)
  }
  
  getMediaFromIds(body: api.GetMediaFromIdsBody, options?: ExtraOptions) {
    return this.request("/api/v1/media/ids", "POST", api.GetMedia, z.any(), body, options)
  }
}
