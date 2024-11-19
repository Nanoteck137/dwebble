import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

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

  deleteTrack: async ({ locals, request }) => {
    const formData = await request.formData();

    const trackId = formData.get("trackId");
    if (!trackId) {
      throw error(400, { message: "trackId is not set" });
    }

    const res = await locals.apiClient.removeTrack(trackId.toString());
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },

  editTrack: async ({ locals, request }) => {
    const formData = await request.formData();

    const trackId = formData.get("trackId");
    if (!trackId) {
      throw error(400, { message: "trackId is not set" });
    }

    const trackName = formData.get("trackName");
    if (!trackName) {
      throw error(400, { message: "trackName is not set" });
    }

    const trackTags = formData.get("trackTags");
    if (trackTags === null) {
      throw error(400, { message: "trackTags is not set" });
    }

    const trackNumber = formData.get("trackNumber");
    if (trackNumber === null) {
      throw error(400, { message: "trackNumber is not set" });
    }

    console.log("Track Number", trackNumber);

    let t = trackTags.toString().split(",");
    t = t.map((t) => t.trim()).filter((t) => t !== "");

    const res = await locals.apiClient.editTrack(trackId.toString(), {
      name: trackName.toString(),
      tags: t,
      number:
        trackNumber.toString() !== ""
          ? parseInt(trackNumber.toString())
          : null,
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },
};
