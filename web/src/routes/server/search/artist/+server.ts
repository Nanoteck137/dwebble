import type { Search } from "$lib/types";
import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ locals, request }) => {
  const url = new URL(request.url);

  const query = url.searchParams.get("query") || "";

  const res = await locals.apiClient.searchArtists({
    query: { query: query },
  });

  if (!res.success) {
    return json({
      message: res.error.message,
      success: false,
    } as Search);
  }

  return json({
    success: true,
    artists: res.data.artists,
  } as Search);
};
