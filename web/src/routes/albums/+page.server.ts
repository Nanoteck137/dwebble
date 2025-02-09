import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const albums = await locals.apiClient.getAlbums({
    query,
  });
  if (!albums.success) {
    throw error(albums.error.code, { message: albums.error.message });
  }

  return {
    page: albums.data.page,
    albums: albums.data.albums,
  };
};
