import type { Track } from "$lib/api/types";
import type { MusicTrack } from "$lib/music-manager";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function formatTime(s: number) {
  const min = Math.floor(s / 60);
  const sec = Math.floor(s % 60);

  return `${min}:${sec.toString().padStart(2, "0")}`;
}

export function trackToMusicTrack(track: Track): MusicTrack {
  return {
    name: track.name,
    artistName: track.artistName,
    source: track.mobileMediaUrl,
    coverArt: track.coverArtUrl,
  };
}

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
