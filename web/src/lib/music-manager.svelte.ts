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

  abstract addFromFilter(
    filter: string,
    settings?: AddToQueueSettings,
  ): Promise<void>;

  abstract addFromArtist(
    artistId: string,
    settings?: AddToQueueSettings,
  ): Promise<void>;

  abstract addFromAlbum(
    albumId: string,
    settings?: AddToQueueSettings,
  ): Promise<void>;

  abstract addFromIds(
    trackIds: string[],
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

  async addFromPlaylist(playlistId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromPlaylist(playlistId, {});
    if (res.success) {
      this.items = [...this.items, ...res.data.items];
    }
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromTaglist(taglistId, {});
    if (res.success) {
      this.items = [...this.items, ...res.data.items];
    }
  }

  async addFromFilter(filter: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromFilter({ filter });
    if (res.success) {
      this.items = [...this.items, ...res.data.items];
    }
  }

  async addFromArtist(artistId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromArtist(artistId, {});
    if (res.success) {
      this.items = [...this.items, ...res.data.items];
    }
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

  async addFromIds(trackIds: string[], settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromIds({ trackIds });
    if (res.success) {
      this.items = [...this.items, ...res.data.items];
    }
  }
}

export class DummyQueue extends Queue {
  constructor() {
    super();
  }

  async initialize() {}

  async clearQueue() {}

  async addFromPlaylist() {}
  async addFromTaglist() {}
  async addFromFilter() {}
  async addFromArtist() {}
  async addFromAlbum() {}
  async addFromIds() {}

  async setQueueIndex() {}
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

  async addFromPlaylist(playlistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromPlaylist(playlistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromTaglist(taglistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromFilter(filter: string, settings?: AddToQueueSettings) {
    await this.queue.addFromFilter(filter, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromArtist(artistId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromArtist(artistId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromAlbum(albumId: string, settings?: AddToQueueSettings) {
    await this.queue.addFromAlbum(albumId, settings);
    this.emitter.emit("onQueueUpdated");
  }

  async addFromIds(trackIds: string[], settings?: AddToQueueSettings) {
    await this.queue.addFromIds(trackIds, settings);
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
