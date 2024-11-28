import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals, params }) => {
  const res = await locals.apiClient.getTaglistById(params.id);
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    taglist: res.data,
  };
};
