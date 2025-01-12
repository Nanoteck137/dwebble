import type { Album, Artist, Track } from "$lib/api/types";
import type { BareBoneError } from "$lib/utils";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ request, locals }) => {
  const url = new URL(request.url);
  const query = url.searchParams.get("query") ?? "";

  let artists = [] as Artist[];
  let artistError: BareBoneError | null = null;

  let albums = [] as Album[];
  let albumError: BareBoneError | null = null;

  let tracks = [] as Track[];
  let trackError: BareBoneError | null = null;

  const [artistQuery, albumQuery, trackQuery] = await Promise.all([
    locals.apiClient.searchArtists({ query: { query } }),
    locals.apiClient.searchAlbums({ query: { query } }),
    locals.apiClient.searchTracks({ query: { query } }),
  ]);

  if (!artistQuery.success) {
    artistError = artistQuery.error;
  } else {
    artists = artistQuery.data.artists;
  }

  if (!albumQuery.success) {
    albumError = albumQuery.error;
  } else {
    albums = albumQuery.data.albums;
  }

  if (!trackQuery.success) {
    trackError = trackQuery.error;
  } else {
    tracks = trackQuery.data.tracks;
  }

  return {
    query,

    artistError,
    artists,

    albumError,
    albums,

    trackError,
    tracks,
  };
};
