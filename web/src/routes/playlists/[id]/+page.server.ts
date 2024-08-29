import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, locals }) => {
  const playlist = await locals.apiClient.getPlaylistById(params.id);
  if (!playlist.success) {
    throw error(playlist.error.code, { message: playlist.error.message });
  }

  return {
    ...playlist.data,
  };
};
