import type { Page, Track } from "$lib/api/types";
import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const tracks = await locals.apiClient.getTracks({
    query,
  });

  if (!tracks.success) {
    // TODO(patrik): Fix this
    if (tracks.error.type === "INVALID_FILTER") {
      return {
        page: {} as Page,
        tracks: [] as Track[],
        filter: query["filter"],
        sort: query["sort"],
        filterError: tracks.error.message,
      };
    }

    if (tracks.error.type === "INVALID_SORT") {
      return {
        page: {} as Page,
        tracks: [] as Track[],
        filter: query["filter"],
        sort: query["sort"],
        sortError: tracks.error.message,
      };
    }

    throw error(tracks.error.code, tracks.error.message);
  }

  return {
    page: tracks.data.page,
    tracks: tracks.data.tracks,
    filter: query["filter"],
    sort: query["sort"],
  };
};
