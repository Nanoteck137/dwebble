import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const formData = await request.formData();

    const currentPassword = formData.get("currentPassword");
    if (currentPassword === null) {
      throw error(400, "'currentPassword' not set");
    }

    const newPassword = formData.get("newPassword");
    if (newPassword === null) {
      throw error(400, "'newPassword' not set");
    }

    const newPasswordConfirm = formData.get("newPasswordConfirm");
    if (newPasswordConfirm === null) {
      throw error(400, "'newPasswordConfirm' not set");
    }

    const res = await locals.apiClient.changePassword({
      currentPassword: currentPassword.toString(),
      newPassword: newPassword.toString(),
      newPasswordConfirm: newPasswordConfirm.toString(),
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, "/account");
  },
};
