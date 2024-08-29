import { error } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const playlists = await locals.apiClient.getPlaylists();
  if (!playlists.success) {
    throw error(playlists.error.code, playlists.error.message);
  }

  return {
    playlists: playlists.data.playlists,
  };
};

export const actions: Actions = {
  setQuickPlaylist: async ({ locals, request, cookies }) => {
    const formData = await request.formData();

    const playlistId = formData.get("playlistId");
    if (!playlistId) {
      throw error(400, "Playlist Id not set");
    }

    const playlist = await locals.apiClient.getPlaylistById(
      playlistId.toString(),
    );
    if (!playlist.success) {
      throw error(playlist.error.code, playlist.error.message);
    }

    cookies.set("quick-playlist", playlist.data.id, { path: "/" });
  },
};
