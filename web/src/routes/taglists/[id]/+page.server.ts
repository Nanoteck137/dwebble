import type { Page, Track } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url, parent }) => {
  const data = await parent();

  const query: Record<string, string> = {};
  const filter = url.searchParams.get("filter");
  query["filter"] = data.taglist.filter;

  const sort = url.searchParams.get("sort");
  if (sort) {
    query["sort"] = sort;
  }

  const page = url.searchParams.get("page");
  if (page) {
    query["page"] = page;
  }

  const tracks = await locals.apiClient.getTracks({
    query,
  });

  if (!tracks.success) {
    // TODO(patrik): Fix this
    if (tracks.error.type === "INVALID_FILTER") {
      return {
        ...data,
        page: {} as Page,
        tracks: [] as Track[],
        filter,
        sort,
        filterError: tracks.error.message,
      };
    }

    if (tracks.error.type === "INVALID_SORT") {
      return {
        ...data,
        page: {} as Page,
        tracks: [] as Track[],
        filter,
        sort,
        sortError: tracks.error.message,
      };
    }
    throw error(tracks.error.code, tracks.error.message);
  }

  return {
    ...data,
    page: tracks.data.page,
    tracks: tracks.data.tracks,
    filter,
    sort,
  };
};
