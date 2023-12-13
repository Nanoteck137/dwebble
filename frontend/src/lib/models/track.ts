import { z } from "zod";

export const Track = z.object({
  id: z.string().cuid2(),
  num: z.number(),
  name: z.string(),

  artist_id: z.string().cuid2(),
  album_id: z.string().cuid2(),

  album_name: z.string(),

  file_quality: z.string(),
  file_mobile: z.string(),
});

export type Track = z.infer<typeof Track>;
