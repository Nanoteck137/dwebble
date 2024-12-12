<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

<p>Change Quick Playlist</p>

<form method="post">
  <input name="playlistId" value="" type="hidden" />
  <Button type="submit" variant="outline">Unset Quick Playlist</Button>
</form>

<div class="h-4"></div>

<div class="flex flex-col gap-2">
  {#each data.playlists as playlist}
    <form method="post">
      <input name="playlistId" value={playlist.id} type="hidden" />
      <Button type="submit" variant="ghost">{playlist.name}</Button>
    </form>
  {/each}
</div>
