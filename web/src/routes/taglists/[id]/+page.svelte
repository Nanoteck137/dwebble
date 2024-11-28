<script lang="ts">
  import { enhance } from "$app/forms";
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { musicManager } from "$lib/music-manager.js";
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
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <div class="group relative">
        <img
          class="aspect-square w-14 min-w-14 rounded object-cover"
          src={track.coverArt.small}
          alt={`${track.name} Cover Art`}
          loading="lazy"
        />
        <button
          class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded bg-black/80 group-hover:flex`}
          onclick={() => {
            musicManager.clearQueue();
            for (const track of data.tracks) {
              musicManager.addTrackToQueue(trackToMusicTrack(track), false);
            }
            musicManager.setQueueIndex(i);

            musicManager.requestPlay();
          }}
        >
          <Play size="25" />
        </button>
      </div>
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <p class="line-clamp-1 w-fit font-medium" title={track.name}>
            {track.name}
          </p>
        </div>

        <a
          class="line-clamp-1 w-fit text-sm font-light hover:underline"
          title={track.artistName}
          href={`/artists/${track.artistId}`}
        >
          {track.artistName}
        </a>

        <p class="line-clamp-1 text-xs font-light">
          {#if track.tags.length > 0}
            {track.tags.join(", ")}
          {:else}
            No Tags
          {/if}
        </p>
      </div>
      <div class="flex items-center">
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
              <DropdownMenu.Item>
                <form
                  class="jusitfy-center flex items-center"
                  action="?/quickAddToPlaylist"
                  method="post"
                  use:enhance
                >
                  <input type="hidden" name="trackId" value={track.id} />
                  <button
                    class="flex w-full items-center gap-2 py-1"
                    title="Quick Add"
                  >
                    <Plus size="16" />
                    Quick add to Playlist
                  </button>
                </form>
              </DropdownMenu.Item>
              <DropdownMenu.Item>
                <a
                  class="flex h-full w-full items-center gap-2 py-1"
                  href="/albums/{track.albumId}"
                >
                  <Pencil size="16" />
                  Go to Album
                </a>
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
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
