import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  deleteAlbum: async ({ locals, request }) => {
    const formData = await request.formData();
    const albumId = formData.get("albumId");
    if (!albumId) {
      throw error(400, { message: "albumId is not set" });
    }

    const res = await locals.apiClient.deleteAlbum(albumId.toString());
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(301, `/server/edit/album`);
  },

  deleteTrack: async ({ locals, request }) => {
    const formData = await request.formData();

    const trackId = formData.get("trackId");
    if (!trackId) {
      throw error(400, { message: "trackId is not set" });
    }

    const res = await locals.apiClient.removeTrack(trackId.toString());
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },
};
