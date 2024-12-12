import { error, redirect } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const res = await locals.apiClient.getPlaylists();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    playlists: res.data.playlists,
  };
};

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const formData = await request.formData();

    const playlistId = formData.get("playlistId");
    console.log(playlistId);
    if (playlistId === null) {
      throw error(400, "'playlistId' not set");
    }

    const res = await locals.apiClient.updateUserSettings({
      quickPlaylist: playlistId.toString(),
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, "/account");
  },
};
