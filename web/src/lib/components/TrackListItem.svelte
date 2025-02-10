<script lang="ts">
  import type { Track } from "$lib/api/types";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import Image from "$lib/components/Image.svelte";
  import { Play } from "lucide-svelte";
  import type { Snippet } from "svelte";

  type Props = {
    showNumber?: boolean;
    track: Track;
    children?: Snippet;

    onPlayClicked?: () => void;
  };

  const { showNumber, track, children, onPlayClicked }: Props = $props();
</script>

<div class="flex items-center gap-2 py-2 pr-2">
  <div class="group relative">
    {#if showNumber}
      <div
        class="group flex min-h-10 min-w-10 flex-col items-end justify-center"
      >
        <p class=" text-right text-sm font-medium group-hover:hidden">
          {track.number}.
        </p>
        {#if onPlayClicked}
          <button
            class={`hidden group-hover:block`}
            onclick={() => {
              onPlayClicked?.();
            }}
          >
            <Play size="25" />
          </button>
        {/if}
      </div>
    {:else}
      <Image class="w-14 min-w-14" src={track.coverArt.small} alt="cover" />
      {#if onPlayClicked}
        <button
          class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded border bg-black/80 group-hover:flex`}
          onclick={() => {
            onPlayClicked?.();
          }}
        >
          <Play size="25" />
        </button>
      {/if}
    {/if}
  </div>
  <div class="flex flex-grow flex-col">
    <div class="flex items-center gap-1">
      <p
        class="line-clamp-1 w-fit text-sm font-medium"
        title={track.name.default}
      >
        {track.name.default}
      </p>
    </div>

    <ArtistList class="text-muted-foreground" artists={track.allArtists} />

    <p class="line-clamp-1 text-xs text-muted-foreground">
      {#if track.tags.length > 0}
        {track.tags.join(", ")}
      {:else}
        No Tags
      {/if}
    </p>
  </div>
  <div class="flex items-center">
    {@render children?.()}
  </div>
</div>
