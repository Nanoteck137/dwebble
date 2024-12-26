// DO NOT EDIT THIS: This file was generated by the Pyrin Typescript Generator
import { z } from "zod";

export const Name = z.object({
  default: z.string(),
  other: z.string().nullable(),
});
export type Name = z.infer<typeof Name>;

export const Images = z.object({
  original: z.string(),
  small: z.string(),
  medium: z.string(),
  large: z.string(),
});
export type Images = z.infer<typeof Images>;

export const Artist = z.object({
  id: z.string(),
  name: Name,
  picture: Images,
  tags: z.array(z.string()),
  created: z.number(),
  updated: z.number(),
});
export type Artist = z.infer<typeof Artist>;

export const GetArtists = z.object({
  artists: z.array(Artist),
});
export type GetArtists = z.infer<typeof GetArtists>;

export const GetArtistById = Artist;
export type GetArtistById = z.infer<typeof GetArtistById>;

export const ArtistInfo = z.object({
  id: z.string(),
  name: Name,
});
export type ArtistInfo = z.infer<typeof ArtistInfo>;

export const Album = z.object({
  id: z.string(),
  name: Name,
  year: z.number().nullable(),
  coverArt: Images,
  artistId: z.string(),
  artistName: Name,
  tags: z.array(z.string()),
  featuringArtists: z.array(ArtistInfo),
  allArtists: z.array(ArtistInfo),
  created: z.number(),
  updated: z.number(),
});
export type Album = z.infer<typeof Album>;

export const GetArtistAlbumsById = z.object({
  albums: z.array(Album),
});
export type GetArtistAlbumsById = z.infer<typeof GetArtistAlbumsById>;

export const EditArtistBody = z.object({
  name: z.string().nullable(),
  tags: z.array(z.string()).nullable().optional(),
});
export type EditArtistBody = z.infer<typeof EditArtistBody>;

export const CreateArtist = z.object({
  id: z.string(),
});
export type CreateArtist = z.infer<typeof CreateArtist>;

export const CreateArtistBody = z.object({
  name: z.string(),
  otherName: z.string(),
});
export type CreateArtistBody = z.infer<typeof CreateArtistBody>;

export const GetAlbums = z.object({
  albums: z.array(Album),
});
export type GetAlbums = z.infer<typeof GetAlbums>;

export const GetAlbumById = Album;
export type GetAlbumById = z.infer<typeof GetAlbumById>;

export const Track = z.object({
  id: z.string(),
  name: Name,
  duration: z.number(),
  number: z.number().nullable(),
  year: z.number().nullable(),
  originalMediaUrl: z.string(),
  mobileMediaUrl: z.string(),
  coverArt: Images,
  albumId: z.string(),
  artistId: z.string(),
  albumName: Name,
  artistName: Name,
  tags: z.array(z.string()),
  featuringArtists: z.array(ArtistInfo),
  allArtists: z.array(ArtistInfo),
  created: z.number(),
  updated: z.number(),
});
export type Track = z.infer<typeof Track>;

export const GetAlbumTracks = z.object({
  tracks: z.array(Track),
});
export type GetAlbumTracks = z.infer<typeof GetAlbumTracks>;

export const EditAlbumBody = z.object({
  name: z.string().nullable().optional(),
  otherName: z.string().nullable().optional(),
  artistId: z.string().nullable().optional(),
  artistName: z.string().nullable().optional(),
  year: z.number().nullable().optional(),
  tags: z.array(z.string()).nullable().optional(),
  featuringArtists: z.array(z.string()).nullable().optional(),
});
export type EditAlbumBody = z.infer<typeof EditAlbumBody>;

export const CreateAlbum = z.object({
  albumId: z.string(),
});
export type CreateAlbum = z.infer<typeof CreateAlbum>;

export const CreateAlbumBody = z.object({
  name: z.string(),
  otherName: z.string(),
  artistId: z.string(),
  year: z.number(),
  tags: z.array(z.string()),
  featuringArtists: z.array(z.string()),
});
export type CreateAlbumBody = z.infer<typeof CreateAlbumBody>;

export const Page = z.object({
  page: z.number(),
  perPage: z.number(),
  totalItems: z.number(),
  totalPages: z.number(),
});
export type Page = z.infer<typeof Page>;

export const GetTracks = z.object({
  page: Page,
  tracks: z.array(Track),
});
export type GetTracks = z.infer<typeof GetTracks>;

export const GetTrackById = Track;
export type GetTrackById = z.infer<typeof GetTrackById>;

export const EditTrackBody = z.object({
  name: z.string().nullable().optional(),
  otherName: z.string().nullable().optional(),
  artistId: z.string().nullable().optional(),
  artistName: z.string().nullable().optional(),
  year: z.number().nullable().optional(),
  number: z.number().nullable().optional(),
  tags: z.array(z.string()).nullable().optional(),
  featuringArtists: z.array(z.string()).nullable().optional(),
});
export type EditTrackBody = z.infer<typeof EditTrackBody>;

export const UploadTrackBody = z.object({
  name: z.string(),
  otherName: z.string(),
  number: z.number(),
  year: z.number(),
  albumId: z.string(),
  artistId: z.string(),
  tags: z.array(z.string()),
  featuringArtists: z.array(z.string()),
});
export type UploadTrackBody = z.infer<typeof UploadTrackBody>;

export const Signup = z.object({
  id: z.string(),
  username: z.string(),
});
export type Signup = z.infer<typeof Signup>;

export const SignupBody = z.object({
  username: z.string(),
  password: z.string(),
  passwordConfirm: z.string(),
});
export type SignupBody = z.infer<typeof SignupBody>;

export const Signin = z.object({
  token: z.string(),
});
export type Signin = z.infer<typeof Signin>;

export const SigninBody = z.object({
  username: z.string(),
  password: z.string(),
});
export type SigninBody = z.infer<typeof SigninBody>;

export const ChangePasswordBody = z.object({
  currentPassword: z.string(),
  newPassword: z.string(),
  newPasswordConfirm: z.string(),
});
export type ChangePasswordBody = z.infer<typeof ChangePasswordBody>;

export const GetMe = z.object({
  id: z.string(),
  username: z.string(),
  role: z.string(),
  displayName: z.string(),
  quickPlaylist: z.string().nullable(),
});
export type GetMe = z.infer<typeof GetMe>;

export const Playlist = z.object({
  id: z.string(),
  name: z.string(),
});
export type Playlist = z.infer<typeof Playlist>;

export const GetPlaylists = z.object({
  playlists: z.array(Playlist),
});
export type GetPlaylists = z.infer<typeof GetPlaylists>;

export const CreatePlaylist = Playlist;
export type CreatePlaylist = z.infer<typeof CreatePlaylist>;

export const CreatePlaylistBody = z.object({
  name: z.string(),
});
export type CreatePlaylistBody = z.infer<typeof CreatePlaylistBody>;

export const PostPlaylistFilterBody = z.object({
  name: z.string(),
  filter: z.string(),
});
export type PostPlaylistFilterBody = z.infer<typeof PostPlaylistFilterBody>;

export const GetPlaylistById = z.object({
  id: z.string(),
  name: z.string(),
  items: z.array(Track),
});
export type GetPlaylistById = z.infer<typeof GetPlaylistById>;

export const AddItemToPlaylistBody = z.object({
  trackId: z.string(),
});
export type AddItemToPlaylistBody = z.infer<typeof AddItemToPlaylistBody>;

export const RemovePlaylistItemBody = z.object({
  trackId: z.string(),
});
export type RemovePlaylistItemBody = z.infer<typeof RemovePlaylistItemBody>;

export const GetSystemInfo = z.object({
  version: z.string(),
});
export type GetSystemInfo = z.infer<typeof GetSystemInfo>;

export const ExportArtist = z.object({
  id: z.string(),
  name: z.string(),
  picture: z.string(),
});
export type ExportArtist = z.infer<typeof ExportArtist>;

export const ExportAlbum = z.object({
  id: z.string(),
  name: z.string(),
  artistId: z.string(),
  coverArt: z.string(),
  year: z.number(),
});
export type ExportAlbum = z.infer<typeof ExportAlbum>;

export const ExportTrack = z.object({
  id: z.string(),
  name: z.string(),
  albumId: z.string(),
  artistId: z.string(),
  duration: z.number(),
  number: z.number(),
  year: z.number(),
  originalFilename: z.string(),
  mobileFilename: z.string(),
  created: z.number(),
  tags: z.array(z.string()),
});
export type ExportTrack = z.infer<typeof ExportTrack>;

export const Export = z.object({
  artists: z.array(ExportArtist),
  albums: z.array(ExportAlbum),
  tracks: z.array(ExportTrack),
});
export type Export = z.infer<typeof Export>;

export const Taglist = z.object({
  id: z.string(),
  name: z.string(),
  filter: z.string(),
  created: z.number(),
  updated: z.number(),
});
export type Taglist = z.infer<typeof Taglist>;

export const GetTaglists = z.object({
  taglists: z.array(Taglist),
});
export type GetTaglists = z.infer<typeof GetTaglists>;

export const GetTaglistById = Taglist;
export type GetTaglistById = z.infer<typeof GetTaglistById>;

export const GetTaglistTracks = z.object({
  page: Page,
  tracks: z.array(Track),
});
export type GetTaglistTracks = z.infer<typeof GetTaglistTracks>;

export const CreateTaglist = z.object({
  id: z.string(),
});
export type CreateTaglist = z.infer<typeof CreateTaglist>;

export const CreateTaglistBody = z.object({
  name: z.string(),
  filter: z.string(),
});
export type CreateTaglistBody = z.infer<typeof CreateTaglistBody>;

export const UpdateTaglistBody = z.object({
  name: z.string().nullable().optional(),
  filter: z.string().nullable().optional(),
});
export type UpdateTaglistBody = z.infer<typeof UpdateTaglistBody>;

export const UpdateUserSettingsBody = z.object({
  displayName: z.string().nullable().optional(),
  quickPlaylist: z.string().nullable().optional(),
});
export type UpdateUserSettingsBody = z.infer<typeof UpdateUserSettingsBody>;

export const TrackId = z.object({
  trackId: z.string(),
});
export type TrackId = z.infer<typeof TrackId>;

export const GetUserQuickPlaylistItemIds = z.object({
  trackIds: z.array(z.string()),
});
export type GetUserQuickPlaylistItemIds = z.infer<typeof GetUserQuickPlaylistItemIds>;

export const CreateApiToken = z.object({
  token: z.string(),
});
export type CreateApiToken = z.infer<typeof CreateApiToken>;

export const CreateApiTokenBody = z.object({
  name: z.string(),
});
export type CreateApiTokenBody = z.infer<typeof CreateApiTokenBody>;

export const ApiToken = z.object({
  id: z.string(),
  name: z.string(),
});
export type ApiToken = z.infer<typeof ApiToken>;

export const GetAllApiTokens = z.object({
  tokens: z.array(ApiToken),
});
export type GetAllApiTokens = z.infer<typeof GetAllApiTokens>;

