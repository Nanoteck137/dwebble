<script lang="ts">
  import { Filter } from "lucide-svelte";
  import { Button, Input, Pagination } from "@nanoteck137/nano-ui";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import Spacer from "$lib/components/Spacer.svelte";

  let { data } = $props();
  const musicManager = getMusicManager();
</script>

<form method="GET">
  <div class="flex flex-col gap-2">
    <Input
      type="text"
      name="filter"
      placeholder="Filter"
      value={data.filter ?? ""}
    />

    <Input
      type="text"
      name="sort"
      placeholder="Sort"
      value={data.sort ?? ""}
    />
  </div>

  {#if data.filterError}
    <p class="text-red-400">{data.filterError}</p>
  {/if}
  {#if data.sortError}
    <p class="text-red-400">{data.sortError}</p>
  {/if}
  <div class="h-2"></div>
  <Button type="submit">
    <Filter />
    Filter Tracks
  </Button>
</form>

<div class="h-2"></div>

<TrackListHeader
  name="Tracks"
  onPlay={async (shuffle) => {
    await musicManager.queueRequest(
      {
        type: "addFilter",
        filter: data.filter ?? "",
      },
      { shuffle },
    );
  }}
/>

<Spacer size="md" />

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  onPlay={async (trackId) => {
    await musicManager.queueRequest(
      {
        type: "addFilter",
        filter: data.filter ?? "",
      },
      { queueIndexToTrackId: trackId },
    );
  }}
/>

<Pagination.Root
  page={data.page.page + 1}
  count={data.page.totalItems}
  perPage={data.page.perPage}
  siblingCount={1}
  onPageChange={(p) => {
    const query = $page.url.searchParams;
    query.set("page", (p - 1).toString());

    goto(`?${query.toString()}`, { invalidateAll: true, keepFocus: true });
  }}
>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.PrevButton />
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link
              href="?page={page.value}"
              {page}
              isActive={currentPage === page.value}
            >
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
