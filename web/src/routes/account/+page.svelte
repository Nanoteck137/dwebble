<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  let { data } = $props();
</script>

<p>{data.user?.username}</p>

<!-- <div class="relative">
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
</div> -->

<Button href="/account/password">Change Password</Button>
