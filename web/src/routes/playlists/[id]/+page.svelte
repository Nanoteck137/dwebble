<script lang="ts">
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Spacer from "$lib/components/Spacer.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { Breadcrumb, DropdownMenu, Pagination } from "@nanoteck137/nano-ui";
  import { Pencil } from "lucide-svelte";

  const { data } = $props();
  const musicManager = getMusicManager();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/playlists">Playlists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.playlist.name}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<TrackListHeader
  name={data.playlist.name}
  onPlay={async () => {
    await musicManager.clearQueue();
    await musicManager.addFromPlaylist(data.playlist.id);
    musicManager.requestPlay();
  }}
>
  {#snippet more()}
    <DropdownMenu.Group>
      <DropdownMenu.Item onSelect={() => {}}>
        <Pencil />
        Edit Playlist
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  {/snippet}
</TrackListHeader>

<Spacer size="md" />

<TrackList
  totalTracks={data.page.totalItems}
  tracks={data.items}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  onPlay={async () => {
    await musicManager.clearQueue();
    await musicManager.addFromPlaylist(data.playlist.id);
    musicManager.requestPlay();
  }}
  onTrackPlay={async (trackId) => {
    await musicManager.clearQueue();
    await musicManager.addFromPlaylist(data.playlist.id);
    await musicManager.setQueueIndex(
      musicManager.queue.items.findIndex((t) => t.track.id === trackId),
    );
    musicManager.requestPlay();
  }}
/>

<Spacer size="md" />

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
