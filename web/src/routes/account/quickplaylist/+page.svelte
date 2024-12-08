<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

<p>Change Quick Playlist</p>

{#each data.playlists as playlist}
  <form method="post">
    <input name="playlistId" value={playlist.id} type="hidden" />
    <Button type="submit" variant="ghost">{playlist.name}</Button>
  </form>
{/each}
