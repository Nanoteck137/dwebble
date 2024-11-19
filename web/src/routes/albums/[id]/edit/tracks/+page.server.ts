import { error, redirect } from "@sveltejs/kit";
import type { Actions } from "./$types";

export const actions: Actions = {
  default: async ({ locals, request, params }) => {
    const formData = await request.formData();

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

    throw redirect(302, `/server/edit/album/${params.id}`);
  },
};
