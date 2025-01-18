// place files you want to import through the `$lib` alias in this folder.

import { ApiClient } from "$lib/api/client";
import AddToPlaylistModal, {
  type Props as AddToPlaylistModalProps,
} from "$lib/components/modals/AddToPlaylistModal.svelte";
import ConfirmModal, {
  type Props as ConfirmModalProps,
} from "$lib/components/modals/ConfirmModal.svelte";
import InputModal, {
  type Props as InputModalProps,
} from "$lib/components/modals/InputModal.svelte";
import QueryArtistModal, {
  type Props as QueryArtistModalProps,
} from "$lib/components/modals/QueryArtistModal.svelte";
import type { UIArtist } from "$lib/types";
import { modals } from "svelte-modals";
import { writable, type Writable } from "svelte/store";

export function openAddToPlaylist(props: AddToPlaylistModalProps) {
  return modals.open(AddToPlaylistModal, props);
}

export function openConfirm(props: ConfirmModalProps) {
  return modals.open(ConfirmModal, props);
}

export function openArtistQuery(props: QueryArtistModalProps) {
  return modals.open(QueryArtistModal, props);
}

export function openInput(props: InputModalProps) {
  return modals.open(InputModal, props);
}

export function createApiClient(data: {
  apiAddress: string;
  userToken?: string;
}) {
  const apiClient = new ApiClient(data.apiAddress);
  apiClient.setToken(data.userToken);
  return apiClient;
}

export function isInQuickPlaylist(
  data: { quickPlaylistIds: string[] },
  trackId: string,
) {
  if (!data.quickPlaylistIds) return false;
  return !!data.quickPlaylistIds.find((v) => v === trackId);
}

type GetApiClient = () => ApiClient;

export function artistQuery(getApiClient: GetApiClient) {
  const open = writable(false);

  const artist: Writable<UIArtist | undefined> = writable();

  const currentQuery = writable("");
  const queryResults = writable([] as UIArtist[]);

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
            name: artist.name.default,
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
