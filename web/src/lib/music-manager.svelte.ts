/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-unused-vars */
import type { ApiClient } from "$lib/api/client";
import type { MediaItem } from "$lib/api/types";
import { type Emitter, createNanoEvents } from "nanoevents";
import { getContext, setContext } from "svelte";

type AddToQueueSettings = {
  shuffle?: boolean;
};

export abstract class Queue {
  items: MediaItem[] = [];
  index: number = 0;

  async initialize() {}

  setQueueItems(items: MediaItem[], index?: number) {
    this.items = items;
    this.index = index ?? 0;
  }

  getCurrentMediaItem() {
    if (this.items.length === 0) return null;
    return this.items[this.index];
  }

  isEndOfQueue() {
    return this.index >= this.items.length - 1;
  }

  isQueueEmpty() {
    return this.items.length === 0;
  }

  async setQueueIndex(index: number) {
    if (index >= this.items.length) {
      return;
    }

    if (index < 0) {
      return;
    }

    this.index = index;
  }

  abstract clearQueue(): Promise<void>;

  abstract addFromPlaylist(
    playlistId: string,
    settings?: AddToQueueSettings,
  ): Promise<void>;

  abstract addFromTaglist(
    taglistId: string,
    settings?: AddToQueueSettings,
  ): Promise<void>;

  abstract addFromAlbum(
    albumId: string,
    settings?: AddToQueueSettings,
  ): Promise<void>;
}

/*
export class BackendQueue extends Queue {
  apiClient: ApiClient;
  queueId: string;

  constructor(apiClient: ApiClient, queueId: string) {
    super();

    this.apiClient = apiClient;
    this.queueId = queueId;

    this.refetchQueue();
  }

  async initialize() {
    await this.refetchQueue();
  }

  async refetchQueue() {
    // TODO(patrik): Handle errors
    const queueItems = await this.apiClient.getQueueItems(this.queueId);
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

  async addFromAlbum(albumId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle errors
    const res = await this.apiClient.addToQueueFromAlbum(this.queueId, {
      albumId,
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      await this.refetchQueue();
    }
  }

  async addFromPlaylist(playlistId: string, settings?: AddToQueueSettings) {
    const res = await this.apiClient.addToQueueFromPlaylist(this.queueId, {
      playlistId,
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      await this.refetchQueue();
    }
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    const res = await this.apiClient.addToQueueFromTaglist(this.queueId, {
      taglistId,
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      await this.refetchQueue();
    }
  }

  async setQueueIndex(index: number) {
    super.setQueueIndex(index);

    await this.apiClient.updateQueue(this.queueId, {
      itemIndex: this.index,
    });
  }
}
*/

export class LocalQueue extends Queue {
  apiClient: ApiClient;

  constructor(apiClient: ApiClient) {
    super();

    this.apiClient = apiClient;
  }

  async clearQueue() {
    this.items = [];
    this.index = 0;
  }

  async addFromAlbum(albumId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromAlbum(albumId, {});
    if (res.success) {
      this.items = [...this.items, ...res.data.items];

      // res.data.items.forEach((track) => {
      //   this.items.push(track);
      // });
    }
  }

  async addFromPlaylist(playlistId: string, settings?: AddToQueueSettings) {
    throw new Error("Method not implemented.");
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    throw new Error("Method not implemented.");
  }
}

export class DummyQueue extends Queue {
  constructor() {
    super();
    console.log("Dummy Queue: Constructor");
  }

  async initialize() {
    console.log("Dummy Queue: initialize");
  }

  async clearQueue() {
    console.log("Dummy Queue: clearQueue");
  }

  async addFromAlbum() {
    console.log("Dummy Queue: addFromAlbum");
  }

  async addFromPlaylist() {
    console.log("Dummy Queue: addFromPlaylist");
  }

  async addFromTaglist() {
    console.log("Dummy Queue: addFromTaglist");
  }

  async setQueueIndex() {
    console.log("Dummy Queue: setQueueIndex");
  }
}

export class MusicManager {
  apiClient: ApiClient;
  queue: Queue;

  emitter: Emitter;

  constructor(apiClient: ApiClient, queue: Queue) {
    this.apiClient = apiClient;
    this.queue = queue;
    this.emitter = createNanoEvents();

    this.queue.initialize().then(() => {
      this.emitter.emit("onQueueUpdated");
    });
  }

  async clearQueue() {
    await this.queue.clearQueue();
    this.emitter.emit("onQueueUpdated");
  }

  async addFromAlbum(albumId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromAlbum(albumId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromPlaylist(playlistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromPlaylist(playlistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromTaglist(taglistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async setQueueIndex(index: number) {
    await this.queue.setQueueIndex(index);
    this.emitter.emit("onQueueUpdated");
  }

  emitQueueUpdate() {
    this.emitter.emit("onQueueUpdated");
  }

  async nextTrack() {
    await this.setQueueIndex(this.queue.index + 1);
    this.requestPlay();
  }

  async previousTrack() {
    await this.setQueueIndex(this.queue.index - 1);
    this.requestPlay();
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

  setQueue(queue: Queue) {
    this.queue = queue;
    this.emitter.emit("onQueueUpdated");

    this.queue.initialize().then(() => {
      this.emitter.emit("onQueueUpdated");
    });
  }
}

const MUSIC_MANAGER_KEY = Symbol("MUSIC_MANAGER");

export function setMusicManager(apiClient: ApiClient, queue: Queue) {
  return setContext(MUSIC_MANAGER_KEY, new MusicManager(apiClient, queue));
}

export function getMusicManager() {
  return getContext<ReturnType<typeof setMusicManager>>(MUSIC_MANAGER_KEY);
}
