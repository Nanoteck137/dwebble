import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const artists = await locals.apiClient.getArtists({
    query,
  });
  if (!artists.success) {
    throw error(artists.error.code, artists.error.message);
  }

  return {
    page: artists.data.page,
    artists: artists.data.artists,
  };
};
