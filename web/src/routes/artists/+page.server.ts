import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const artists = await locals.apiClient.getArtists();
  if (!artists.success) {
    throw error(artists.error.code, artists.error.message);
  }

  return {
    artists: artists.data.artists,
  };
};
