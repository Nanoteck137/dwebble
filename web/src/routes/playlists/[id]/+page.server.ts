import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, locals }) => {
  const playlist = await locals.apiClient.getPlaylistById(params.id);
  if (!playlist.success) {
    throw error(playlist.error.code, { message: playlist.error.message });
  }

  const items = await locals.apiClient.getPlaylistItems(params.id);
  if (!items.success) {
    throw error(items.error.code, { message: items.error.message });
  }

  return {
    playlist: playlist.data,
    items: items.data.items,
  };
};

// export const actions: Actions = {
//   remove: async ({ locals, request }) => {
//     const formData = await request.formData();

//     const playlistId = formData.get("playlistId");
//     if (!playlistId) {
//       throw error(500, "playlistId not set");
//     }

//     const tracks = formData.getAll("tracks[]");
//     console.log(tracks);

//     if (tracks.length <= 0) return;

//     const trackIds = tracks.map((e) => e.toString());
//     const res = await locals.apiClient.deletePlaylistItems(
//       playlistId.toString(),
//       { trackIds },
//     );
//     if (!res.success) {
//       throw error(res.error.code, { message: res.error.message });
//     }
//   },
// };
