<script lang="ts">
  import { invalidate } from "$app/navigation";
  import { createApiResponse } from "$lib/api/base-client";
  import { GetSync } from "$lib/api/types";
  import { onMount } from "svelte";
  import { z } from "zod";
  import type { SubmitFunction } from "./$types";
  import { enhance } from "$app/forms";
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";

  let syncing = $state(false);

  async function updateSyncStatus() {
    const res = await fetch("/status/sync");
    const Schema = createApiResponse(GetSync, z.undefined());

    const data = await res.json();
    const parsedData = await Schema.parseAsync(data);
    if (!parsedData.success) {
      console.error("Failed to get sync status", parsedData.error.message);
      return;
    }

    syncing = parsedData.data.isSyncing;
  }

  onMount(() => {
    updateSyncStatus();

    const interval = setInterval(async () => {
      await updateSyncStatus();
    }, 1000);

    return () => {
      clearInterval(interval);
    };
  });

  const submitRunSync: SubmitFunction = ({}) => {
    return async ({ result }) => {
      switch (result.type) {
        case "success":
          await updateSyncStatus();
      }
    };
  };
</script>

<p>Server Page (W.I.P)</p>

<form action="?/runSync" method="post" use:enhance={submitRunSync}>
  <button class="button-primary rounded px-4 py-2">Run Sync</button>
</form>

<p>Syncing: {syncing}</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

<div class="flex flex-col">
  <a href="/server/import">Import Album</a>
  <a href="/server/edit/album">Edit Albums</a>
</div>

<form action="?/runExport" method="post">
  <button>Run Export</button>
</form>
