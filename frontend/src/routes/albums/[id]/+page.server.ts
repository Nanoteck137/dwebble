import { FullAlbum } from "$lib/models/album";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, fetch }) => {
  const id = params.id;

  const res = await fetch(`http://127.0.0.1:3000/api/albums/${id}`);
  const d = await res.json();
  const album = FullAlbum.parse(d);

  return {
    album,
  };
};
