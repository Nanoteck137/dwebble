import { z } from "zod";
import { Track } from "./track";

export const Album = z.object({
  id: z.string().cuid2(),
  name: z.string(),
  picture: z.string(),
});

export type Album = z.infer<typeof Album>;

export const FullAlbum = Album.extend({ tracks: z.array(Track) });
export type FullAlbum = z.infer<typeof FullAlbum>;
