import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const res = await locals.apiClient.getArtistById(params.id);
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  // TODO(patrik): Use Promise.all
  const albums = await locals.apiClient.getAlbums({
    query: {
      filter: `artistId == "${params.id}" || hasFeaturingArtist("${params.id}")`,
    },
  });
  if (!albums.success) {
    throw error(albums.error.code, {
      message: albums.error.message,
    });
  }

  const tracks = await locals.apiClient.getTracks({
    query: {
      filter: `artistId == "${params.id}" || hasFeaturingArtist("${params.id}")`,
    },
  });
  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
    });
  }

  return {
    artist: res.data,
    albums: albums.data.albums,
    trackPage: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
