<script lang="ts">
  import type { Album, Artist } from "$lib/api/types";
  import type { Snippet } from "svelte";

  type Props = {
    artist: Artist;
    link?: boolean;
    children?: Snippet;
  };

  const { artist, link, children }: Props = $props();
</script>

<div class="flex items-center gap-2 border-b py-2 pr-2">
  <div class="group relative">
    <img
      class="inline-flex aspect-square w-14 min-w-14 items-center justify-center rounded border object-cover text-xs"
      src={artist.picture.small}
      alt="cover"
      loading="lazy"
    />
  </div>
  <div class="flex flex-grow flex-col">
    <div class="flex items-center gap-1">
      {#if link}
        <a
          class="line-clamp-1 w-fit text-sm font-medium hover:underline"
          title={artist.name.default}
          href="/artists/{artist.id}"
        >
          {artist.name.default}
        </a>
      {:else}
        <p
          class="line-clamp-1 w-fit text-sm font-medium"
          title={artist.name.default}
        >
          {artist.name.default}
        </p>
      {/if}
    </div>
  </div>
  <div class="flex items-center">
    {@render children?.()}
  </div>
</div>
