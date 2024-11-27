import type { Album, Artist, Track } from "$lib/api/types";
import { z } from "zod";

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
