import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request, params }) => {
    const formData = await request.formData();
    console.log(formData);

    const albumId = formData.get("albumId");
    if (albumId === null) {
      throw error(400, "Missing 'albumId'");
    }

    const f = formData.get("file");
    if (f === null) {
      throw error(400, "Missing 'file'");
    }
    const file = f as File;

    if (file.size === 0) {
      throw error(400, "No file selected");
    }

    const data = new FormData();
    data.set("cover", file);

    const res = await locals.apiClient.changeAlbumCover(
      albumId.toString(),
      data,
    );
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/albums/${params.id}/edit`);
  },
};
