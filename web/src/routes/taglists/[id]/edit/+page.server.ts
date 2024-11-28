import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ request, locals, params }) => {
    const formData = await request.formData();

    const taglistName = formData.get("taglistName");
    if (taglistName === null) {
      throw error(500, "'taglistName' not set");
    }

    const taglistFilter = formData.get("taglistFilter");
    if (taglistFilter === null) {
      throw error(500, "'taglistFilter' not set");
    }

    const res = await locals.apiClient.updateTaglist(params.id, {
      name: taglistName.toString(),
      filter: taglistFilter.toString(),
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/taglists/${params.id}`);
  },
};
