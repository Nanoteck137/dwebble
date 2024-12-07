import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request, params }) => {
    const formData = await request.formData();

    const name = formData.get("name");
    if (!name) {
      throw error(400, { message: "name is not set" });
    }

    const year = formData.get("year");
    if (!year) {
      throw error(400, { message: "year is not set" });
    }

    const artist = formData.get("artist");
    if (!artist) {
      throw error(400, { message: "artist is not set" });
    }

    const res = await locals.apiClient.editAlbum(params.id, {
      name: name.toString(),
      artistId: null,
      artistName: artist.toString(),
      year: year ? parseInt(year.toString()) : null,
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/albums/${params.id}/edit`);
  },
};
