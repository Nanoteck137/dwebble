/* eslint-disable no-unused-vars */
import type { ApiClient } from "$lib/api/client";
import type { MediaItem } from "$lib/api/types";
import { type Emitter, createNanoEvents } from "nanoevents";
import { getContext, setContext } from "svelte";

type AddToQueueSettings = {
  shuffle?: boolean;
  front?: boolean;
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
      this.index = 0;
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

type SavedQueue = {
  items: string[];
  index: number;
};

export class LocalQueue extends Queue {
  apiClient: ApiClient;

  constructor(apiClient: ApiClient) {
    super();

    this.apiClient = apiClient;
  }

  async initialize() {
    await this.loadQueue();
  }

  async clearQueue() {
    this.items = [];
    this.index = 0;

    this.saveQueue();
  }

  async setQueueIndex(index: number) {
    super.setQueueIndex(index);
    this.saveQueue();
  }

  saveQueue() {
    const q: SavedQueue = {
      items: this.items.map((i) => i.track.id),
      index: this.index,
    };

    localStorage.setItem("queue", JSON.stringify(q));
  }

  async loadQueue() {
    const s = localStorage.getItem("queue");
    if (s) {
      const q: SavedQueue = JSON.parse(s);
      await this.loadIds(q.items);
      this.index = q.index;
    }
  }

  async loadIds(trackIds: string[]) {
    const res = await this.apiClient.getMediaFromIds({
      trackIds,
      keepOrder: true,
    });
    if (res.success) {
      this.items = res.data.items;
    }
  }

  // TODO(patrik): Fix every add method to use settings
  async addFromPlaylist(playlistId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromPlaylist(playlistId, {
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      if (settings?.front) {
        this.items = [...res.data.items, ...this.items];
      } else {
        this.items = [...this.items, ...res.data.items];
      }
    }

    this.saveQueue();
  }

  async addFromTaglist(taglistId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromTaglist(taglistId, {
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      if (settings?.front) {
        this.items = [...res.data.items, ...this.items];
      } else {
        this.items = [...this.items, ...res.data.items];
      }
    }

    this.saveQueue();
  }

  async addFromFilter(filter: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromFilter({
      filter,
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      if (settings?.front) {
        this.items = [...res.data.items, ...this.items];
      } else {
        this.items = [...this.items, ...res.data.items];
      }
    }

    this.saveQueue();
  }

  async addFromArtist(artistId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromArtist(artistId, {
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      if (settings?.front) {
        this.items = [...res.data.items, ...this.items];
      } else {
        this.items = [...this.items, ...res.data.items];
      }
    }

    this.saveQueue();
  }

  async addFromAlbum(albumId: string, settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromAlbum(albumId, {
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      if (settings?.front) {
        this.items = [...res.data.items, ...this.items];
      } else {
        this.items = [...this.items, ...res.data.items];
      }
    }

    this.saveQueue();
  }

  async addFromIds(trackIds: string[], settings?: AddToQueueSettings) {
    // TODO(patrik): Handle error
    const res = await this.apiClient.getMediaFromIds({
      trackIds,
      shuffle: settings?.shuffle,
    });
    if (res.success) {
      if (settings?.front) {
        this.items = [...res.data.items, ...this.items];
      } else {
        this.items = [...this.items, ...res.data.items];
      }
    }

    this.saveQueue();
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

type QueueRequestOptions = {
  queueIndex?: number;
  queueIndexToTrackId?: string;
  shuffle?: boolean;
  skipPlay?: boolean;
  append?: "front" | "back";
};

type QueueRequestAddPlaylist = {
  type: "addPlaylist";
  playlistId: string;
};

type QueueRequestAddAlbum = {
  type: "addAlbum";
  albumId: string;
};

type QueueRequestAddArtist = {
  type: "addArtist";
  artistId: string;
};

type QueueRequestAddTaglist = {
  type: "addTaglist";
  taglistId: string;
};

type QueueRequestAddFilter = {
  type: "addFilter";
  filter: string;
};

type QueueRequest =
  | QueueRequestAddPlaylist
  | QueueRequestAddAlbum
  | QueueRequestAddArtist
  | QueueRequestAddTaglist
  | QueueRequestAddFilter;

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

  async queueRequest(request: QueueRequest, options: QueueRequestOptions) {
    // onPlay={async (shuffle) => {
    //   await musicManager.clearQueue();
    //   await musicManager.addFromPlaylist(data.playlist.id, { shuffle });
    //   musicManager.requestPlay();
    // }}
    // onTrackPlay={async (trackId) => {
    //   await musicManager.clearQueue();
    //   await musicManager.addFromPlaylist(data.playlist.id);
    //   await musicManager.setQueueIndex(
    //     musicManager.queue.items.findIndex((t) => t.track.id === trackId),
    //   );
    //   musicManager.requestPlay();
    // }}

    if (!options.append) {
      await this.clearQueue();
    }

    switch (request.type) {
      case "addPlaylist": {
        await this.addFromPlaylist(request.playlistId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addAlbum": {
        await this.addFromAlbum(request.albumId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addArtist": {
        await this.addFromArtist(request.artistId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addTaglist": {
        await this.addFromTaglist(request.taglistId, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
      case "addFilter": {
        await this.addFromFilter(request.filter, {
          shuffle: options.shuffle,
          front: options.append === "front",
        });
        break;
      }
    }

    if (options.queueIndex) {
      await this.setQueueIndex(options.queueIndex);
    }

    if (options.queueIndexToTrackId) {
      await this.setQueueIndex(
        this.queue.items.findIndex(
          (t) => t.track.id === options.queueIndexToTrackId,
        ),
      );
    }

    if (!options.skipPlay) {
      this.requestPlay();
    }
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
