import { SigninBody } from "$lib/api/types";
import { capitilize } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import { fail, setError, superValidate } from "sveltekit-superforms";
import { zod } from "sveltekit-superforms/adapters";
import { assert, type Equals } from "tsafe";
import type { z } from "zod";
import type { Actions, PageServerLoad } from "./$types";

const Body = SigninBody;
const schema = Body.extend({
  username: Body.shape.username,
  password: Body.shape.password,
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
  default: async ({ locals, request, cookies }) => {
    const form = await superValidate(request, zod(schema));

    if (!form.valid) {
      return fail(400, { form });
    }

    const res = await locals.apiClient.signin(form.data);
    if (!res.success) {
      switch (res.error.type) {
        case "VALIDATION_ERROR": {
          const extra = res.error.extra as Record<
            keyof z.infer<typeof schema>,
            string | undefined
          >;

          setError(form, "username", capitilize(extra.username ?? ""));
          setError(form, "password", capitilize(extra.password ?? ""));

          return fail(400, { form });
        }
        case "USER_NOT_FOUND": {
          setError(form, "username", "The user does not exist");
          return fail(400, { form });
        }
        case "INVALID_CREDENTIALS": {
          setError(form, "password", "Invalid credentials");
          return fail(400, { form });
        }
        default:
          throw error(res.error.code, { message: res.error.message });
      }
    }

    locals.apiClient.setToken(res.data.token);
    const user = await locals.apiClient.getMe();
    if (!user.success) {
      throw error(user.error.code, { message: user.error.message });
    }

    const data = {
      token: res.data.token,
      user: {
        id: user.data.id,
        username: user.data.username,
      },
    };

    cookies.set("auth", JSON.stringify(data), {
      path: "/",
      sameSite: "strict",
    });

    throw redirect(302, "/");
  },
};
