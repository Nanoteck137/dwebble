import type { PostAlbumImportBody } from "$lib/api/types";
import { error, fail, redirect } from "@sveltejs/kit";
import { z } from "zod";
import type { Actions } from "./$types";

const ImportSchema = z.object({
  albumName: z.string().min(1),
  artistName: z.string().min(1),
});

export const actions: Actions = {
  import: async ({ locals, request }) => {
    const formData = await request.formData();
    console.log(formData);

    const parsed = await ImportSchema.safeParseAsync(
      Object.fromEntries(formData),
    );

    if (!parsed.success) {
      const flatten = parsed.error.flatten();
      return fail(400, {
        errors: flatten.fieldErrors,
      });
    }

    const data = parsed.data;

    const bodyData: PostAlbumImportBody = {
      name: data.albumName,
      artist: data.artistName,
    };

    const body = new FormData();
    body.set("data", JSON.stringify(bodyData));

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
      error(res.error.code, { message: res.error.message });
    }

    redirect(302, `/server/edit/album/${res.data.albumId}`);
  },
};
