import { browser } from "$app/environment";
import { derived, get, writable, type Writable } from "svelte/store";
import type { Track } from "./models/track";

type PlayQueue = {
  index: number;
  songs: Track[];
};

type PlayerState = "paused" | "loading" | "playing";

// TODO(patrik): Change name of isPlaying
export const isPlaying = writable<PlayerState>("paused");
export const playQueue = writable<PlayQueue>({ index: 0, songs: [] });
export const currentPlayingSong = derived<
  Writable<PlayQueue>,
  Track | undefined
>(
  playQueue,
  (playQueue: PlayQueue) => playQueue.songs[playQueue.index],
  undefined,
);

let audio: HTMLAudioElement | null = null;

function init() {
  audio = new Audio();
  audio.addEventListener("timeupdate", () => {
    const currentTime = audio!.currentTime;
    let duration = audio!.duration;
    duration = isNaN(duration) ? 0 : duration;

    time.set({ currentTime, duration });
  });
  // audio.addEventListener("volumechange", () => {});
  audio.addEventListener("pause", () => {
    isPlaying.set("paused");
  });
  audio.addEventListener("playing", () => {
    isPlaying.set("playing");
  });
  audio.addEventListener("loadstart", () => {
    isPlaying.set("loading");
  });
  audio.addEventListener("loadeddata", () => {
    isPlaying.set("playing");
  });
  audio.addEventListener("canplay", () => {});
  audio.addEventListener("ended", () => AudioHandler.nextSong());

  volume.subscribe((volume) => {
    setVolume(volume);
  });
}

function setVolume(volume: number) {
  if (!audio) return;

  audio.volume = volume;
  localStorage.setItem("volume", volume.toString());
}

export const time = writable<{ currentTime: number; duration: number }>({
  currentTime: 0.0,
  duration: 0.0,
});

export const volume = writable(
  browser ? parseFloat(localStorage.getItem("volume") ?? "1.0") : 0.0,
);

export const AudioHandler = {
  init,

  play() {
    audio?.play();
  },

  pause() {
    audio?.pause();
  },

  playPause() {
    if (get(isPlaying) == "playing") {
      audio?.pause();
    } else {
      audio?.play();
    }
  },

  setQueue(songs: Track[], index?: number) {
    playQueue.set({ index: index || 0, songs });
    if (audio) {
      audio.src = `http://localhost:3000/tracks/${
        songs[index || 0].file_mobile
      }`;
    }
    // handlers.setSong(songs[index || 0]);
  },

  nextSong() {
    const queue = get(playQueue);

    if (queue.index >= queue.songs.length - 1) return;

    const newIndex = queue.index + 1;
    playQueue.update((old) => ({ ...old, index: newIndex }));
    // handlers.setSong(queue.songs[newIndex]);
  },

  prevSong() {
    const queue = get(playQueue);

    if (queue.index <= 0) return;

    const newIndex = queue.index - 1;
    playQueue.update((old) => ({ ...old, index: newIndex }));
    // handlers.setSong(queue.songs[newIndex]);
  },
};
