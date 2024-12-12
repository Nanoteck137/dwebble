import { CreateApiTokenBody } from "$lib/api/types";
import { error, redirect } from "@sveltejs/kit";
import { fail, setError, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { assert, type Equals } from "tsafe";
import { z } from "zod";
import type { Actions, PageServerLoad } from "./$types";

const Body = CreateApiTokenBody;
const schema = Body.extend({
  name: Body.shape.name,
});

// eslint-disable-next-line @typescript-eslint/no-unused-expressions
assert<Equals<keyof z.infer<typeof Body>, keyof z.infer<typeof schema>>>;

export const load: PageServerLoad = async () => {
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

    const res = await locals.apiClient.createApiToken(form.data);
    if (!res.success) {
      if (res.error.type === "VALIDATION_ERROR") {
        const extra = res.error.extra as Record<
          keyof z.infer<typeof schema>,
          string | undefined
        >;

        setError(form, "name", extra.name ?? "");

        return fail(400, { form });
      } else {
        throw error(res.error.code, { message: res.error.message });
      }
    }

    throw redirect(302, "/account/tokens");
  },
};
