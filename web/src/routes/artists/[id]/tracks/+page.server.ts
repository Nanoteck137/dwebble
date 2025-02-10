import { getPagedQueryOptions } from "$lib/utils";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, url, params }) => {
  const query = getPagedQueryOptions(url.searchParams);

  const tracks = await locals.apiClient.getTracks({
    query: {
      ...query,
      filter: `artistId == "${params.id}" || hasFeaturingArtist("${params.id}")`,
    },
  });
  if (!tracks.success) {
    throw error(tracks.error.code, {
      message: tracks.error.message,
    });
  }

  return {
    page: tracks.data.page,
    tracks: tracks.data.tracks,
  };
};
