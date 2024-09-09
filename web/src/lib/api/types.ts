// DO NOT EDIT THIS: This file was generated by the Pyrin Typescript Generator
import { z } from "zod";

export const Artist = z.object({
  id: z.string(),
  name: z.string(),
  picture: z.string(),
});
export type Artist = z.infer<typeof Artist>;

export const GetArtists = z.object({
  artists: z.array(Artist),
});
export type GetArtists = z.infer<typeof GetArtists>;

export const GetArtistById = Artist;
export type GetArtistById = z.infer<typeof GetArtistById>;

export const Album = z.object({
  id: z.string(),
  name: z.string(),
  coverArt: z.string(),
  artistId: z.string(),
  artistName: z.string(),
  created: z.number(),
  updated: z.number(),
});
export type Album = z.infer<typeof Album>;

export const GetArtistAlbumsById = z.object({
  albums: z.array(Album),
});
export type GetArtistAlbumsById = z.infer<typeof GetArtistAlbumsById>;

export const GetAlbums = z.object({
  albums: z.array(Album),
});
export type GetAlbums = z.infer<typeof GetAlbums>;

export const GetAlbumById = Album;
export type GetAlbumById = z.infer<typeof GetAlbumById>;

export const Track = z.object({
  id: z.string(),
  name: z.string(),
  albumId: z.string(),
  artistId: z.string(),
  number: z.number().nullable(),
  duration: z.number().nullable(),
  year: z.number().nullable(),
  originalMediaUrl: z.string(),
  mobileMediaUrl: z.string(),
  coverArtUrl: z.string(),
  albumName: z.string(),
  artistName: z.string(),
  created: z.number(),
  updated: z.number(),
  tags: z.array(z.string()),
  available: z.boolean(),
});
export type Track = z.infer<typeof Track>;

export const GetAlbumTracksById = z.object({
  tracks: z.array(Track),
});
export type GetAlbumTracksById = z.infer<typeof GetAlbumTracksById>;

export const PatchAlbumBody = z.object({
  name: z.string().nullable(),
  artistId: z.string().nullable(),
  artistName: z.string().nullable(),
  year: z.number().nullable(),
});
export type PatchAlbumBody = z.infer<typeof PatchAlbumBody>;

export const PostAlbumImport = z.object({
  albumId: z.string(),
});
export type PostAlbumImport = z.infer<typeof PostAlbumImport>;

export const PostAlbumImportBody = z.object({
  name: z.string(),
  artist: z.string(),
});
export type PostAlbumImportBody = z.infer<typeof PostAlbumImportBody>;

export const GetTracks = z.object({
  tracks: z.array(Track),
});
export type GetTracks = z.infer<typeof GetTracks>;

export const GetTrackById = Track;
export type GetTrackById = z.infer<typeof GetTrackById>;

export const PatchTrackBody = z.object({
  name: z.string().nullable().optional(),
  artistId: z.string().nullable().optional(),
  artistName: z.string().nullable().optional(),
  year: z.number().nullable().optional(),
  number: z.number().nullable().optional(),
  tags: z.array(z.string()).nullable().optional(),
});
export type PatchTrackBody = z.infer<typeof PatchTrackBody>;

export const GetSync = z.object({
  isSyncing: z.boolean(),
});
export type GetSync = z.infer<typeof GetSync>;

export const Tag = z.object({
  id: z.string(),
  name: z.string(),
});
export type Tag = z.infer<typeof Tag>;

export const GetTags = z.object({
  tags: z.array(Tag),
});
export type GetTags = z.infer<typeof GetTags>;

export const PostAuthSignup = z.object({
  id: z.string(),
  username: z.string(),
});
export type PostAuthSignup = z.infer<typeof PostAuthSignup>;

export const PostAuthSignupBody = z.object({
  username: z.string(),
  password: z.string(),
  passwordConfirm: z.string(),
});
export type PostAuthSignupBody = z.infer<typeof PostAuthSignupBody>;

export const PostAuthSignin = z.object({
  token: z.string(),
});
export type PostAuthSignin = z.infer<typeof PostAuthSignin>;

export const PostAuthSigninBody = z.object({
  username: z.string(),
  password: z.string(),
});
export type PostAuthSigninBody = z.infer<typeof PostAuthSigninBody>;

export const GetAuthMe = z.object({
  id: z.string(),
  username: z.string(),
  isOwner: z.boolean(),
});
export type GetAuthMe = z.infer<typeof GetAuthMe>;

export const Playlist = z.object({
  id: z.string(),
  name: z.string(),
});
export type Playlist = z.infer<typeof Playlist>;

export const GetPlaylists = z.object({
  playlists: z.array(Playlist),
});
export type GetPlaylists = z.infer<typeof GetPlaylists>;

export const PostPlaylist = Playlist;
export type PostPlaylist = z.infer<typeof PostPlaylist>;

export const PostPlaylistBody = z.object({
  name: z.string(),
});
export type PostPlaylistBody = z.infer<typeof PostPlaylistBody>;

export const PostPlaylistFilterBody = z.object({
  name: z.string(),
  filter: z.string(),
  sort: z.string(),
});
export type PostPlaylistFilterBody = z.infer<typeof PostPlaylistFilterBody>;

export const GetPlaylistById = z.object({
  id: z.string(),
  name: z.string(),
  items: z.array(Track),
});
export type GetPlaylistById = z.infer<typeof GetPlaylistById>;

export const PostPlaylistItemsByIdBody = z.object({
  tracks: z.array(z.string()),
});
export type PostPlaylistItemsByIdBody = z.infer<typeof PostPlaylistItemsByIdBody>;

export const DeletePlaylistItemsByIdBody = z.object({
  trackIds: z.array(z.string()),
});
export type DeletePlaylistItemsByIdBody = z.infer<typeof DeletePlaylistItemsByIdBody>;

export const PostPlaylistsItemMoveByIdBody = z.object({
  trackId: z.string(),
  toIndex: z.number(),
});
export type PostPlaylistsItemMoveByIdBody = z.infer<typeof PostPlaylistsItemMoveByIdBody>;

export const GetSystemInfo = z.object({
  version: z.string(),
  isSetup: z.boolean(),
});
export type GetSystemInfo = z.infer<typeof GetSystemInfo>;

