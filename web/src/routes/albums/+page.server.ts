import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query: Record<string, string> = {};
  const filter = url.searchParams.get("filter");
  if (filter) {
    query["filter"] = filter;
  }

  const sort = url.searchParams.get("sort");
  if (sort) {
    query["sort"] = sort;
  }

  const albums = await locals.apiClient.getAlbums({
    query,
  });
  if (!albums.success) {
    throw error(albums.error.code, { message: albums.error.message });
  }

  return {
    albums: albums.data.albums,
  };
};
