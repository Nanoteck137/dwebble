import { createApiResponse } from "$lib/api/base-client";
import { GetTracks } from "$lib/api/types";
import { z } from "zod";

type FetchOptions = {
  filter?: string;
  sort?: string;
  perPage?: number;
  page?: number;
};

const Schema = createApiResponse(GetTracks, z.undefined());
type T = z.infer<typeof Schema>;

export async function getTracks(opts: FetchOptions, fetchFunc: typeof fetch) {
  const query = new URLSearchParams();
  if (opts.filter) {
    query.set("filter", opts.filter);
  }

  if (opts.sort) {
    query.set("sort", opts.sort);
  }

  if (opts.perPage) {
    query.set("perPage", opts.perPage.toString());
  }

  if (opts.page) {
    query.set("page", opts.page.toString());
  }

  const res = await fetchFunc(`/server/tracks?${query.toString()}`);
  const tracks = (await res.json()) as T;

  return tracks;
}
