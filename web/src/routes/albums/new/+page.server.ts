import { CreateAlbumBody } from "$lib/api/types";
import { isRoleAdmin } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import { fail, setError, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { assert, type Equals } from "tsafe";
import type { z } from "zod";
import type { Actions, PageServerLoad } from "./$types";

const Body = CreateAlbumBody;
const schema = Body.extend({
  name: Body.shape.name,
  artistId: Body.shape.artistId,
});

// eslint-disable-next-line @typescript-eslint/no-unused-expressions
assert<Equals<keyof z.infer<typeof Body>, keyof z.infer<typeof schema>>>;

export const load: PageServerLoad = async ({ locals }) => {
  if (!locals.user || !isRoleAdmin(locals.user.role)) {
    throw redirect(301, "/albums");
  }

  const form = await superValidate(zod(schema));

  return {
    form,
  };
};

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const form = await superValidate(request, zod(schema));

    if (!form.valid) {
      return fail(400, { form });
    }

    const res = await locals.apiClient.createAlbum(form.data);
    if (!res.success) {
      if (res.error.type === "VALIDATION_ERROR") {
        const extra = res.error.extra as Record<
          keyof z.infer<typeof schema>,
          string | undefined
        >;

        setError(form, "name", extra.artistId ?? "");
        setError(form, "artistId", extra.artistId ?? "");

        return fail(400, { form });
      } else {
        throw error(res.error.code, { message: res.error.message });
      }
    }

    throw redirect(302, `/albums/${res.data.albumId}/edit`);
  },
};
