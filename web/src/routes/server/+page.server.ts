import { error } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  return {};
};

export const actions: Actions = {
  runExport: async ({ locals }) => {
    const res = await locals.apiClient.systemExport();
    if (!res.success) {
      throw error(res.error.code, res.error.message);
    }
  },
};
