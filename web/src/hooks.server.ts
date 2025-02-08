import { env } from "$env/dynamic/private";
import { ApiClient } from "$lib/api/client";
import { redirect, type Handle } from "@sveltejs/kit";

const apiAddress = env.API_ADDRESS ? env.API_ADDRESS : "";

export const handle: Handle = async ({ event, resolve }) => {
  const url = new URL(event.request.url);

  let addr = apiAddress;
  if (addr === "") {
    addr = url.origin;
  }
  const client = new ApiClient(addr);
  event.locals.apiAddress = addr;

  const auth = event.cookies.get("auth");
  if (auth) {
    const obj = JSON.parse(auth);
    client.setToken(obj.token);
    event.locals.token = obj.token;

    const me = await client.getMe();
    if (!me.success) {
      event.cookies.delete("auth", { path: "/" });
      throw redirect(301, "/");
    }

    event.locals.user = me.data;
  }

  event.locals.apiClient = client;

  if (
    event.locals.user &&
    (url.pathname === "/login" || url.pathname === "/register")
  ) {
    throw redirect(301, "/");
  }

  if (
    !event.locals.user &&
    (url.pathname.startsWith("/taglists") ||
      url.pathname.startsWith("/playlists") ||
      url.pathname.startsWith("/account"))
  ) {
    throw redirect(301, "/");
  }

  const response = await resolve(event);
  return response;
};
