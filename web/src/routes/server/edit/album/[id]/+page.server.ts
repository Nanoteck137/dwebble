import { error, redirect } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const album = await locals.apiClient.getAlbumById(params.id);
  if (!album.success) {
    throw error(album.error.code, { message: album.error.message });
  }

  const tracks = await locals.apiClient.getAlbumTracks(params.id);
  if (!tracks.success) {
    throw error(tracks.error.code, { message: tracks.error.message });
  }

  return {
    album: album.data,
    tracks: tracks.data.tracks,
  };
};

export const actions: Actions = {
  deleteAlbum: async ({ locals, request }) => {
    const formData = await request.formData();
    const albumId = formData.get("albumId");
    if (!albumId) {
      throw error(400, { message: "albumId is not set" });
    }

    const res = await locals.apiClient.deleteAlbum(albumId.toString());
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(301, `/server/edit/album`);
  },

  importTracks: async ({ locals, request }) => {
    const formData = await request.formData();

    const albumId = formData.get("albumId");
    if (!albumId) {
      throw error(400, { message: "albumId is not set" });
    }

    const body = new FormData();
    const files = formData.getAll("files");
    files.forEach((f) => {
      const file = f as File;
      body.append("files", file);
    });

    const res = await locals.apiClient.importTrackToAlbum(
      albumId.toString(),
      body,
    );
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },
};
