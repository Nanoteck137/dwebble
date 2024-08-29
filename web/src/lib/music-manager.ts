import { type Emitter, createNanoEvents } from "nanoevents";

export type MusicTrack = {
  name: string;
  artistName: string;
  source: string;
  coverArt: string;
};

export type Queue = {
  index: number;
  items: MusicTrack[];
};

export class MusicManager {
  queue: Queue = { index: 0, items: [] };
  emitter: Emitter;

  constructor() {
    this.emitter = createNanoEvents();
  }

  getCurrentTrack() {
    if (this.queue.items.length === 0) return null;

    return this.queue.items[this.queue.index];
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

const musicManager = new MusicManager();
export { musicManager };
