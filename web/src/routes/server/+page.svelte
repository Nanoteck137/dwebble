<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";
  import { onDestroy, onMount } from "svelte";
  import { z } from "zod";

  const { data } = $props();
  const apiClient = getApiClient();

  let refillSearch = $state(false);

  let syncState = $state<SyncStateTy>();

  const SyncError = z.object({
    type: z.string(),
    message: z.string(),
    fullMessage: z.string().optional(),
  });

  const MissingAlbum = z.object({
    id: z.string(),
    name: z.string(),
    artistName: z.string(),
  });

  const MissingTrack = z.object({
    id: z.string(),
    name: z.string(),
    albumName: z.string(),
    artistName: z.string(),
  });

  const SyncState = z.object({
    isSyncing: z.boolean(),
    isRetrivingPaths: z.boolean(),

    paths: z.array(
      z.object({
        name: z.string(),
        path: z.string(),
        isDir: z.boolean(),
        depth: z.number(),
      }),
    ),

    report: z.object({
      syncErrors: z.array(SyncError).nullable(),
      missingAlbums: z.array(MissingAlbum).nullable(),
      missingTracks: z.array(MissingTrack).nullable(),
    }),
  });
  type SyncStateTy = z.infer<typeof SyncState>;

  const Event = z.discriminatedUnion("type", [
    z.object({
      type: z.literal("sync-state"),
      data: SyncState,
    }),
    z.object({
      type: z.literal("syncing"),
      data: z.object({
        syncing: z.boolean(),
      }),
    }),
    z.object({
      type: z.literal("report"),
      data: z.object({
        syncErrors: z.array(SyncError).nullable(),
        missingAlbums: z.array(MissingAlbum).nullable(),
        missingTracks: z.array(MissingTrack).nullable(),
      }),
    }),
  ]);

  onMount(() => {
    console.log("Mount");
    const eventSource = new EventSource(
      data.apiAddress + "/api/v1/system/library/sse",
    );

    eventSource.onmessage = (e) => {
      const event = Event.parse(JSON.parse(e.data));
      console.log(event);

      switch (event.type) {
        case "sync-state":
          syncState = event.data;
          break;
        // case "syncing":
        //   syncing = event.data.syncing;
        //   break;
        case "report":
          console.log("Report", event.data);
          // const mapped =
          //   event.data.reports?.map((t) => {
          //     if (t.fullMessage) return t.fullMessage;
          //     return t.message;
          //   }) ?? [];
          // test = mapped;
          break;
      }
    };

    return () => {
      console.log("Cleanup");
      eventSource.close();
    };
  });
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

<Button
  onclick={async () => {
    refillSearch = true;
    const res = await apiClient.refillSearch();
    if (!res.success) {
      handleApiError(res.error);
    }
    refillSearch = false;
  }}
  disabled={refillSearch}
>
  Refill Search
</Button>

<p>Library Syncing: {syncState?.isSyncing}</p>
<p>Library Retriving Paths: {syncState?.isRetrivingPaths}</p>

<Button
  onclick={async () => {
    const res = await apiClient.syncLibrary({});
    if (!res.success) {
      handleApiError(res.error);
      return;
    }
  }}
>
  Sync Library
</Button>

<Button
  onclick={async () => {
    const res = await apiClient.retrivePaths({});
    if (!res.success) {
      handleApiError(res.error);
      return;
    }
  }}
>
  Retrive Paths
</Button>

<Button
  onclick={async () => {
    const res = await apiClient.cleanupLibrary();
    if (!res.success) {
      handleApiError(res.error);
      return;
    }
  }}
>
  Cleanup Library
</Button>

{#if syncState?.report.syncErrors}
  {#each syncState?.report.syncErrors as err}
    <p class="whitespace-pre font-mono">{err.fullMessage}</p>
    <br />
  {/each}
{/if}

{#if syncState?.paths}
  <div class="flex flex-col items-start">
    {#each syncState?.paths as path}
      <button
        class="hover:underline"
        style={`padding-left: ${path.depth * 20}px`}
        onclick={async () => {
          const res = await apiClient.syncLibrary({ path: path.path });
          if (!res.success) {
            handleApiError(res.error);
            return;
          }
        }}
      >
        {path.path}
      </button>
    {/each}
  </div>
{/if}
