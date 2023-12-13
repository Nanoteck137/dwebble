import { FullArtist } from "$lib/models/artist";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, fetch }) => {
  const id = params.id;

  const res = await fetch(`http://127.0.0.1:3000/api/artists/${id}`);
  const artist = FullArtist.parse(await res.json());

  return {
    artist,
  };
};
