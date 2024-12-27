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

<div class="flex flex-col">
  <div class="flex gap-2">
    <Button
      size="sm"
      variant="ghost"
      onclick={() => {
        musicManager.clearQueue();
        for (const track of data.tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track));
        }
      }}
    >
      <Play />
      Play
    </Button>

    <form action="?/newPlaylist" method="post">
      <input name="filter" value={data.filter} type="hidden" />
      <input name="sort" value={data.sort} type="hidden" />
      <Button type="submit" size="sm" variant="ghost">
        <ListPlus />
        Create Playlist
      </Button>
    </form>

    <Button
      size="sm"
      variant="ghost"
      onclick={() => {
        let tracks = shuffle([...data.tracks]);

        musicManager.clearQueue();
        for (const track of tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track));
        }
      }}
    >
      <Shuffle />
      Shuffle Play
    </Button>
  </div>

  <div class="flex items-center justify-between">
    <p class="text-bold text-xl">Tracks</p>
    <p class="text-sm">{data.page.totalItems} track(s)</p>
  </div>

  {#each data.tracks as track, i}
    <TrackListItem
      {track}
      onPlayClicked={() => {
        musicManager.clearQueue();
        for (const track of data.tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track), false);
        }
        musicManager.setQueueIndex(i);
        musicManager.requestPlay();
      }}
    >
      <QuickAddButton
        show={!!(data.user && data.user.quickPlaylist)}
        {track}
        {apiClient}
        isInQuickPlaylist={(trackId) => {
          if (!data.quickPlaylistIds) return false;
          return !!data.quickPlaylistIds.find((v) => v === trackId);
        }}
      />

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={cn(
            buttonVariants({ variant: "ghost", size: "icon-lg" }),
            "rounded-full",
          )}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="end">
          <DropdownMenu.Group>
            <DropdownMenu.Item class="p-0">
              <!-- svelte-ignore a11y_invalid_attribute -->
              <a
                class="flex h-full w-full items-center gap-2 px-3 py-2"
                href="#"
              >
                <DiscAlbum size="16" />
                Go to Album
              </a>
            </DropdownMenu.Item>
            <DropdownMenu.Item class="p-0">
              <!-- svelte-ignore a11y_invalid_attribute -->
              <a
                class="flex h-full w-full items-center gap-2 px-3 py-2"
                href="#"
              >
                <Users size="16" />
                Go to Artist
              </a>
            </DropdownMenu.Item>
            <DropdownMenu.Item
              onSelect={async () => {
                if (!data.userPlaylists) return;

                await openAddToPlaylist({
                  apiClient,
                  playlists: data.userPlaylists,
                  track,
                });
                await invalidateAll();
              }}
            >
              <Users size="16" />
              Save to Playlist
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </TrackListItem>
  {/each}
</div>

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
