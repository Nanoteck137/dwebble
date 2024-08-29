import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions = {
  default: async ({ locals, cookies, request }) => {
    const formData = await request.formData();

    const username = formData.get("username")!;
    const password = formData.get("password")!;

    const res = await locals.apiClient.signin({
      username: username.toString(),
      password: password.toString(),
    });

    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
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
} satisfies Actions;
