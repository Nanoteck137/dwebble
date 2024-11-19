import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  const album = await locals.apiClient.getAlbumById(params.id);
  if (!album.success) {
    throw error(album.error.code, { message: album.error.message });
  }

  const tracks = await locals.apiClient.getAlbumTracks(params.id);
  if (!tracks.success) {
    throw error(tracks.error.code, { message: tracks.error.message });
  }

  return {
    album: album.data,
    tracks: tracks.data.tracks,
  };
};
