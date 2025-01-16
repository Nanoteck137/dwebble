import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const res = await locals.apiClient.getDetailedTrackById(params.trackId);
  if (!res.success) {
    throw error(res.error.code, { message: res.error.message });
  }

  return {
    track: res.data,
  };
};
