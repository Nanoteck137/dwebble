<script lang="ts">
  import { Track } from "$lib/api/types";
  import { musicManager } from "$lib/music-manager";
  import { cn, formatTime, trackToMusicTrack } from "$lib/utils";
  import { Edit, EllipsisVertical, Play, Trash } from "lucide-svelte";
  import {
    Breadcrumb,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";

  const { data } = $props();

  let confirmAlbumDeletionDialog = $state<HTMLDialogElement>();

  let deleteTrack = $state<Track>();
  let confirmTrackDeletionDialog = $state<HTMLDialogElement>();

  function openDeleteConfirm(track: Track) {
    deleteTrack = track;
    confirmTrackDeletionDialog?.showModal();
  }
</script>

<Breadcrumb.Root>
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/albums/{data.album.id}">
        {data.album.name}
      </Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page>Edit</Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<div class="flex gap-2">
  <img
    class="aspect-square w-48 rounded object-cover"
    src={data.album.coverArt.medium}
    alt="Album Cover Art"
  />

  <div class="flex flex-col py-2">
    <p class="font-bold">
      {data.album.name}
    </p>
    <a
      class="w-fit text-xs hover:underline"
      href="/artists/{data.album.artistId}"
    >
      {data.album.artistName}
    </a>

    <div class="flex-grow"></div>

    <div class="flex gap-2">
      <a class="text-sm hover:underline" href="edit/details">
        Edit Album Details
      </a>
      {#if data.tracks.length > 0}
        <a class="text-sm hover:underline" href="edit/tracks">Edit Tracks</a>
      {/if}
      <a class="text-sm hover:underline" href="edit/import">Import Tracks</a>
    </div>
  </div>
</div>

<div class="flex flex-col">
  {#each data.tracks as track (track.id)}
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <div class="flex flex-grow flex-col">
        <p class="text-sm font-medium" title={track.name}>
          {#if track.number}
            <span>{track.number}.</span>
          {/if}
          {track.name}
        </p>
        <p class="text-xs" title={track.artistName}>
          Artist: <a class="hover:underline" href="/artists/{track.artistId}"
            >{track.artistName}</a
          >
        </p>
        {#if track.year}
          <p class="text-xs">Year: {track.year}</p>
        {/if}
        {#if track.tags.length > 0}
          <p class="text-xs">Tags: {track.tags.join(", ")}</p>
        {/if}
      </div>
      <div class="flex items-center gap-2">
        <p class="text-xs">
          {formatTime(track.duration ?? 0)}
        </p>

        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item>
                <button
                  class="flex w-full items-center gap-2"
                  onclick={() => {
                    musicManager.clearQueue();
                    musicManager.addTrackToQueue(
                      trackToMusicTrack(track),
                      true,
                    );
                  }}
                >
                  <Play size="16" />
                  Play
                </button>
              </DropdownMenu.Item>

              <DropdownMenu.Item>
                <button
                  class="flex w-full items-center gap-2"
                  onclick={() => {
                    openDeleteConfirm(track);
                  }}
                >
                  <Trash size="16" />
                  Remove
                </button>
              </DropdownMenu.Item>

              <DropdownMenu.Item>
                <a
                  class="flex w-full items-center gap-2"
                  href="edit/tracks#track-{track.id}"
                >
                  <Edit size="16" />
                  Edit
                </a>
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  {/each}
</div>

<dialog
  class="rounded bg-[--bg-color] p-4 text-[--fg-color] backdrop:bg-black/45"
  bind:this={confirmAlbumDeletionDialog}
>
  <p>Are you sure?</p>

  <button
    onclick={() => {
      confirmAlbumDeletionDialog?.close();
    }}>Close</button
  >

  <form action="?/deleteAlbum" method="post">
    <input name="albumId" value={data.album.id} type="hidden" />
    <button>DELETE</button>
  </form>
</dialog>

<dialog
  class="rounded bg-[--bg-color] p-4 text-[--fg-color] backdrop:bg-black/45"
  bind:this={confirmTrackDeletionDialog}
>
  <p>Are you sure?</p>

  <button
    onclick={() => {
      confirmTrackDeletionDialog?.close();
    }}>Close</button
  >

  <form action="?/deleteTrack" method="post">
    <input name="trackId" value={deleteTrack?.id} type="hidden" />
    <button>DELETE</button>
  </form>
</dialog>
