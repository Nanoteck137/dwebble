<script lang="ts">
  import { Track } from "$lib/api/types.js";
  import { musicManager } from "$lib/music-manager.js";
  import { trackToMusicTrack } from "$lib/utils.js";
  import { onMount } from "svelte";

  const { data } = $props();

  let files = $state<File[]>([]);
  let fileSelector = $state<HTMLInputElement>();

  let confirmAlbumDeletionDialog = $state<HTMLDialogElement>();

  let deleteTrack = $state<Track>();
  let confirmTrackDeletionDialog = $state<HTMLDialogElement>();
</script>

<p>Edit Album</p>

<button
  class="rounded bg-red-500 px-4 py-2 text-white hover:bg-red-400 active:scale-95"
  onclick={() => {
    confirmAlbumDeletionDialog?.showModal();
  }}>Delete Album</button
>

<p>{data.album.name}</p>
<p>{data.album.artistName}</p>
<p>{data.album.coverArt}</p>

<div class="flex flex-col">
  {#each data.tracks as track (track.id)}
    <div class="flex gap-2">
      <button
        class="text-blue-400"
        onclick={() => {
          musicManager.clearQueue();
          musicManager.addTrackToQueue(trackToMusicTrack(track), true);
        }}>Play</button
      >

      <p>{track.name}</p>
      <button
        class="text-red-400"
        onclick={() => {
          deleteTrack = track;
          confirmTrackDeletionDialog?.showModal();
        }}>Delete</button
      >
    </div>
  {/each}
</div>

<button
  onclick={() => {
    fileSelector?.click();
  }}>Add Files</button
>

<form action="?/importTracks" method="post" enctype="multipart/form-data">
  <input name="albumId" value={data.album.id} type="hidden" />

  {#each files as file}
    <p>{file.name}</p>
  {/each}

  {#if files.length > 0}
    <button>Import Tracks</button>
  {/if}

  <input
    class="hidden"
    name="files"
    bind:this={fileSelector}
    type="file"
    multiple
    accept="audio/*"
    onchange={(e) => {
      console.log(e);
      const target = e.target as HTMLInputElement;

      if (!target.files) {
        return;
      }

      let list = [];

      for (let i = 0; i < target.files.length; i++) {
        const item = target.files.item(i);
        if (!item) {
          continue;
        }

        list.push(item);
      }

      list.sort((a, b) => a.name.localeCompare(b.name));

      files = list;
    }}
  />
</form>

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
  <p>{deleteTrack?.name}</p>

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
