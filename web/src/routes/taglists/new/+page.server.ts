import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request }) => {
    const formData = await request.formData();

    const taglistName = formData.get("taglistName");
    if (!taglistName) {
      throw error(400, { message: "taglistName is not set" });
    }

    const taglistFilter = formData.get("taglistFilter");
    if (!taglistFilter) {
      throw error(400, { message: "taglistFilter is not set" });
    }

    const res = await locals.apiClient.createTaglist({
      name: taglistName.toString(),
      filter: taglistFilter.toString(),
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/taglists/${res.data.id}`);
  },
};
