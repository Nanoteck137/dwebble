import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, locals, url }) => {
  const playlist = await locals.apiClient.getPlaylistById(params.id);
  if (!playlist.success) {
    throw error(playlist.error.code, { message: playlist.error.message });
  }

  const query = getPagedQueryOptions(url.searchParams);
  const items = await locals.apiClient.getPlaylistItems(params.id, {
    query,
  });
  if (!items.success) {
    throw error(items.error.code, { message: items.error.message });
  }

  return {
    playlist: playlist.data,
    page: items.data.page,
    items: items.data.items,
  };
};
