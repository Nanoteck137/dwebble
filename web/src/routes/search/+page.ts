import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ data }) => {
  console.log("LOAD");
  return data;
};
