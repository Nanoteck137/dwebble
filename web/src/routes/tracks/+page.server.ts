import { error } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

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

  const tracks = await locals.apiClient.getTracks({ query });
  if (!tracks.success) {
    // TODO(patrik): Fix this
    if (tracks.error.type === "INVALID_FILTER") {
      return {
        tracks: [],
        filter,
        sort,
        filterError: tracks.error.message,
      };
    }

    if (tracks.error.type === "INVALID_SORT") {
      return {
        tracks: [],
        filter,
        sort,
        sortError: tracks.error.message,
      };
    }
    throw error(tracks.error.code, tracks.error.message);
  }

  return {
    tracks: tracks.data.tracks,
    filter,
    sort,
  };
};

export const actions: Actions = {
  quickAddToPlaylist: async ({ request, locals, cookies }) => {
    const formData = await request.formData();

    const trackId = formData.get("trackId");
    if (!trackId) {
      throw error(400, "No track id set");
    }

    const playlistId = cookies.get("quick-playlist");
    if (!playlistId) {
      throw error(400, "No quick playlist set");
    }

    const res = await locals.apiClient.addItemsToPlaylist(playlistId, {
      tracks: [trackId.toString()],
    });
    if (!res.success) {
      throw error(res.error.code, res.error.message);
    }
  },
};
