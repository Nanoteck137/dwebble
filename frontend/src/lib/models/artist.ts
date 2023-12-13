import { z } from "zod";
import { Album } from "./album";

export const Artist = z.object({
  id: z.string().cuid2(),
  name: z.string(),
  picture: z.string(),
});

type Artist = z.infer<typeof Artist>;

export const FullArtist = Artist.extend({
  albums: z.array(Album),
});
export type FullArtist = z.infer<typeof FullArtist>;
