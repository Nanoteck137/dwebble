import type { TrackFilter } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ parent, params, url }) => {
  const data = await parent();

  let filters: TrackFilter[] | null = null;
  const res = await data.apiClient.getTrackFilters();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }
  filters = res.data.filters;

  const query = getPagedQueryOptions(url.searchParams);

  const filterId = url.searchParams.get("filterId");
  if (filterId) {
    query["filterId"] = filterId;
  }

  const favorites = await data.apiClient.getUserTrackFavoritesById(params.id, {
    query,
  });
  if (!favorites.success) {
    throw error(favorites.error.code, { message: favorites.error.message });
  }

  return {
    ...data,
    filters,
    page: favorites.data.page,
    tracks: favorites.data.items,
  };
};
