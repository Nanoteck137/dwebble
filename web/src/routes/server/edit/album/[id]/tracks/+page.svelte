<script lang="ts">
  import TrackEdit from "../TrackEdit.svelte";

  const { data } = $props();
</script>

<form class="flex flex-col px-4 py-1" method="post">
  <div class="flex flex-col gap-2">
    {#each data.tracks as track (track.id)}
      <TrackEdit {track} />
    {/each}
  </div>

  <div class="flex gap-2 py-4">
    <button>Save</button>
    <a href="/server/edit/album/{data.album.id}">Back</a>
  </div>
</form>
