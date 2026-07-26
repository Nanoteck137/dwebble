import { handleApiError } from "$lib";
import type { ApiClient } from "$lib/api/client";
import * as api from "$lib/api/types";
import { getContext, setContext } from "svelte";

export type MediaRef = {
  id: string;
  name: string;
};

export type MediaItem = {
  queueItemId: string;
  trackId: string;
  name: string;
  album: MediaRef;
  artists: MediaRef[];
  coverArt: string;
};

type QueueEntry = {
  queueItemId: string;
  trackId: string;
};

const QUEUE_PAGE_SIZE = 50;

function trackToMediaItem(item: api.QueueItem): MediaItem {
  return {
    queueItemId: item.queueItemId,
    trackId: item.track.id,
    name: item.track.name,
    album: { id: item.track.albumId, name: item.track.albumName },
    artists: item.track.artists.map((a) => ({ id: a.id, name: a.name })),
    coverArt: item.track.coverArt.small,
  };
}

class Queue {
  entries = $state<QueueEntry[]>([]);
  index = $state(0);
  totalItems = $state(0);

  private apiClient: ApiClient;
  private queueId: string;
  private loadedItems = $state(new Map<number, MediaItem>());

  constructor(apiClient: ApiClient, queueId: string) {
    this.apiClient = apiClient;
    this.queueId = queueId;
  }

  setQueue(entries: QueueEntry[], index: number, totalItems: number) {
    this.entries = entries;
    this.index = index;
    this.totalItems = totalItems;
    this.loadedItems.clear();
  }

  async loadEntries() {
    const res = await this.apiClient.getQueueIds(this.queueId);
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    this.setQueue(
      res.data.items,
      res.data.currentIndex,
      res.data.items.length,
    );
  }

  async loadItem(index: number): Promise<MediaItem | null> {
    if (index < 0 || index >= this.totalItems) return null;

    const cached = this.loadedItems.get(index);
    if (cached) return cached;

    const res = await this.apiClient.getQueueItemAtIndex(
      this.queueId,
      index.toString(),
    );
    if (!res.success) {
      handleApiError(res.error);
      return null;
    }

    const mediaItem = trackToMediaItem(res.data.item);
    this.loadedItems.set(index, mediaItem);
    return mediaItem;
  }

  async loadItems(indices: number[]): Promise<Map<number, MediaItem>> {
    if (indices.length === 0) return new Map();

    const missing = indices.filter((i) => !this.loadedItems.has(i));

    if (missing.length > 0) {
      const pages = new Map<number, number[]>();
      for (const index of missing) {
        const page = Math.floor(index / QUEUE_PAGE_SIZE);
        if (!pages.has(page)) pages.set(page, []);
        pages.get(page)!.push(index);
      }

      for (const page of pages.keys()) {
        const res = await this.apiClient.getQueue(this.queueId, {
          query: {
            page: page.toString(),
            perPage: QUEUE_PAGE_SIZE.toString(),
          },
        });
        if (!res.success) {
          handleApiError(res.error);
          continue;
        }

        for (let i = 0; i < res.data.items.length; i++) {
          const position = page * QUEUE_PAGE_SIZE + i;
          const mediaItem = trackToMediaItem(res.data.items[i]);
          this.loadedItems.set(position, mediaItem);
        }
      }
    }

    const result = new Map<number, MediaItem>();
    for (const index of indices) {
      const item = this.loadedItems.get(index);
      if (item) result.set(index, item);
    }
    return result;
  }

  getCurrentMediaItem(): MediaItem | null {
    return this.loadedItems.get(this.index) ?? null;
  }

  async #loadRange(start: number, end: number): Promise<MediaItem[]> {
    const indices: number[] = [];
    for (let i = start; i < end; i++) indices.push(i);

    const loaded = await this.loadItems(indices);
    const result: MediaItem[] = [];
    for (let i = start; i < end; i++) {
      const item = loaded.get(i);
      if (item) result.push(item);
    }
    return result;
  }

  async getPreviousItems(
    page: number,
    perPage: number,
  ): Promise<{ items: MediaItem[]; startIndex: number }> {
    const end = this.index - page * perPage;
    const start = Math.max(0, end - perPage);
    const items = await this.#loadRange(start, end);
    return { items, startIndex: start };
  }

  async getNextItems(page: number, perPage: number): Promise<MediaItem[]> {
    const start = this.index + 1 + page * perPage;
    const end = Math.min(this.totalItems, start + perPage);
    return this.#loadRange(start, end);
  }

  isEndOfQueue() {
    return this.index >= this.totalItems - 1;
  }

  isQueueEmpty() {
    return this.totalItems === 0;
  }
}

function generateRandomId(): string {
  const bytes = new Uint8Array(16);
  crypto.getRandomValues(bytes);
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

function getDeviceId(): string {
  let deviceId = localStorage.getItem("device-id");
  if (!deviceId) {
    deviceId = generateRandomId();
    localStorage.setItem("device-id", deviceId);
  }
  return deviceId;
}

function getVolume(): number {
  const volume = localStorage.getItem("player-volume");
  if (volume) {
    return parseFloat(volume);
  }

  return 1.0;
}

function getMuted(): boolean {
  const muted = localStorage.getItem("player-muted");
  if (muted) {
    return muted === "true";
  }

  return false;
}

export class MusicManager {
  #audio: HTMLAudioElement;

  private apiClient: ApiClient;
  queue: Queue;

  loading = $state(false);
  playing = $state(false);

  currentTime = $state(0);
  duration = $state(0);
  buffered = $state(0);

  #volume = $state(getVolume());
  #muted = $state(getMuted());

  get volume() {
    return this.#volume;
  }

  set volume(v: number) {
    this.#volume = v;
    this.#audio.volume = v;
    localStorage.setItem("player-volume", v.toString());
  }

  get muted() {
    return this.#muted;
  }

  set muted(m: boolean) {
    this.#muted = m;
    this.#audio.volume = m ? 0 : this.#volume;
    localStorage.setItem("player-muted", m.toString());
  }

  showPlayer = $state(false);

  currentItem = $state<MediaItem | null>(null);

  #trackEventSent = $state(false);

  #deviceId = getDeviceId();

  constructor(apiClient: ApiClient) {
    this.#audio = new Audio();

    this.apiClient = apiClient;
    this.queue = new Queue(apiClient, this.#deviceId);

    this.#setupAudio();
    this.#setupMediaSession();

    this.#audio.volume = this.#muted ? 0 : this.#volume;
  }

  async initQueue() {
    await this.queue.loadEntries();
    await this.#queueUpdate();
  }

  reset() {
    this.#audio.pause();
    this.#audio.removeAttribute("src");
    this.#audio.load();

    this.queue.setQueue([], 0, 0);

    this.playing = false;
    this.loading = false;
    this.showPlayer = false;
    this.currentItem = null;
    this.currentTime = 0;
    this.duration = 0;
    this.buffered = 0;
    this.#trackEventSent = false;

    this.#updateMediaSession();
  }

  async #refreshQueue() {
    await this.queue.loadEntries();
    await this.#queueUpdate();
  }

  #setupAudio() {
    this.#audio.addEventListener("canplay", () => {
      this.loading = false;
    });

    this.#audio.addEventListener("loadstart", () => {
      this.loading = true;
    });

    this.#audio.addEventListener("loadedmetadata", () => {
      this.currentTime = this.#audio.currentTime;
      this.duration = this.#audio.duration;
    });

    this.#audio.addEventListener("progress", () => {
      const tr = this.#audio.buffered;
      if (tr.length > 0) {
        const end = tr.end(tr.length - 1);
        this.buffered = this.duration > 0 ? end / this.duration : 0;
      }
    });

    this.#audio.addEventListener("timeupdate", () => {
      this.currentTime = this.#audio.currentTime;
    });

    this.#audio.addEventListener("playing", () => {
      this.playing = true;
      this.#updateMediaSession();
    });

    this.#audio.addEventListener("pause", () => {
      this.playing = false;
      this.#updateMediaSession();
    });

    this.#audio.addEventListener("ended", async () => {
      await this.nextTrack();
    });
  }

  #setupMediaSession() {
    if (!("mediaSession" in navigator)) return;

    navigator.mediaSession.setActionHandler("play", () => this.play());
    navigator.mediaSession.setActionHandler("pause", () => this.pause());
    navigator.mediaSession.setActionHandler("nexttrack", () =>
      this.nextTrack(),
    );
    navigator.mediaSession.setActionHandler("previoustrack", () =>
      this.previousTrack(),
    );
    navigator.mediaSession.setActionHandler("seekto", (details) => {
      if (details.seekTime != null) {
        this.setPosition(details.seekTime);
      }
    });
    navigator.mediaSession.setActionHandler("seekforward", () => {
      this.setPosition(Math.min(this.#audio.currentTime + 10, this.duration));
    });
    navigator.mediaSession.setActionHandler("seekbackward", () => {
      this.setPosition(Math.max(this.#audio.currentTime - 10, 0));
    });
  }

  #updateMediaSession() {
    if (!("mediaSession" in navigator)) return;

    const item = this.currentItem;
    if (!item) {
      navigator.mediaSession.playbackState = "none";
      return;
    }

    navigator.mediaSession.metadata = new MediaMetadata({
      title: item.name,
      artist: item.artists.map((a) => a.name).join(", "),
      album: item.album.name,
      artwork: [{ src: item.coverArt }],
    });

    navigator.mediaSession.playbackState = this.playing ? "playing" : "paused";
  }

  async #checkSendTrackEvent() {
    if (this.#trackEventSent) return;
    if (!this.currentItem) return;

    const pct = this.duration > 0 ? this.currentTime / this.duration : 0;
    const percentPlayed = Math.round(pct * 100);

    const res = await this.apiClient.pushTrackHistory({
      trackId: this.currentItem.trackId,
      playbackType: "normal",
      percentPlayed,
    });

    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    this.#trackEventSent = true;
  }

  setPosition(position: number) {
    this.#audio.currentTime = position;
  }

  async removeQueueItem(index: number) {
    const entry = this.queue.entries[index];
    if (!entry) return;

    const res = await this.apiClient.removeQueueItem(
      this.#deviceId,
      entry.queueItemId,
    );
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    await this.#refreshQueue();
  }

  async clearQueue(update = true) {
    await this.#checkSendTrackEvent();

    const res = await this.apiClient.clearQueue(this.#deviceId);
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    if (update) {
      await this.#refreshQueue();
    }
  }

  async setQueueIndex(index: number) {
    await this.#checkSendTrackEvent();

    const res = await this.apiClient.setQueuePosition(this.#deviceId, {
      index,
    });
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    await this.queue.loadEntries();
    await this.#queueUpdate();
  }

  async #resolveCurrentMediaItem(): Promise<MediaItem | null> {
    let mediaItem = this.queue.getCurrentMediaItem();
    if (!mediaItem && !this.queue.isQueueEmpty()) {
      mediaItem = await this.queue.loadItem(this.queue.index);
    }
    return mediaItem;
  }

  #applyMediaItem(mediaItem: MediaItem | null) {
    if (mediaItem) {
      if (mediaItem.trackId !== this.currentItem?.trackId) {
        this.#trackEventSent = false;
      }

      const src = this.apiClient.url.streamTrack(mediaItem.trackId).toString();
      this.#audio.src = src;

      this.currentItem = mediaItem;
    } else {
      this.#audio.removeAttribute("src");
      this.#audio.load();
      this.currentItem = null;
    }
  }

  async #queueUpdate() {
    this.showPlayer = !this.queue.isQueueEmpty();
    const mediaItem = await this.#resolveCurrentMediaItem();
    this.#applyMediaItem(mediaItem);
    this.#updateMediaSession();
  }

  async nextTrack() {
    await this.setQueueIndex(this.queue.index + 1);
    this.play();
  }

  async previousTrack() {
    await this.setQueueIndex(this.queue.index - 1);
    this.play();
  }

  play() {
    this.#audio.play();
  }

  pause() {
    this.#audio.pause();
  }

  async queueRequest(
    request:
      | { type: "addArtist"; artistId: string }
      | { type: "addAlbum"; albumId: string }
      | { type: "addPlaylist"; playlistId: string; filterId?: string }
      | { type: "addFavorites"; userId: string; filterId?: string }
      | { type: "addFilter"; filterId: string },
    options: {
      shuffle?: boolean;
      append?: "back";
      queueIndexToTrackId?: string;
    } = {},
  ) {
    const position = options.append === "back" ? "end" : "replace";
    const body = {
      position,
      shuffle: options.shuffle ?? false,
      currentIndex: undefined as number | undefined,
      queueIndexToTrackId: options.queueIndexToTrackId ?? undefined,
    };

    let res: Awaited<ReturnType<typeof this.apiClient.addAlbumToQueue>>;

    switch (request.type) {
      case "addAlbum":
        res = await this.apiClient.addAlbumToQueue(
          this.#deviceId,
          request.albumId,
          body,
        );
        break;
      case "addArtist":
        res = await this.apiClient.addArtistToQueue(
          this.#deviceId,
          request.artistId,
          body,
        );
        break;
      case "addPlaylist":
        res = await this.apiClient.addPlaylistToQueue(
          this.#deviceId,
          request.playlistId,
          body,
        );
        break;
      case "addFavorites":
        res = await this.apiClient.addFavoritesToQueue(
          this.#deviceId,
          request.userId,
          { ...body, filterId: request.filterId },
        );
        break;
      case "addFilter":
        res = await this.apiClient.addTracksToQueue(this.#deviceId, {
          ...body,
          trackIds: [],
          filterId: request.filterId,
        });
        break;
    }

    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    await this.#refreshQueue();

    if (!options.append) {
      this.play();
    }
  }

  async addAlbumTracks(params: {
    albumId: string;
    clear?: boolean;
    trackId?: string;
    shuffle?: boolean;
  }) {
    await this.queueRequest(
      { type: "addAlbum", albumId: params.albumId },
      {
        queueIndexToTrackId: params.trackId,
        shuffle: params.shuffle,
      },
    );
  }

  async addTracks(params: {
    trackIds?: string[];
    trackId?: string;
    clear?: boolean;
  }) {
    if (!params.trackIds || params.trackIds.length === 0) {
      return;
    }

    const res = await this.apiClient.addTracksToQueue(this.#deviceId, {
      trackIds: params.trackIds,
      position: "replace",
      shuffle: false,
      currentIndex: undefined,
      queueIndexToTrackId: params.trackId ?? undefined,
    });
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    await this.#refreshQueue();
    this.play();
  }
}

const MUSIC_MANAGER_KEY = Symbol("MUSIC_MANAGER");

export function setMusicManager(apiClient: ApiClient) {
  return setContext(MUSIC_MANAGER_KEY, new MusicManager(apiClient));
}

export function getMusicManager() {
  return getContext<ReturnType<typeof setMusicManager>>(MUSIC_MANAGER_KEY);
}
