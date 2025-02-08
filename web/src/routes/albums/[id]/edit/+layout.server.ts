import { isRoleAdmin } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  if (!locals.user || !isRoleAdmin(locals.user.role)) {
    throw redirect(301, `/albums/${params.id}`);
  }

  const album = await locals.apiClient.getAlbumById(params.id);
  if (!album.success) {
    throw error(album.error.code, { message: album.error.message });
  }

  const tracks = await locals.apiClient.getDetailedTracks({
    query: {
      filter: `albumId == "${params.id}"`,
      sort: "sort=number",
    },
  });
  if (!tracks.success) {
    throw error(tracks.error.code, { message: tracks.error.message });
  }

  // const tracks = await locals.apiClient.getAlbumTracks(params.id);
  // if (!tracks.success) {
  //   throw error(tracks.error.code, { message: tracks.error.message });
  // }

  return {
    album: album.data,
    tracks: tracks.data.tracks,
  };
};
