<script lang="ts">
  import { onMount } from "svelte";

  const { data } = $props();

  let files = $state<File[]>([]);
  let fileSelector = $state<HTMLInputElement>();

  let confirmAlbumDeletionDialog = $state<HTMLDialogElement>();
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

{#each data.tracks as track (track.id)}
  <p>{track.name}</p>
{/each}

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
