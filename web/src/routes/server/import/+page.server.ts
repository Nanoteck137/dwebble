import { PostAlbumImportBody } from "$lib/api/types";
import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  import: async ({ locals, request }) => {
    const formData = await request.formData();
    console.log(formData);

    const albumName = formData.get("albumName");
    if (!albumName) {
      throw error(400, "'albumName' needs to be set");
    }

    const artistName = formData.get("artistName");
    if (!artistName) {
      throw error(400, "'artistName' needs to be set");
    }

    const data: PostAlbumImportBody = {
      name: albumName.toString(),
      artist: artistName.toString(),
    };

    const body = new FormData();
    body.set("data", JSON.stringify(data));

    const coverArt = formData.get("coverArt");
    if (coverArt) {
      body.set("coverArt", coverArt);
    }

    const files = formData.getAll("files");
    files.forEach((f) => {
      const file = f as File;
      body.append("files", file);
    });

    const res = await locals.apiClient.importAlbum(body);
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/server/edit/album/${res.data.albumId}`);
  },
};
