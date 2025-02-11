import type { Artist, Page } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const artists = await locals.apiClient.getArtists({
    query,
  });
  if (!artists.success) {
    // TODO(patrik): Fix this
    if (artists.error.type === "INVALID_FILTER") {
      return {
        page: {} as Page,
        artists: [] as Artist[],
        filter: query["filter"],
        sort: query["sort"],
        filterError: artists.error.message,
      };
    }

    if (artists.error.type === "INVALID_SORT") {
      return {
        page: {} as Page,
        artists: [] as Artist[],
        filter: query["filter"],
        sort: query["sort"],
        sortError: artists.error.message,
      };
    }

    throw error(artists.error.code, artists.error.message);
  }

  return {
    page: artists.data.page,
    artists: artists.data.artists,
    filter: query["filter"],
    sort: query["sort"],
  };
};
