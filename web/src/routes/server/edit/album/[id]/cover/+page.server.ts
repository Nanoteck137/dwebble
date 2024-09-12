import { error } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request, params }) => {
    const formData = await request.formData();
    console.log(formData);

    const f = formData.get("file");
    if (f === null) {
      throw error(400, "Missing 'file'");
    }
    const file = f as File;

    if (file.size === 0) {
      throw error(400, "No file selected");
    }
  },
};
