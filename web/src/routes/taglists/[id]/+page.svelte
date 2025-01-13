<script lang="ts">
  import { enhance } from "$app/forms";
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { createApiClient, openConfirm } from "$lib";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { musicManager } from "$lib/music-manager.svelte.js";
  import { cn, shuffle, trackToMusicTrack } from "$lib/utils.js";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Input,
    Pagination,
  } from "@nanoteck137/nano-ui";
  import {
    EllipsisVertical,
    Filter,
    ListPlus,
    Pencil,
    Play,
    Plus,
    Shuffle,
  } from "lucide-svelte";

  const { data } = $props();
  const apiClient = createApiClient(data);
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/taglists">Taglists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.taglist.name}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

{#if data.filterError}
  <p class="text-red-400">{data.filterError}</p>
{/if}

<Button href="/taglists/{data.taglist.id}/edit">Edit Taglist</Button>
<Button
  onclick={async () => {
    // TODO(patrik): Better title and desc
    const confirmed = await openConfirm({ title: "Are you sure?" });

    if (confirmed) {
      const res = await apiClient.deleteTaglist(data.taglist.id);
      if (!res.success) {
        // TODO(patrik): Toast
        throw res.error.message;
      }

      goto("/taglists");
    }
  }}>Delete Taglist</Button
>

<form
  method="GET"
  onsubmit={() => {
    // TODO(patrik): Temp Fix
    invalidateAll();
  }}
>
  <div class="flex flex-col gap-2">
    <Input
      type="text"
      name="sort"
      placeholder="Sort"
      value={data.sort ?? ""}
    />
  </div>

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
