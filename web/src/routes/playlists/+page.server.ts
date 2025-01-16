import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const playlists = await locals.apiClient.getPlaylists();
  if (!playlists.success) {
    throw error(playlists.error.code, { message: playlists.error.message });
  }

  return {
    playlists: playlists.data.playlists,
  };
};
