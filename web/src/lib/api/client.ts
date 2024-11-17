import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";

export const GET_ARTISTS_URL = "/api/v1/artists"
export const GET_ARTIST_BY_ID_URL = "/api/v1/artists/:id"
export const GET_ARTIST_ALBUMS_URL = "/api/v1/artists/:id/albums"
export const GET_ALBUMS_URL = "/api/v1/albums"
export const GET_ALBUM_BY_ID_URL = "/api/v1/albums/:id"
export const GET_ALBUM_TRACKS_URL = "/api/v1/albums/:id/tracks"
export const EDIT_ALBUM_URL = "/api/v1/albums/:id"
export const DELETE_ALBUM_URL = "/api/v1/albums/:id"
export const IMPORT_ALBUM_URL = "/api/v1/albums/import"
export const IMPORT_TRACK_TO_ALBUM_URL = "/api/v1/albums/:id/import/track"
export const GET_TRACKS_URL = "/api/v1/tracks"
export const GET_TRACK_BY_ID_URL = "/api/v1/tracks/:id"
export const REMOVE_TRACK_URL = "/api/v1/tracks/:id"
export const EDIT_TRACK_URL = "/api/v1/tracks/:id"
export const SIGNUP_URL = "/api/v1/auth/signup"
export const SIGNIN_URL = "/api/v1/auth/signin"
export const GET_ME_URL = "/api/v1/auth/me"
export const GET_PLAYLISTS_URL = "/api/v1/playlists"
export const CREATE_PLAYLIST_URL = "/api/v1/playlists"
export const CREATE_PLAYLIST_FROM_FILTER_URL = "/api/v1/playlists/filter"
export const GET_PLAYLIST_BY_ID_URL = "/api/v1/playlists/:id"
export const ADD_ITEMS_TO_PLAYLIST_URL = "/api/v1/playlists/:id/items"
export const DELETE_PLAYLIST_ITEMS_URL = "/api/v1/playlists/:id/items"
export const MOVE_PLAYLIST_ITEM_URL = "/api/v1/playlists/:id/items/move"
export const GET_SYSTEM_INFO_URL = "/api/v1/system/info"
export const SYSTEM_EXPORT_URL = "/api/v1/system/export"
export const SYSTEM_IMPORT_URL = "/api/v1/system/import"
export const PROCESS_URL = "/api/v1/system/process"

export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  getArtists(options?: ExtraOptions) {
    return this.request("/api/v1/artists", "GET", api.GetArtists, z.undefined(), undefined, options)
  }
  
  getArtistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}`, "GET", api.GetArtistById, z.undefined(), undefined, options)
  }
  
  getArtistAlbums(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/artists/${id}/albums`, "GET", api.GetArtistAlbumsById, z.undefined(), undefined, options)
  }
  
  getAlbums(options?: ExtraOptions) {
    return this.request("/api/v1/albums", "GET", api.GetAlbums, z.undefined(), undefined, options)
  }
  
  getAlbumById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "GET", api.GetAlbumById, z.undefined(), undefined, options)
  }
  
  getAlbumTracks(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}/tracks`, "GET", api.GetAlbumTracksById, z.undefined(), undefined, options)
  }
  
  editAlbum(id: string, body: api.PatchAlbumBody, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "PATCH", z.undefined(), z.undefined(), body, options)
  }
  
  deleteAlbum(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "DELETE", z.undefined(), z.undefined(), undefined, options)
  }
  
  importAlbum(formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData("/api/v1/albums/import", "POST", api.PostAlbumImport, z.undefined(), formData, options)
  }
  
  importTrackToAlbum(id: string, formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData(`/api/v1/albums/${id}/import/track`, "POST", z.undefined(), z.undefined(), formData, options)
  }
  
  getTracks(options?: ExtraOptions) {
    return this.request("/api/v1/tracks", "GET", api.GetTracks, z.undefined(), undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, z.undefined(), undefined, options)
  }
  
  removeTrack(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "DELETE", z.undefined(), z.undefined(), undefined, options)
  }
  
  editTrack(id: string, body: api.PatchTrackBody, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "PATCH", z.undefined(), z.undefined(), body, options)
  }
  
  signup(body: api.PostAuthSignupBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signup", "POST", api.PostAuthSignup, z.undefined(), body, options)
  }
  
  signin(body: api.PostAuthSigninBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signin", "POST", api.PostAuthSignin, z.undefined(), body, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetAuthMe, z.undefined(), undefined, options)
  }
  
  getPlaylists(options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "GET", api.GetPlaylists, z.undefined(), undefined, options)
  }
  
  createPlaylist(body: api.PostPlaylistBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists", "POST", api.PostPlaylist, z.undefined(), body, options)
  }
  
  createPlaylistFromFilter(body: api.PostPlaylistFilterBody, options?: ExtraOptions) {
    return this.request("/api/v1/playlists/filter", "POST", api.PostPlaylist, z.undefined(), body, options)
  }
  
  getPlaylistById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}`, "GET", api.GetPlaylistById, z.undefined(), undefined, options)
  }
  
  addItemsToPlaylist(id: string, body: api.PostPlaylistItemsByIdBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "POST", z.undefined(), z.undefined(), body, options)
  }
  
  deletePlaylistItems(id: string, body: api.DeletePlaylistItemsByIdBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "DELETE", z.undefined(), z.undefined(), body, options)
  }
  
  movePlaylistItem(id: string, body: api.PostPlaylistsItemMoveByIdBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items/move`, "POST", z.undefined(), z.undefined(), body, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.undefined(), undefined, options)
  }
  
  systemExport(options?: ExtraOptions) {
    return this.request("/api/v1/system/export", "POST", z.undefined(), z.undefined(), undefined, options)
  }
  
  systemImport(options?: ExtraOptions) {
    return this.request("/api/v1/system/import", "POST", z.undefined(), z.undefined(), undefined, options)
  }
  
  process(options?: ExtraOptions) {
    return this.request("/api/v1/system/process", "POST", z.undefined(), z.undefined(), undefined, options)
  }
}
