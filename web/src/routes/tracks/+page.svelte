<script lang="ts">
  import {
    DiscAlbum,
    EllipsisVertical,
    Filter,
    ListPlus,
    Play,
    Shuffle,
    Star,
    Users,
  } from "lucide-svelte";
  import { musicManager } from "$lib/music-manager";
  import { cn, formatError, shuffle, trackToMusicTrack } from "$lib/utils";
  import {
    DropdownMenu,
    Button,
    buttonVariants,
    Input,
    Pagination,
  } from "@nanoteck137/nano-ui";
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { enhance } from "$app/forms";
  import type { UIArtist } from "$lib/types.js";
  import { createApiClient, openAddToPlaylist } from "$lib";
  import toast from "svelte-5-french-toast";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";

  let { data } = $props();
  const apiClient = createApiClient(data);
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

<TrackList
  {apiClient}
  totalTracks={data.page.totalItems}
  tracks={data.tracks}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  isInQuickPlaylist={(trackId) => {
    if (!data.quickPlaylistIds) return false;
    return !!data.quickPlaylistIds.find((v) => v === trackId);
  }}
/>

<Pagination.Root
  page={data.page.page + 1}
  count={data.page.totalItems}
  perPage={data.page.perPage}
  siblingCount={0}
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
