import type { ApiClient } from "$lib/api/client";
import type { MusicTrack } from "$lib/api/types";
import { type Emitter, createNanoEvents } from "nanoevents";
import { getContext, setContext } from "svelte";

export type Queue = {
  index: number;
  items: MusicTrack[];
};

export class MusicManager {
  apiClient: ApiClient;
  queue: Queue = { index: 0, items: [] };
  emitter: Emitter;

  constructor(apiClient: ApiClient) {
    this.apiClient = apiClient;
    this.emitter = createNanoEvents();
  }

  getCurrentTrack() {
    if (this.queue.items.length === 0) return null;

    return this.queue.items[this.queue.index];
  }

  setQueueItems(tracks: MusicTrack[], index?: number) {
    this.queue.items = tracks;
    this.queue.index = index ?? 0;

    this.emitter.emit("onQueueUpdated");
    this.emitter.emit("onTrackChanged");
  }

  addTrackToQueue(track: MusicTrack, requestPlay: boolean = true) {
    const play = this.queue.items.length === 0;

    this.queue.items.push(track);
    this.emitter.emit("onQueueUpdated");

    if (play && requestPlay) {
      this.emitter.emit("onTrackChanged");
      this.requestPlay();
    }
  }

  isEndOfQueue() {
    return this.queue.index >= this.queue.items.length - 1;
  }

  isQueueEmpty() {
    return this.queue.items.length === 0;
  }

  clearQueue() {
    this.queue.index = 0;
    this.queue.items = [];

    this.emitter.emit("onQueueUpdated");
    this.emitter.emit("onTrackChanged");
  }

  setQueueIndex(index: number) {
    if (index >= this.queue.items.length) {
      return;
    }

    if (index < 0) {
      return;
    }

    this.queue.index = index;

    this.emitter.emit("onQueueUpdated");
    this.emitter.emit("onTrackChanged");
  }

  nextTrack() {
    this.setQueueIndex(this.queue.index + 1);
  }

  prevTrack() {
    this.setQueueIndex(this.queue.index - 1);
  }

  requestPlay() {
    this.emitter.emit("requestPlay");
  }

  requestPause() {
    this.emitter.emit("requestPause");
  }

  requestPlayPause() {
    this.emitter.emit("requestPlayPause");
  }

  markAsListened() {
    console.log("Update server");
  }
}

const MUSIC_MANAGER_KEY = Symbol("MUSIC_MANAGER");

export function setMusicManager(apiClient: ApiClient) {
  return setContext(MUSIC_MANAGER_KEY, new MusicManager(apiClient));
}

export function getMusicManager() {
  return getContext<ReturnType<typeof setMusicManager>>(MUSIC_MANAGER_KEY);
}
