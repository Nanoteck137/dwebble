import { Artist } from "$lib/models/artist";
import { z } from "zod";
import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";

export const load: PageServerLoad = async ({ fetch }) => {
  console.log(env.DATABASE_URL);

  const res = await fetch("http://127.0.0.1:3000/api/artists");
  const data = z.array(Artist).parse(await res.json());

  console.log(data);

  return {
    artists: data,
  };
};
