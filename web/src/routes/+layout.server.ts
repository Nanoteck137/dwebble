import { error } from "@sveltejs/kit";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
  const user = locals.user;

  let quickPlaylistIds = [] as string[];

  if (user && user.quickPlaylist) {
    const res = await locals.apiClient.getUserQuickPlaylistItemIds();
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    quickPlaylistIds = res.data.trackIds;
  }

  return {
    user,
    quickPlaylistIds,
  };
};
