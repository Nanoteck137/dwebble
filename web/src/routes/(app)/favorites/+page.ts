import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, url }) => {
  const data = await parent();

  if (!data.user) {
    throw error(401, { message: "Not authenticated" });
  }

  const query = getPagedQueryOptions(url.searchParams);

  const favorites = await data.apiClient.getUserTrackFavoritesById(
    data.user.id,
    { query },
  );
  if (!favorites.success) {
    throw error(favorites.error.code, { message: favorites.error.message });
  }

  return {
    ...data,
    page: favorites.data.page,
    tracks: favorites.data.items,
  };
};
