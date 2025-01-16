import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, locals }) => {
  const playlist = await locals.apiClient.getPlaylistById(params.id);
  if (!playlist.success) {
    throw error(playlist.error.code, { message: playlist.error.message });
  }

  const items = await locals.apiClient.getPlaylistItems(params.id);
  if (!items.success) {
    throw error(items.error.code, { message: items.error.message });
  }

  return {
    playlist: playlist.data,
    items: items.data.items,
  };
};
