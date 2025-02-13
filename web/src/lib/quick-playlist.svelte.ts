import { handleApiError } from "$lib";
import type { ApiClient } from "$lib/api/client";
import { getContext, setContext } from "svelte";
import toast from "svelte-5-french-toast";

class QuickPlaylist {
  apiClient: ApiClient;
  playlistId = $state<string>("");
  ids = $state<string[]>([]);

  constructor(apiClient: ApiClient, playlistId: string, ids: string[]) {
    this.apiClient = apiClient;
    this.playlistId = playlistId;
    this.ids = ids;
  }

  async toggleTrack(trackId: string) {
    if (this.playlistId === "") {
      toast.error("No playlist id set for quick playlist");
      return;
    }

    if (this.hasTrack(trackId)) {
      const res = await this.apiClient.removeItemFromUserQuickPlaylist({
        trackId: trackId,
      });

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    } else {
      const res = await this.apiClient.addToUserQuickPlaylist({
        trackId: trackId,
      });

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    }

    const ids = await this.apiClient.getUserQuickPlaylistItemIds();
    if (!ids.success) {
      handleApiError(ids.error);
      return;
    }

    this.ids = ids.data.trackIds;
  }

  hasTrack(trackId: string) {
    return !!this.ids.find((v) => v === trackId);
  }
}

const QUICK_PLAYLIST_KEY = Symbol("QUICK_PLAYLIST");

export function setQuickPlaylist(
  apiClient: ApiClient,
  playlistId: string,
  ids: string[],
) {
  return setContext(
    QUICK_PLAYLIST_KEY,
    new QuickPlaylist(apiClient, playlistId, ids),
  );
}

export function getQuickPlaylist() {
  return getContext<ReturnType<typeof setQuickPlaylist>>(QUICK_PLAYLIST_KEY);
}
