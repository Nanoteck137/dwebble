import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const formData = await request.formData();

    const playlistName = formData.get("playlistName");
    if (playlistName === null) {
      throw error(400, { message: "'playlistName' not set" });
    }

    const res = await locals.apiClient.createPlaylist({
      name: playlistName.toString(),
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/playlists/${res.data.id}`);
  },
};
