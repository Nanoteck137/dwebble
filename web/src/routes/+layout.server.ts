import type { Playlist } from "$lib/api/types";
import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
  const user = locals.user;

  let quickPlaylistIds = [] as string[];
  let playlists: Playlist[] | null = null;

  if (user && user.quickPlaylist) {
    {
      const res = await locals.apiClient.getUserQuickPlaylistItemIds();
      if (!res.success) {
        throw error(res.error.code, { message: res.error.message });
      }

      quickPlaylistIds = res.data.trackIds;
    }

    {
      const res = await locals.apiClient.getPlaylists();
      if (!res.success) {
        throw error(res.error.code, { message: res.error.message });
      }

      playlists = res.data.playlists;
    }
  }

  return {
    apiAddress: locals.apiAddress,
    userToken: locals.token,
    user,
    quickPlaylistIds,
    userPlaylists: playlists,
  };
};
