<script lang="ts">
  import { applyAction, enhance } from "$app/forms";
  import { invalidate, invalidateAll } from "$app/navigation";
  import { Track } from "$lib/api/types.js";
  import { musicManager } from "$lib/music-manager.js";
  import { trackToMusicTrack } from "$lib/utils.js";

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

<!-- <button
        class="text-blue-400"
        onclick={() => {
          musicManager.clearQueue();
          musicManager.addTrackToQueue(trackToMusicTrack(track), true);
        }}>Play</button
      > -->

<!-- <p>
          {#if track.number}
            <span>{track.number} - </span>
          {/if}
          <span>{track.name}</span>
        </p>
        <p>Tags: {track.tags.join(", ")}</p> -->

<!-- -->

<div class="flex flex-col">
  {#each data.tracks as track (track.id)}
    <div class="border-b">
      <p>{track.name}</p>
      <div class="px-4">
        <p>Number: {track.number == null ? "NULL" : track.number}</p>
        <p>Year: {track.year == null ? "NULL" : track.year}</p>
        <p>Tags: {track.tags.join(", ")}</p>

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
  class="w-[400px] rounded bg-[--bg-color] p-4 text-[--fg-color] backdrop:bg-black/45"
  bind:this={editTrackDialog}
>
  <form
    class="flex flex-col"
    action="?/editTrack"
    method="post"
    use:enhance={() => {
      return async ({ result }) => {
        await applyAction(result);

        if (result.type === "success") {
          editTrackDialog?.close();
          await invalidateAll();
        }
      };
    }}
  >
    <input
      class="text-black"
      name="trackId"
      value={editTrack?.id ?? ""}
      type="hidden"
    />
    <input
      class="text-black"
      name="trackName"
      value={editTrack?.name}
      type="text"
    />
    <input
      class="text-black"
      name="trackTags"
      value={editTrack?.tags.join(",")}
      type="text"
    />

    <div>
      <label for="trackNumber">Number</label>
      <input
        class="w-full text-black"
        id="trackNumber"
        name="trackNumber"
        value={editTrack?.number === 0 ? undefined : editTrack?.number}
        type="number"
      />
    </div>

    <div>
      <label for="trackYear">Year</label>
      <input
        class="w-full text-black"
        id="trackYear"
        name="trackYear"
        value={editTrack?.year === 0 ? undefined : editTrack?.year}
        type="number"
      />
    </div>

    <button
      type="button"
      onclick={() => {
        editTrackDialog?.close();
      }}>Close</button
    >

    <button>Save</button>
  </form>
</dialog>
