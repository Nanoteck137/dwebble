<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/state";
  import { Button, Separator } from "@nanoteck137/nano-ui";
  import { Play, Shuffle, X } from "lucide-svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import FilterButton from "../tracks/FilterButton.svelte";

  let { data } = $props();

  const musicManager = getMusicManager();

  let filterId = $derived(page.url.searchParams.get("filterId"));

  function clearFilter() {
    const query = page.url.searchParams;
    query.delete("filterId");
    goto("?" + query.toString(), {
      invalidateAll: true,
      replaceState: true,
    });
  }

  async function playAll() {
    await musicManager.queueRequest(
      { type: "addFavorites", userId: data.user.id, filterId: filterId ?? undefined },
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
            { type: "addFavorites", userId: data.user.id, filterId: filterId ?? undefined },
            { shuffle: true },
          );
        }}
      >
        <Shuffle size={14} />
        Shuffle
      </Button>
    </div>
  </div>

  {#if data.filters && data.filters.length > 0}
    <div class="flex flex-wrap items-center gap-2">
      {#each data.filters as filter (filter.filterId)}
        <FilterButton {filter} />
      {/each}

      {#if filterId}
        <Button variant="ghost" size="sm" onclick={clearFilter}>
          <X size={14} />
          Clear
        </Button>
      {/if}
    </div>
  {/if}
</div>

<Spacer size="lg" />

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  onPlay={async (trackId) => {
    await musicManager.queueRequest(
      { type: "addFavorites", userId: data.user.id, filterId: filterId ?? undefined },
      { queueIndexToTrackId: trackId },
    );
  }}
/>

<Spacer size="lg" />

<Separator />

<Spacer size="lg" />

<Pagination page={data.page} />
