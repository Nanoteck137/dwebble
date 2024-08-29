import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, locals }) => {
  const album = await locals.apiClient.getAlbumById(params.id);
  if (!album.success) {
    throw error(album.error.code, {
      message: album.error.message,
      type: album.error.type,
    });
  }

  const tracks = await locals.apiClient.getAlbumTracks(params.id);
  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
      type: tracks.error.type,
    });
  }

  return {
    album: album.data,
    tracks: tracks.data.tracks,
  };
};
