<script lang="ts">
  import { enhance } from "$app/forms";

  let showQuickPlaylistDropDown = false;

  export let data;
</script>

<p>Account Page (W.I.P)</p>

<div class="relative">
  <button
    onclick={() => {
      showQuickPlaylistDropDown = true;
    }}>Set Quick Playlist</button
  >

  {#if showQuickPlaylistDropDown}
    <button
      class="fixed inset-0 z-50"
      onclick={() => {
        showQuickPlaylistDropDown = false;
      }}
    ></button>
  {/if}
  <div
    class={`absolute z-50 flex max-h-56 w-48 flex-col overflow-scroll bg-white ${showQuickPlaylistDropDown ? "block" : "hidden"}`}
  >
    {#each data.playlists as playlist}
      <form action="?/setQuickPlaylist" method="post" use:enhance>
        <input type="hidden" name="playlistId" value={playlist.id} />
        <button>{playlist.name}</button>
      </form>
    {/each}
  </div>
</div>
