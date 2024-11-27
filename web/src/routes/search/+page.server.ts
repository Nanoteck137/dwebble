import type { Search } from "$lib/types";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ request, fetch }) => {
  const url = new URL(request.url);
  const query = url.searchParams.get("query") ?? "";

  const res = await fetch(`/server/search?query=${encodeURIComponent(query)}`);
  const d = (await res.json()) as Search;

  if (!d.success) {
    throw error(res.status, { message: d.message });
  }

  return {
    query,
    artists: d.artists,
    albums: d.albums,
    tracks: d.tracks,
  };
};
