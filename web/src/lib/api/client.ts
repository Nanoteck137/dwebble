import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";

export const GET_ARTISTS_URL = "/api/v1/artists"
export const SEARCH_ARTISTS_URL = "/api/v1/artists/search"
export const GET_ARTIST_BY_ID_URL = "/api/v1/artists/:id"
export const GET_ARTIST_ALBUMS_URL = "/api/v1/artists/:id/albums"
export const EDIT_ARTIST_URL = "/api/v1/artists/:id"
export const CHANGE_ARTIST_PICTURE_URL = "/api/v1/artists/:id/picture"
export const GET_ALBUMS_URL = "/api/v1/albums"
export const SEARCH_ALBUMS_URL = "/api/v1/albums/search"
export const GET_ALBUM_BY_ID_URL = "/api/v1/albums/:id"
export const GET_ALBUM_TRACKS_URL = "/api/v1/albums/:id/tracks"
export const EDIT_ALBUM_URL = "/api/v1/albums/:id"
export const DELETE_ALBUM_URL = "/api/v1/albums/:id"
export const CREATE_ALBUM_URL = "/api/v1/albums"
export const CHANGE_ALBUM_COVER_URL = "/api/v1/albums/:id/cover"
export const UPLOAD_TRACKS_URL = "/api/v1/albums/:id/upload"
export const GET_TRACKS_URL = "/api/v1/tracks"
export const SEARCH_TRACKS_URL = "/api/v1/tracks/search"
export const GET_TRACK_BY_ID_URL = "/api/v1/tracks/:id"
export const REMOVE_TRACK_URL = "/api/v1/tracks/:id"
export const EDIT_TRACK_URL = "/api/v1/tracks/:id"
export const SIGNUP_URL = "/api/v1/auth/signup"
export const SIGNIN_URL = "/api/v1/auth/signin"
export const CHANGE_PASSWORD_URL = "/api/v1/auth/password"
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
export const PROCESS_URL = "/api/v1/system/process"
export const GET_TAGLISTS_URL = "/api/v1/taglists"
export const GET_TAGLIST_BY_ID_URL = "/api/v1/taglists/:id"
export const GET_TAGLIST_TRACKS_URL = "/api/v1/taglists/:id/tracks"
export const CREATE_TAGLIST_URL = "/api/v1/taglists"
export const UPDATE_TAGLIST_URL = "/api/v1/taglists/:id"

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
  
  changeArtistPicture(id: string, formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData(`/api/v1/artists/${id}/picture`, "PATCH", z.undefined(), z.undefined(), formData, options)
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
  
  editAlbum(id: string, body: api.EditAlbumBody, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  deleteAlbum(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/albums/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  createAlbum(body: api.CreateAlbumBody, options?: ExtraOptions) {
    return this.request("/api/v1/albums", "POST", api.CreateAlbum, z.any(), body, options)
  }
  
  changeAlbumCover(id: string, formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData(`/api/v1/albums/${id}/cover`, "POST", z.undefined(), z.undefined(), formData, options)
  }
  
  uploadTracks(id: string, formData: FormData, options?: ExtraOptions) {
    return this.requestWithFormData(`/api/v1/albums/${id}/upload`, "POST", z.undefined(), z.undefined(), formData, options)
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
  
  removeTrack(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  editTrack(id: string, body: api.EditTrackBody, options?: ExtraOptions) {
    return this.request(`/api/v1/tracks/${id}`, "PATCH", z.undefined(), z.any(), body, options)
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
  
  addItemsToPlaylist(id: string, body: api.AddItemsToPlaylistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "POST", z.undefined(), z.any(), body, options)
  }
  
  deletePlaylistItems(id: string, body: api.DeletePlaylistItemsBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items`, "DELETE", z.undefined(), z.any(), body, options)
  }
  
  movePlaylistItem(id: string, body: api.MovePlaylistItemBody, options?: ExtraOptions) {
    return this.request(`/api/v1/playlists/${id}/items/move`, "POST", z.undefined(), z.any(), body, options)
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
  
  updateTaglist(id: string, body: api.UpdateTaglistBody, options?: ExtraOptions) {
    return this.request(`/api/v1/taglists/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
}
