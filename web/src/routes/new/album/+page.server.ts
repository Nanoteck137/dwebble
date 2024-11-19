import { error, fail, redirect } from "@sveltejs/kit";
import { z } from "zod";
import type { Actions } from "./$types";

// TODO(patrik): Set error messages
const ImportSchema = z.object({
  albumName: z.string().trim().min(1),
  artistName: z.string().trim().min(1),
});

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const formData = await request.formData();

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

    const res = await locals.apiClient.createAlbum({
      name: data.albumName,
      artist: data.artistName,
    });
    if (!res.success) {
      error(res.error.code, { message: res.error.message });
    }

    redirect(302, `/albums/${res.data.albumId}/edit`);
  },
};
