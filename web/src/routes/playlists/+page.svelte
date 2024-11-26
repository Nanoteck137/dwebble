<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";
  import type { PageData } from "./$types";

  interface Props {
    data: PageData;
  }

  let { data }: Props = $props();
</script>

<Button href="/playlists/new">New Playlist</Button>

<div class="flex flex-col gap-2">
  {#each data.playlists as playlist}
    <a href={`/playlists/${playlist.id}`}>{playlist.name}</a>
  {/each}
</div>
