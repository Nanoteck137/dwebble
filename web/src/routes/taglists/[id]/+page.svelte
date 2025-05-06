<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { getApiClient, handleApiError } from "$lib";
  import Spacer from "$lib/components/Spacer.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { Breadcrumb, Button, Input, Pagination } from "@nanoteck137/nano-ui";
  import { Filter } from "lucide-svelte";

  const { data } = $props();
  const musicManager = getMusicManager();
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
<!--TODO(patrik): Fix this-->
<!-- <Button
  onclick={async () => {
    // TODO(patrik): Better title and desc
    const confirmed = await openConfirm({ title: "Are you sure?" });

    if (confirmed) {
      const res = await apiClient.deleteTaglist(data.taglist.id);
      if (!res.success) {
        handleApiError(res.error);
        return;
      }

      goto("/taglists");
    }
  }}
>
  Delete Taglist
</Button> -->

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

<TrackListHeader
  name={data.taglist.name}
  onPlay={async (shuffle) => {
    await musicManager.queueRequest(
      {
        type: "addTaglist",
        taglistId: data.taglist.id,
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
        type: "addTaglist",
        taglistId: data.taglist.id,
      },
      { queueIndexToTrackId: trackId },
    );
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
