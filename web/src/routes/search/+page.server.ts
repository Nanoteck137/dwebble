import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ request, locals }) => {
  const url = new URL(request.url);
  const query = url.searchParams.get("query") ?? "";

  // TODO(patrik): Use Promise.all
  const artists = await locals.apiClient.searchArtists({ query: { query } });
  if (!artists.success) {
    // TODO(patrik): Should we throw the error?
    throw error(artists.error.code, { message: artists.error.message });
  }

  const albums = await locals.apiClient.searchAlbums({ query: { query } });
  if (!albums.success) {
    // TODO(patrik): Should we throw the error?
    throw error(albums.error.code, { message: albums.error.message });
  }

  const tracks = await locals.apiClient.searchTracks({ query: { query } });
  if (!tracks.success) {
    // TODO(patrik): Should we throw the error?
    throw error(tracks.error.code, { message: tracks.error.message });
  }

  return {
    query,
    artists: artists.data.artists,
    albums: albums.data.albums,
    tracks: tracks.data.tracks,
  };
};
