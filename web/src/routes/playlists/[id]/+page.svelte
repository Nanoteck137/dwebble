<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { getApiClient, handleApiError } from "$lib";
  import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { Breadcrumb, DropdownMenu, Pagination } from "@nanoteck137/nano-ui";
  import { Pencil, Trash } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const musicManager = getMusicManager();
  const apiClient = getApiClient();

  let openConfirmDeleteAlbum = $state(false);
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
  onPlay={async (shuffle) => {
    await musicManager.queueRequest(
      { type: "addPlaylist", playlistId: data.playlist.id },
      { shuffle },
    );
  }}
>
  {#snippet more()}
    <DropdownMenu.Group>
      <DropdownMenu.Item onSelect={() => {}}>
        <Pencil />
        Edit Playlist
      </DropdownMenu.Item>

      <DropdownMenu.Item
        onSelect={() => {
          openConfirmDeleteAlbum = true;
        }}
      >
        <Trash />
        Delete Playlist
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
  onPlay={async (trackId) => {
    await musicManager.queueRequest(
      { type: "addPlaylist", playlistId: data.playlist.id },
      { queueIndexToTrackId: trackId },
    );
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

<ConfirmModal
  bind:open={openConfirmDeleteAlbum}
  removeTrigger
  confirmDelete
  onResult={async () => {
    const res = await apiClient.deletePlaylist(data.playlist.id);
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Successfully deleted album");
    goto("/playlists", { invalidateAll: true });
  }}
/>
