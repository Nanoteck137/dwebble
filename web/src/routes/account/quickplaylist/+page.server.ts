import { error, redirect } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const res = await locals.apiClient.getPlaylists();
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    playlists: res.data.playlists,
  };
};

export const actions: Actions = {
  default: async () => {
    throw redirect(302, "/account");
  },
};
