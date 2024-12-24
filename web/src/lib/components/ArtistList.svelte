<script lang="ts">
  import type { ArtistInfo } from "$lib/api/types";
  import { cn } from "$lib/utils";

  type Props = {
    class?: string;
    artists: ArtistInfo[];
  };

  const { class: className, artists }: Props = $props();
</script>

<p class={cn("line-clamp-1 text-xs", className)}>
  {#each artists as artist, i}
    {#if i > 0}
      {" & "}
    {/if}
    <a class="hover:underline" href="/artists/{artist.id}">
      {artist.name.default}
    </a>
  {/each}
</p>
