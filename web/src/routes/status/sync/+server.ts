import { json } from "@sveltejs/kit";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ locals }) => {
  const res = await locals.apiClient.getSyncStatus();
  return json(res);
};
