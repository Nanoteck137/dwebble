<script lang="ts">
  import type { Artist } from "$lib/api/types";
  import Image from "$lib/components/Image.svelte";
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
    <Image class="w-14 min-w-14" src={artist.picture.small} alt="cover" />
  </div>
  <div class="flex flex-grow flex-col">
    <div class="flex items-center gap-1">
      {#if link}
        <a
          class="line-clamp-1 w-fit text-sm font-medium hover:underline"
          title={artist.name}
          href="/artists/{artist.id}"
        >
          {artist.name}
        </a>
      {:else}
        <p class="line-clamp-1 w-fit text-sm font-medium" title={artist.name}>
          {artist.name}
        </p>
      {/if}
    </div>
  </div>
  <div class="flex items-center">
    {@render children?.()}
  </div>
</div>
