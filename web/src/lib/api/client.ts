import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, createError, type ExtraOptions } from "./base-client";

export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  getArtists(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/artists", "GET", api.GetArtists, error, undefined, options)
  }
  
  getArtistById(id: string, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR", "ARTIST_NOT_FOUND"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/artists/${id}`, "GET", api.GetArtistById, error, undefined, options)
  }
  
  getArtistAlbums(id: string, options?: ExtraOptions) {
    const error = createError(
      z.enum(["ARTIST_NOT_FOUND", "UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/artists/${id}/albums`, "GET", api.GetArtistAlbumsById, error, undefined, options)
  }
  
  getAlbums(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/albums", "GET", api.GetAlbums, error, undefined, options)
  }
  
  getAlbumById(id: string, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR", "ALBUM_NOT_FOUND"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/albums/${id}`, "GET", api.GetAlbumById, error, undefined, options)
  }
  
  getAlbumTracks(id: string, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR", "ALBUM_NOT_FOUND"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/albums/${id}/tracks`, "GET", api.GetAlbumTracksById, error, undefined, options)
  }
  
  getTracks(options?: ExtraOptions) {
    const error = createError(
      z.enum(["INVALID_SORT", "UNKNOWN_ERROR", "INVALID_FILTER"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/tracks", "GET", api.GetTracks, error, undefined, options)
  }
  
  getTrackById(id: string, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR", "TRACK_NOT_FOUND"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/tracks/${id}`, "GET", api.GetTrackById, error, undefined, options)
  }
  
  getSyncStatus(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/sync", "GET", api.GetSync, error, undefined, options)
  }
  
  runSync(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/sync", "POST", z.undefined(), error, undefined, options)
  }
  
  getTags(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/tags", "GET", api.GetTags, error, undefined, options)
  }
  
  signup(body: api.PostAuthSignupBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/auth/signup", "POST", api.PostAuthSignup, error, body, options)
  }
  
  signin(body: api.PostAuthSigninBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/auth/signin", "POST", api.PostAuthSignin, error, body, options)
  }
  
  getMe(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/auth/me", "GET", api.GetAuthMe, error, undefined, options)
  }
  
  getPlaylists(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/playlists", "GET", api.GetPlaylists, error, undefined, options)
  }
  
  createPlaylist(body: api.PostPlaylistBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/playlists", "POST", api.PostPlaylist, error, body, options)
  }
  
  createPlaylistFromFilter(body: api.PostPlaylistFilterBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/playlists/filter", "POST", api.PostPlaylist, error, body, options)
  }
  
  getPlaylistById(id: string, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/playlists/${id}`, "GET", api.GetPlaylistById, error, undefined, options)
  }
  
  addItemsToPlaylist(id: string, body: api.PostPlaylistItemsByIdBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/playlists/${id}/items`, "POST", z.undefined(), error, body, options)
  }
  
  deletePlaylistItems(id: string, body: api.DeletePlaylistItemsByIdBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/playlists/${id}/items`, "DELETE", z.undefined(), error, body, options)
  }
  
  movePlaylistItem(id: string, body: api.PostPlaylistsItemMoveByIdBody, options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request(`/api/v1/playlists/${id}/items/move`, "POST", z.undefined(), error, body, options)
  }
  
  getSystemInfo(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, error, undefined, options)
  }
  
  systemExport(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/system/export", "POST", api.PostSystemExport, error, undefined, options)
  }
  
  systemImport(options?: ExtraOptions) {
    const error = createError(
      z.enum(["UNKNOWN_ERROR"]),
      z.map(z.string(), z.string()).optional(),
    )
    return this.request("/api/v1/system/import", "POST", z.undefined(), error, undefined, options)
  }
}
