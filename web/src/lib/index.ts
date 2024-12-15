// place files you want to import through the `$lib` alias in this folder.

import type { ApiClient } from "$lib/api/client";
import { writable, type Writable } from "svelte/store";

export type Artist = {
  id: string;
  name: string;
};

type GetApiClient = () => ApiClient;

export function artistQuery(getApiClient: GetApiClient) {
  const open = writable(false);

  const artist: Writable<Artist | undefined> = writable();

  const currentQuery = writable("");
  const queryResults = writable([] as Artist[]);

  open.subscribe((v) => {
    if (v) {
      queryResults.set([]);
      currentQuery.set("");
    }
  });

  let timer: NodeJS.Timeout;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;

    queryResults.set([]);
    currentQuery.set(current);

    clearTimeout(timer);
    timer = setTimeout(async () => {
      const apiClient = getApiClient();
      const res = await apiClient.searchArtists({
        query: {
          query: current,
        },
      });

      if (res.success) {
        queryResults.set(
          res.data.artists.map((artist) => ({
            id: artist.id,
            name: artist.name,
          })),
        );
      }
    }, 500);
  }

  return {
    artist,
    open,
    onInput,
    queryResults,
    currentQuery,
  };
}
