// See https://kit.svelte.dev/docs/types#app

import type { ApiClient } from "$lib/api/client";
import type { GetAuthMe } from "$lib/api/types";

// for information about these interfaces
declare global {
  namespace App {
    interface Error {
      type?: string;
    }
    interface Locals {
      apiClient: ApiClient;
      user?: GetAuthMe;
    }
    // interface PageData {}
    // interface PageState {}
    // interface Platform {}
  }
}

export {};
