import type { Name, Track } from "$lib/api/types";
import type { MusicTrack } from "$lib/music-manager";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function formatName(name: Name) {
  let s = name.default;
  if (name.other) {
    s += ` - ${name.other}`;
  }

  return s;
}

type BareBoneError = {
  code: number;
  message: string;
  type: string;
};

export function formatError(err: BareBoneError) {
  return `${err.type} (${err.code}): ${err.message}`;
}

export function formatTime(s: number) {
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  return `${min}:${sec.toString().padStart(2, "0")}`;
}

export function trackToMusicTrack(track: Track): MusicTrack {
  return {
    name: track.name.default,
    artistName: track.artistName.default,
    source: track.mobileMediaUrl,
    coverArt: track.coverArt.small,
  };
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function shuffle<T>(arr: T[]): T[] {
  for (let i = arr.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    const temp = arr[i];
    arr[i] = arr[j];
    arr[j] = temp;
  }
  return arr;
}
