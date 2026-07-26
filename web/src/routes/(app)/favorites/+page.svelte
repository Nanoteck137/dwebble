<script lang="ts">
  import { Button, Separator } from "@nanoteck137/nano-ui";
  import { Play, Shuffle } from "lucide-svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import Spacer from "$lib/components/Spacer.svelte";

  let { data } = $props();

  const musicManager = getMusicManager();

  async function playAll() {
    await musicManager.queueRequest(
      { type: "addFavorites", userId: data.user.id },
    );
  }
</script>

<div class="flex flex-col gap-4">
  <div
    class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between"
  >
    <div class="flex items-baseline gap-2">
      <h1 class="text-xl font-bold">Favorites</h1>
      {#if data.page}
        <span class="text-sm text-muted-foreground"
          >{data.page.totalItems}</span
        >
      {/if}
    </div>

    <div class="flex items-center gap-2">
      <Button size="sm" onclick={() => playAll()}>
        <Play size={14} />
        Play All
      </Button>
      <Button
        variant="outline"
        size="sm"
        onclick={async () => {
          await musicManager.queueRequest(
            { type: "addFavorites", userId: data.user.id },
            { shuffle: true },
          );
        }}
      >
        <Shuffle size={14} />
        Shuffle
      </Button>
    </div>
  </div>
</div>

<Spacer size="lg" />

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  onPlay={() => {}}
/>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />
