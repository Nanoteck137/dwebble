<script lang="ts">
  import { musicManager } from "$lib/music-manager";
  import { trackToMusicTrack } from "$lib/utils";
  import { EllipsisVertical, Play, Plus } from "lucide-svelte";
  import type { PageData } from "./$types";
  import { enhance } from "$app/forms";

  export let data: PageData;
</script>

<p>Playlist Page (W.I.P)</p>
{data.id}

<div class="flex flex-col">
  <button
    onclick={() => {
      musicManager.clearQueue();
      for (const track of data.items) {
        musicManager.addTrackToQueue(trackToMusicTrack(track));
      }
    }}>Play</button
  >
  {#each data.items as track}
    <div class="group flex items-center gap-2 border-b p-2 pr-4">
      <div class="group relative">
        <img
          class="aspect-square w-14 min-w-14 rounded object-cover"
          src={track.coverArt}
          alt={`${track.name} Cover Art`}
          loading="lazy"
        />
        <button
          class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded bg-[--overlay-bg] group-hover:flex`}
          onclick={() => {
            musicManager.clearQueue();
            musicManager.addTrackToQueue(trackToMusicTrack(track));
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
          <a
            class="line-clamp-1 w-fit text-sm hover:underline"
            title={track.artistName}
            href={`/artists/${track.artistId}`}
          >
            {track.artistName}
          </a>
        </div>

        <p class="line-clamp-1 text-xs">
          {track.genres.join(", ")}
        </p>

        <p class="line-clamp-1 text-xs">
          {#if track.tags.length > 0}
            {track.tags.join(", ")}
          {:else}
            No Tags
          {/if}
        </p>
      </div>
      <div class="flex items-center">
        <button>
          <EllipsisVertical size="28" />
        </button>
      </div>
    </div>
  {/each}
</div>
