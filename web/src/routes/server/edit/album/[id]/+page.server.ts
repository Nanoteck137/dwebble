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

  editAlbumName: async ({ locals, request }) => {
    const formData = await request.formData();

    const albumId = formData.get("albumId");
    if (!albumId) {
      throw error(400, { message: "albumId is not set" });
    }

    const albumName = formData.get("albumName");
    if (!albumName) {
      throw error(400, { message: "albumName is not set" });
    }

    const res = await locals.apiClient.editAlbum(albumId.toString(), {
      name: albumName.toString(),
      artistId: null,
      artistName: null,
      year: null,
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }
  },

  editAlbumArtist: async ({ locals, request }) => {
    const formData = await request.formData();

    const albumId = formData.get("albumId");
    if (!albumId) {
      throw error(400, { message: "albumId is not set" });
    }

    const artistName = formData.get("artistName");
    if (!artistName) {
      throw error(400, { message: "artistName is not set" });
    }

    const res = await locals.apiClient.editAlbum(albumId.toString(), {
      name: null,
      artistId: null,
      artistName: artistName.toString(),
      year: null,
    });
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

  editTracks: async ({ locals, request }) => {
    const formData = await request.formData();
    console.log(formData);

    type TrackEdit = {
      id: string;
      num: number;
      name: string;
      artist: string;
      year: number;
      tags: string[];
    };

    const trackIds = formData.getAll("trackId");
    const trackNumbers = formData.getAll("trackNumber");
    const trackName = formData.getAll("trackName");
    const trackYear = formData.getAll("trackYear");
    const trackTags = formData.getAll("trackTags");
    const trackArtist = formData.getAll("trackArtist");

    const trackEdits: TrackEdit[] = [];
    for (let i = 0; i < trackIds.length; i++) {
      let tags = trackTags[i].toString().split(",");
      tags = tags.map((t) => t.trim()).filter((t) => t !== "");
      trackEdits.push({
        id: trackIds[i].toString(),
        num: parseInt(trackNumbers[i].toString()),
        name: trackName[i].toString(),
        artist: trackArtist[i].toString(),
        year: parseInt(trackYear[i].toString()),
        tags: tags,
      });
    }

    for (const track of trackEdits) {
      const res = await locals.apiClient.editTrack(track.id, {
        name: track.name,
        number: track.num,
        year: track.year,
        artistName: track.artist,
        tags: track.tags,
      });
      if (!res.success) {
        error(res.error.code, { message: res.error.message });
      }
    }
  },
};
