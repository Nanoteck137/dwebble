<script lang="ts">
  import { Track } from "$lib/api/types.js";
  import { musicManager } from "$lib/music-manager.js";
  import { formatTime, trackToMusicTrack } from "$lib/utils.js";
  import { Edit, EllipsisVertical, Play, Trash } from "lucide-svelte";
  import { DropdownMenu } from "bits-ui";

  const { data } = $props();

  let confirmAlbumDeletionDialog = $state<HTMLDialogElement>();

  let deleteTrack = $state<Track>();
  let confirmTrackDeletionDialog = $state<HTMLDialogElement>();

  function openDeleteConfirm(track: Track) {
    deleteTrack = track;
    confirmTrackDeletionDialog?.showModal();
  }
</script>

<div class="flex gap-2 p-2">
  <img
    class="aspect-square w-60 rounded object-cover"
    src={data.album.coverArt.medium}
    alt="Album Cover Art"
  />

  <div class="flex flex-col py-2">
    <a
      class="text-3xl font-medium hover:underline"
      href="/albums/{data.album.id}"
    >
      {data.album.name}
    </a>
    <a class="text-lg hover:underline" href="/artists/{data.album.artistId}">
      {data.album.artistName}
    </a>

    <div class="flex-grow"></div>

    <div class="flex gap-2">
      <a class="text-sm hover:underline" href="{data.album.id}/details">
        Edit Album Details
      </a>
      <a class="text-sm hover:underline" href="{data.album.id}/tracks">
        Edit Tracks
      </a>
      <a class="text-sm hover:underline" href="{data.album.id}/import">
        Import Tracks
      </a>
    </div>
  </div>
</div>

<!-- <div class="flex flex-col">
  <p>Edit Album</p>

  <a href="{data.album.id}/import">Import Tracks</a>
  <a href="{data.album.id}/details">Edit Album Details</a>
  <a href="{data.album.id}/tracks">Edit Tracks</a>

  <button
    class="rounded bg-red-500 px-4 py-2 text-white hover:bg-red-400 active:scale-95"
    onclick={() => {
      confirmAlbumDeletionDialog?.showModal();
    }}
  >
    Delete Album
  </button>

  <div class="flex gap-2">
    <p>{data.album.name}</p>
  </div>

  <p>Album Artist: {data.album.artistName}</p>

  <p>{data.album.coverArt}</p>
</div> -->

<div class="flex flex-col">
  {#each data.tracks as track (track.id)}
    <div class="flex items-center gap-2 border-b py-2 pl-2 pr-4">
      <div class="flex flex-grow flex-col">
        <p title={track.name}>
          {#if track.number}
            <span>{track.number}.</span>
          {/if}
          {track.name}
        </p>
        <a
          class="text-xs"
          title={track.artistName}
          href="/artists/{track.artistId}"
        >
          Artist: {track.artistName}
        </a>
        {#if track.year}
          <p class="text-xs">Year: {track.year}</p>
        {/if}
        {#if track.tags.length > 0}
          <p class="text-xs">Tags: {track.tags.join(", ")}</p>
        {/if}
        <!--  -->
      </div>
      <div class="flex items-center gap-2">
        <p class="">
          {formatTime(track.duration ?? 0)}
        </p>
        <DropdownMenu.Root disableFocusFirstItem={true}>
          <DropdownMenu.Trigger>
            <EllipsisVertical size="24" />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content
            class="w-full max-w-[240px] rounded border border-gray-600 bg-[--bg-color] px-2 py-2"
          >
            <DropdownMenu.Item
              class="select-none rounded px-2 py-2 data-[highlighted]:bg-gray-400 data-[highlighted]:text-black"
            >
              <button
                class="flex h-full w-full items-center gap-2"
                onclick={() => {
                  musicManager.clearQueue();
                  musicManager.addTrackToQueue(trackToMusicTrack(track), true);
                }}
              >
                <Play size="20" />
                <p>Play</p>
              </button>
            </DropdownMenu.Item>
            <DropdownMenu.Item
              class="select-none rounded px-2 py-2 data-[highlighted]:bg-red-400 data-[highlighted]:text-black"
            >
              <button
                class="flex h-full w-full items-center gap-2"
                onclick={() => {
                  openDeleteConfirm(track);
                }}
              >
                <Trash size="20" />
                <p>Remove</p>
              </button>
            </DropdownMenu.Item>

            <DropdownMenu.Item
              class="select-none rounded px-2 py-2 data-[highlighted]:bg-red-400 data-[highlighted]:text-black"
            >
              <a
                class="flex h-full w-full items-center gap-2"
                href="{data.album.id}/tracks#track-{track.id}"
              >
                <Edit size="20" />
                <p>Edit</p>
              </a>
            </DropdownMenu.Item>
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
