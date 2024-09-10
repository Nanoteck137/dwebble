<script lang="ts">
  import { EllipsisVertical, Play, Plus, X } from "lucide-svelte";
  import type { PageData } from "./$types";
  import { musicManager } from "$lib/music-manager";
  import { trackToMusicTrack } from "$lib/utils";
  import { enhance } from "$app/forms";

  export let data: PageData;
</script>

<p>Num Tracks: {data.tracks.length}</p>

<form class="px-4" method="GET">
  <div class="flex flex-col gap-2">
    <input
      class="border-1 h-8 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-sm text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
      type="text"
      name="filter"
      placeholder="Filter"
      value={data.filter ?? ""}
    />
    <input
      class="border-1 h-8 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-sm text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
      type="text"
      name="sort"
      placeholder="Sort"
      value={data.sort ?? ""}
    />
  </div>

  {#if data.filterError}
    <p class="text-red-400">{data.filterError}</p>
  {/if}
  {#if data.sortError}
    <p class="text-red-400">{data.sortError}</p>
  {/if}
  <button>Filter</button>
</form>

<div class="flex flex-col">
  <div class="flex gap-2 px-4">
    <button
      class="rounded bg-blue-500 px-4 py-2 text-lg hover:bg-blue-600"
      onclick={() => {
        musicManager.clearQueue();
        for (const track of data.tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track));
        }
      }}>Play</button
    >
    <form action="?/newPlaylist" method="post">
      <input name="filter" value={data.filter} type="hidden" />
      <input name="sort" value={data.sort} type="hidden" />
      <button class="rounded bg-pink-500 px-4 py-2 text-lg hover:bg-pink-600">
        New playlist from tracks
      </button>
    </form>
  </div>
  {#each data.tracks as track, i}
    <div class="flex items-center gap-2 border-b p-2 pr-4">
      <div class="group relative">
        <img
          class="aspect-square w-14 min-w-14 rounded object-cover"
          src={track.coverArtUrl}
          alt={`${track.name} Cover Art`}
          loading="lazy"
        />
        <button
          class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded bg-[--overlay-bg] group-hover:flex`}
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
        <!-- TODO(patrik): Move to dropdown menu -->
        <form
          class="jusitfy-center flex items-center"
          action="?/quickAddToPlaylist"
          method="post"
          use:enhance
        >
          <input type="hidden" name="trackId" value={track.id} />
          <button title="Quick Add"><Plus size="28" /></button>
        </form>

        <button>
          <EllipsisVertical size="28" />
        </button>
      </div>
    </div>
  {/each}
</div>
