import { isRoleAdmin } from "$lib/utils";
import { error, redirect } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  if (!locals.user || !isRoleAdmin(locals.user.role)) {
    throw redirect(301, `/artists/${params.id}`);
  }

  const res = await locals.apiClient.getArtistById(params.id);
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    artist: res.data,
  };
};
