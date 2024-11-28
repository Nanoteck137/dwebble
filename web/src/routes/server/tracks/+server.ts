import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ locals, request }) => {
  const url = new URL(request.url);

  const query: Record<string, string> = {};

  url.searchParams.forEach((v, k) => {
    query[k] = v;
  });

  // TODO(patrik): We could do a raw fetch to get the data and then
  // let the client parse the api response
  const res = await locals.apiClient.getTracks({
    query,
  });

  if (!res.success) {
    return json(res, { status: res.error.code });
  }

  return json(res);
};
