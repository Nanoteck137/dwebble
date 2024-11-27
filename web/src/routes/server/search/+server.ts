import type { Search } from "$lib/types";
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ locals, request }) => {
  const url = new URL(request.url);

  const query = url.searchParams.get("query") || "";

  const [artists, albums, tracks] = await Promise.all([
    locals.apiClient.searchArtists({
      query: { query: query },
    }),
    locals.apiClient.searchAlbums({
      query: { query: query },
    }),
    locals.apiClient.searchTracks({
      query: { query: query },
    }),
  ]);

  if (!artists.success) {
    return json({
      message: artists.error.message,
      success: false,
    } as Search);
  }

  if (!albums.success) {
    return json({
      message: albums.error.message,
      success: false,
    } as Search);
  }

  if (!tracks.success) {
    return json({
      message: tracks.error.message,
      success: false,
    } as Search);
  }

  return json({
    success: true,
    artists: artists.data.artists,
    albums: albums.data.albums,
    tracks: tracks.data.tracks,
  } as Search);
};
