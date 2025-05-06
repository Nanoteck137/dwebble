<script lang="ts">
  import type { Album } from "$lib/api/types";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import Image from "$lib/components/Image.svelte";
  import type { Snippet } from "svelte";

  type Props = {
    album: Album;
    link?: boolean;
    children?: Snippet;
  };

  const { album, link, children }: Props = $props();
</script>

<div class="flex items-center gap-2 border-b py-2 pr-2">
  <div class="group relative">
    <Image class="w-14 min-w-14" src={album.coverArt.small} alt="cover" />
  </div>
  <div class="flex flex-grow flex-col">
    <div class="flex items-center gap-1">
      {#if link}
        <a
          class="line-clamp-1 w-fit text-sm font-medium hover:underline"
          title={album.name}
          href="/albums/{album.id}"
        >
          {album.name}
        </a>
      {:else}
        <p class="line-clamp-1 w-fit text-sm font-medium" title={album.name}>
          {album.name}
        </p>
      {/if}

      <p>•</p>

      <ArtistList artists={album.artists} />
      <!-- <a
        class="line-clamp-1 text-xs font-light hover:underline"
        title={album.artistName}
        href={`/artists/${album.artistId}`}
      >
        {album.artistName}
      </a> -->
    </div>

    <p class="line-clamp-1 text-xs font-light">
      {#if album.year}
        {album.year}
      {/if}
    </p>
  </div>
  <div class="flex items-center">
    {@render children?.()}
  </div>
</div>
