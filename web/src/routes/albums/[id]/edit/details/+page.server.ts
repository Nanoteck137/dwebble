import { error, redirect } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  return {
    apiAddress: locals.apiAddress,
    token: locals.token,
  };
};

export const actions: Actions = {
  submitEdit: async ({ locals, request, params }) => {
    const formData = await request.formData();

    const name = formData.get("name");
    if (!name) {
      throw error(400, { message: "name is not set" });
    }

    const otherName = formData.get("otherName");
    if (otherName === null) {
      throw error(400, { message: "otherName is not set" });
    }

    const year = formData.get("year");
    if (!year) {
      throw error(400, { message: "year is not set" });
    }

    const artistId = formData.get("artistId");
    if (artistId === null) {
      throw error(400, { message: "artistId is not set" });
    }

    const res = await locals.apiClient.editAlbum(params.id, {
      name: name.toString(),
      otherName: otherName.toString(),
      artistId: artistId.toString(),
      artistName: null,
      year: year ? parseInt(year.toString()) : null,
    });
    if (!res.success) {
      throw error(res.error.code, { message: res.error.message });
    }

    throw redirect(302, `/albums/${params.id}/edit`);
  },
  createArtist: async ({ locals, request, params }) => {
    const formData = await request.formData();

    const name = formData.get("name");
    if (!name) {
      throw error(400, { message: "name is not set" });
    }

    console.log("Create artist");

    // const res = await locals.apiClient.create(params.id, {
    //   name: name.toString(),
    //   otherName: otherName.toString(),
    //   artistId: artistId.toString(),
    //   artistName: null,
    //   year: year ? parseInt(year.toString()) : null,
    // });
    // if (!res.success) {
    //   throw error(res.error.code, { message: res.error.message });
    // }

    return {
      artistId: "CREATED_ARTIST",
      artistName: "Created Artist",
    };
  },
};
