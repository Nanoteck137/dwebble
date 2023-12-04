import { z } from "zod";
import type { PageServerLoad } from "./$types";

const Song = z.object({
  name: z.string(),
  url: z.string(),
});
type Song = z.infer<typeof Song>;

const Playlist = z.object({
  name: z.string(),
  songs: z.array(Song),
});
type Playlist = z.infer<typeof Playlist>;

export const load: PageServerLoad = async ({ fetch }) => {
  const res = await fetch("http://127.0.0.1:3000/api/playlists");
  const data = z.array(Playlist).parse(await res.json());

  return {
    playlists: data,
  };
};
