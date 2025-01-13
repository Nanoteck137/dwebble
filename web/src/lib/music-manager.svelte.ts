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

  queueId: string;
  queueItems: MusicTrack[] = [];
  queueIndex: number = 0;

  emitter: Emitter;

  constructor(apiClient: ApiClient, queueId: string) {
    this.apiClient = apiClient;
    this.queueId = queueId;
    this.emitter = createNanoEvents();

    this.refetchQueue();
  }

  async refetchQueue() {
    // TODO(patrik): Handle errors
    const queueItems = await this.apiClient.getQueueItems(this.queueId);
    console.log("Get queue items", queueItems);

    if (queueItems.success) {
      this.setQueueItems(
        queueItems.data.items.map((i) => i.track),
        queueItems.data.index,
      );
    }
  }

  async clearQueue() {
    // TODO(patrik): Handle errors
    const res = await this.apiClient.clearQueue(this.queueId);
    if (res.success) {
      await this.refetchQueue();
    }
  }

  async addFromAlbum(albumId: string) {
    // TODO(patrik): Handle errors
    const res = await this.apiClient.addToQueueFromAlbum(
      this.queueId,
      albumId,
    );
    if (res.success) {
      await this.refetchQueue();
    }
  }

  getCurrentTrack() {
    if (this.queueItems.length === 0) return null;
    return this.queueItems[this.queueIndex];
  }

  setQueueItems(tracks: MusicTrack[], index?: number) {
    this.queueItems = tracks;
    this.queueIndex = index ?? 0;
    this.emitter.emit("onQueueUpdated");
  }

  addTrackToQueue(track: MusicTrack, requestPlay: boolean = true) {
    throw "remove this function";
    // const play = this.queueItems.length === 0;

    // this.queue.items.push(track);
    // this.emitter.emit("onQueueUpdated");

    // if (play && requestPlay) {
    //   this.emitter.emit("onTrackChanged");
    //   this.requestPlay();
    // }
  }

  isEndOfQueue() {
    return this.queueIndex >= this.queueItems.length - 1;
  }

  isQueueEmpty() {
    return this.queueItems.length === 0;
  }

  // clearQueue() {
  //   this.queueIndex = 0;
  //   this.queueItems = [];

  //   this.emitter.emit("onQueueUpdated");
  // }

  async setQueueIndex(index: number) {
    if (index >= this.queueItems.length) {
      return;
    }

    if (index < 0) {
      return;
    }

    this.queueIndex = index;

    this.emitter.emit("onQueueUpdated");

    await this.apiClient.updateQueue(this.queueId, {
      itemIndex: this.queueIndex,
    });
  }

  async nextTrack() {
    await this.setQueueIndex(this.queueIndex + 1);
  }

  async previousTrack() {
    await this.setQueueIndex(this.queueIndex - 1);
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

export function setMusicManager(apiClient: ApiClient, queueId: string) {
  return setContext(MUSIC_MANAGER_KEY, new MusicManager(apiClient, queueId));
}

export function getMusicManager() {
  return getContext<ReturnType<typeof setMusicManager>>(MUSIC_MANAGER_KEY);
}
