import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const albums = await locals.apiClient.getAlbums();
  if (!albums.success) {
    throw error(albums.error.code, { message: albums.error.message });
  }

  return {
    albums: albums.data.albums,
  };
};
