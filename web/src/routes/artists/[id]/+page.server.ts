import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  // TODO(patrik): Use Promise.all
  const albums = await locals.apiClient.getAlbums({
    query: {
      filter: `artistId == "${params.id}" || hasFeaturingArtist("${params.id}")`,
      perPage: "5",
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
      perPage: "5",
    },
  });
  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
    });
  }

  return {
    albums: albums.data.albums,
    trackPage: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
