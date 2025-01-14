import type { ApiClient } from "$lib/api/client";
import type { MusicTrack } from "$lib/api/types";
import { type Emitter, createNanoEvents } from "nanoevents";
import { getContext, setContext } from "svelte";

export abstract class Queue {
  items: MusicTrack[] = [];
  index: number = 0;

  async initialize() {}

  setQueueItems(tracks: MusicTrack[], index?: number) {
    this.items = tracks;
    this.index = index ?? 0;
  }

  getCurrentTrack() {
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

  // eslint-disable-next-line no-unused-vars
  abstract addFromAlbum(albumId: string): Promise<void>;
}

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

  async setQueueIndex(index: number) {
    super.setQueueIndex(index);

    await this.apiClient.updateQueue(this.queueId, {
      itemIndex: this.index,
    });
  }
}

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

  async addFromAlbum(albumId: string) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getAlbumTracksForPlay(albumId);
    if (res.success) {
      res.data.tracks.forEach((track) => {
        this.items.push(track);
      });
    }
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

  async addFromAlbum(albumId: string) {
    await this.queue.addFromAlbum(albumId);
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
