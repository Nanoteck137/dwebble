<script lang="ts">
  import Separator from "$lib/components/ui/separator/separator.svelte";
  import TrackEdit from "../TrackEdit.svelte";

  const { data } = $props();
</script>

<form class="flex flex-col px-4 py-1" method="post">
  <div class="flex flex-col gap-2">
    {#each data.tracks as track (track.id)}
      <TrackEdit {track} />
      <Separator />
    {/each}
  </div>

  <div class="flex gap-2 py-4">
    <button>Save</button>
    <a href="/albums/{data.album.id}/edit">Back</a>
  </div>
</form>
