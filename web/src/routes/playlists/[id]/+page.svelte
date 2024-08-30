<script lang="ts">
  import { applyAction, enhance } from "$app/forms";
  import { invalidate } from "$app/navigation";
  import { Track } from "$lib/api/types";
  import { musicManager } from "$lib/music-manager";
  import { trackToMusicTrack } from "$lib/utils";
  import { BookmarkMinus, Check, EllipsisVertical, Play } from "lucide-svelte";

  const { data } = $props();

  let popupOpen = $state<string>();
  let selected = $state<string[]>([]);

  function isSelected(chapterSlug: string) {
    for (let i = 0; i < selected.length; i++) {
      if (selected[i] === chapterSlug) {
        return true;
      }
    }

    return false;
  }
</script>

<p>Playlist Page (W.I.P)</p>
{data.id}
<p>{data.items.length} items</p>

<div class="flex flex-col">
  <button
    onclick={() => {
      musicManager.clearQueue();
      for (const track of data.items) {
        musicManager.addTrackToQueue(trackToMusicTrack(track));
      }
    }}>Play</button
  >

  {#snippet normalItem(track: Track, isEditing: boolean)}
    <div class="flex items-center gap-2 border-b p-2 pr-4">
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
      {#if !isEditing}
        <div class="relative flex items-center">
          <button
            onclick={() => {
              if (popupOpen === track.id) {
                popupOpen = undefined;
              } else {
                popupOpen = track.id;
              }
            }}
          >
            <EllipsisVertical size="28" />
          </button>

          {#if popupOpen === track.id}
            <div
              class="absolute right-2 top-8 z-50 rounded border border-gray-400 bg-gray-800 p-4"
            >
              <form action="?/remove" method="post" use:enhance>
                <input name="playlistId" value={data.id} type="hidden" />
                <input name="tracks[]" value={track.id} type="hidden" />
                <button class="flex items-center gap-2">
                  <BookmarkMinus size="24" />
                  <span>Remove</span>
                </button>
              </form>
            </div>
          {/if}
        </div>
      {/if}
      <div class="flex items-center">
        <button
          class="flex h-6 w-6 items-center justify-center rounded border"
          onclick={() => {
            if (isSelected(track.id)) {
              selected = selected.filter((t) => t !== track.id);
            } else {
              selected.push(track.id);
            }
          }}
        >
          {#if isSelected(track.id)}
            <Check size="18" />
          {/if}
        </button>
      </div>
    </div>
  {/snippet}

  {#snippet editItem(track: Track)}
    <div class="flex items-center gap-2 border-b p-2 pr-4">
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
      <button
        class="flex w-full items-center justify-between"
        onclick={() => {
          if (isSelected(track.id)) {
            selected = selected.filter((t) => t !== track.id);
          } else {
            selected.push(track.id);
          }
        }}
      >
        <div class="flex flex-col">
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

          <p class="line-clamp-1 text-start text-xs">
            {track.genres.join(", ")}
          </p>

          <p class="line-clamp-1 text-start text-xs">
            {#if track.tags.length > 0}
              {track.tags.join(", ")}
            {:else}
              No Tags
            {/if}
          </p>
        </div>
        <div class="flex items-center">
          <div class="flex h-6 w-6 items-center justify-center rounded border">
            {#if isSelected(track.id)}
              <Check size="18" />
            {/if}
          </div>
        </div>
      </button>
    </div>
  {/snippet}

  {#each data.items as track (track.id)}
    {#if selected.length > 0}
      {@render editItem(track)}
    {:else}
      {@render normalItem(track, selected.length > 0)}
    {/if}
  {/each}
</div>

{#if selected.length > 0}
  <div class="fixed bottom-4 left-1/2 -translate-x-1/2 bg-red-500 p-4">
    <form
      action="?/remove"
      method="post"
      use:enhance={() => {
        return async ({ result, update }) => {
          if (result.type === "success") {
            selected = [];
          }

          await applyAction(result);
          await update();
        };
      }}
    >
      <input name="playlistId" value={data.id} type="hidden" />
      {#each selected as sel}
        <input name="tracks[]" value={sel} type="hidden" />
      {/each}

      <button>Delete Selected</button>
    </form>
  </div>
{/if}
