<script lang="ts">
  import type { Track } from "$lib/api/types";
  import { Play } from "lucide-svelte";
  import type { Snippet } from "svelte";

  type Props = {
    track: Track;
    children?: Snippet;

    onPlayClicked?: () => void;
  };

  const { track, children, onPlayClicked }: Props = $props();
</script>

<div class="flex items-center gap-2 border-b py-2 pr-2">
  <div class="group relative">
    <img
      class="inline-flex aspect-square w-14 min-w-14 items-center justify-center rounded border object-cover text-xs"
      src={track.coverArt.small}
      alt="cover"
      loading="lazy"
    />
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
  </div>
  <div class="flex flex-grow flex-col">
    <div class="flex items-center gap-1">
      <p class="line-clamp-1 w-fit text-sm font-medium" title={track.name}>
        {track.name}
      </p>

      <p>•</p>

      <a
        class="line-clamp-1 text-xs font-light hover:underline"
        title={track.artistName}
        href={`/artists/${track.artistId}`}
      >
        {track.artistName}
      </a>
    </div>

    <p class="line-clamp-1 text-xs font-light">
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
