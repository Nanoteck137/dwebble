import type { Album, Artist, Track } from "$lib/api/types";
import { z } from "zod";

export interface CheckedValue<T> {
  checked: boolean;
  value: T;
}

// TODO(patrik): Better name
export type UIArtist = {
  id: string;
  name: string;
};

export type EditTrackData = {
  name: string;
  otherName: string;
  artist: UIArtist;
  num?: number;
  year?: number;
  tags: string;
  featuringArtists: UIArtist[];
};

// TODO(patrik): Replace with UIArtist
export type QueryArtist = {
  id: string;
  name: string;
};

export const SetupSchema = z.object({
  user: z
    .object({
      username: z.string(),
      password: z.string(),
    })
    .optional(),
});

export type Setup = z.infer<typeof SetupSchema>;

export type SuccessSearch = {
  success: true;
  artists: Artist[];
  albums: Album[];
  tracks: Track[];
};

export type ErrorSearch = {
  success: false;
  message: string;
};

export type Search = SuccessSearch | ErrorSearch;
