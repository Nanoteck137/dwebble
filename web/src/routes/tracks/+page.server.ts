import type { Page, Track } from "$lib/api/types";
import { error, redirect } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url, fetch }) => {
  const query: Record<string, string> = {};
  const filter = url.searchParams.get("filter");
  if (filter) {
    query["filter"] = filter;
  }

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
        page: {} as Page,
        tracks: [] as Track[],
        filter,
        sort,
        filterError: tracks.error.message,
      };
    }

    if (tracks.error.type === "INVALID_SORT") {
      return {
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
    page: tracks.data.page,
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
  newPlaylist: async ({ locals, request }) => {
    const formData = await request.formData();

    const filter = formData.get("filter");
    if (filter === null) {
      throw error(500, "'filter' not set");
    }
    console.log("Filter", filter);

    const sort = formData.get("sort");
    if (sort === null) {
      throw error(500, "'sort' not set");
    }

    console.log("Sort", sort);

    const res = await locals.apiClient.createPlaylistFromFilter({
      name: "Generated Playlist",
      filter: filter.toString(),
      sort: sort.toString(),
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/playlists/${res.data.id}`);
  },
};
