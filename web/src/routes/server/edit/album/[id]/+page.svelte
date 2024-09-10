<script lang="ts">
  import { applyAction, enhance } from "$app/forms";
  import { invalidateAll } from "$app/navigation";
  import { Track } from "$lib/api/types.js";
  import { musicManager } from "$lib/music-manager.js";
  import { formatTime, trackToMusicTrack } from "$lib/utils.js";
  import { EllipsisVertical, Play } from "lucide-svelte";
  import TrackEdit from "./TrackEdit.svelte";

  const { data } = $props();

  let editAlbumName = $state(false);

  let files = $state<File[]>([]);
  let fileSelector = $state<HTMLInputElement>();

  let confirmAlbumDeletionDialog = $state<HTMLDialogElement>();

  let deleteTrack = $state<Track>();
  let confirmTrackDeletionDialog = $state<HTMLDialogElement>();

  function openDeleteConfirm(track: Track) {
    deleteTrack = track;
    confirmTrackDeletionDialog?.showModal();
  }

  let editAlbumArtist = $state<HTMLDialogElement>();

  let editTrack = $state<Track>();
  let editTrackDialog = $state<HTMLDialogElement>();

  function openTrackEditor(track: Track) {
    editTrack = track;
    editTrackDialog?.showModal();
  }
</script>

<p>Edit Album</p>

<button
  class="rounded bg-red-500 px-4 py-2 text-white hover:bg-red-400 active:scale-95"
  onclick={() => {
    confirmAlbumDeletionDialog?.showModal();
  }}>Delete Album</button
>

<div>
  {#if editAlbumName}
    <form action="?/editAlbumName" method="post">
      <input name="albumId" value={data.album.id} type="hidden" />
      <input
        class="text-black"
        name="albumName"
        value={data.album.name}
        type="text"
      />
      <button>Save</button>
      <button
        type="button"
        onclick={() => {
          editAlbumName = !editAlbumName;
        }}>Edit</button
      >
    </form>
  {:else}
    <div class="flex gap-2">
      <p>{data.album.name}</p>
      <button
        onclick={() => {
          editAlbumName = !editAlbumName;
        }}>Edit</button
      >
    </div>
  {/if}
</div>

<p>Album Artist: {data.album.artistName}</p>

<button
  onclick={() => {
    editAlbumArtist?.showModal();
  }}>Change Album Artist</button
>

<p>{data.album.coverArt}</p>

<button
  onclick={() => {
    editTrackDialog?.showModal();
  }}>Edit Tracks</button
>

<div class="flex flex-col">
  {#each data.tracks as track (track.id)}
    <div class="flex items-center gap-2 border-b p-2">
      <div class="flex flex-grow flex-col">
        <p title={track.name}>
          {#if track.number}
            <span>{track.number}.</span>
          {/if}
          {track.name}
        </p>
        <!-- <a title={track.artistName} href={`/artist/${track.artistId}`}>
          {track.artistName}
        </a> -->
      </div>
      <div class="flex items-center gap-2">
        <p class="">
          {formatTime(track.duration ?? 0)}
        </p>
        <button class="rounded-full p-1 hover:bg-black/20">
          <EllipsisVertical size="30" />
        </button>
      </div>
    </div>
  {/each}
</div>

<!-- <div class="border-b">
      <p>{track.name}</p>
      <div class="px-4">
        <p>Number: {track.number == null ? "NULL" : track.number}</p>
        <p>Year: {track.year == null ? "NULL" : track.year}</p>
        <p>Tags: {track.tags.join(", ")}</p>

        <button
          class="rounded bg-purple-400 px-2 py-1"
          onclick={() => {
            musicManager.clearQueue();
            musicManager.addTrackToQueue(trackToMusicTrack(track), true);
          }}>Play</button
        >

        <button
          class="rounded bg-blue-400 px-2 py-1"
          onclick={() => {
            openTrackEditor(track);
          }}>Edit</button
        >

        <button
          class="rounded bg-red-400 px-2 py-1"
          onclick={() => {
            openDeleteConfirm(track);
          }}>Delete</button
        >
      </div>
    </div> -->

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

<dialog
  class="rounded bg-[--bg-color] p-4 text-[--fg-color] backdrop:bg-black/45"
  bind:this={editAlbumArtist}
>
  <form action="?/editAlbumArtist" method="post">
    <input name="albumId" value={data.album.id} type="hidden" />
    <input name="artistName" value={data.album.artistName} type="text" />

    <button
      type="button"
      onclick={() => {
        editAlbumArtist?.close();
      }}>Close</button
    >

    <button>Change Artist</button>
  </form>
</dialog>

<dialog
  class="w-full rounded bg-[--bg-color] p-4 text-[--fg-color] backdrop:bg-black/45"
  bind:this={editTrackDialog}
>
  <!-- use:enhance={() => {
      return async ({ result, formElement }) => {
        await applyAction(result);

        if (result.type === "success") {
          editTrackDialog?.close();
          await invalidateAll();
          formElement.reset();
          editTrack = undefined;
        }
      };
    }} -->

  <form class="flex flex-col" action="?/editTracks" method="post">
    <div class="flex flex-col gap-2">
      {#each data.tracks as track (track.id)}
        <TrackEdit {track} />
      {/each}
    </div>

    <div>
      <button>Save</button>
      <button
        type="button"
        onclick={() => {
          editTrackDialog?.close();
        }}>Close</button
      >
    </div>
  </form>
</dialog>
