<script lang="ts">
  import { AudioHandler } from "$lib/audio";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<div class="p-2">
  <p>{data.album.name}</p>
  <button
    on:click={() => {
      AudioHandler.setQueue(data.album.tracks);
      AudioHandler.play();
    }}>Play Album</button
  >

  <div class="flex flex-wrap">
    {#each data.album.tracks as track, index}
      <div class="flex gap-4">
        <p>{track.num}</p>
        <p>{track.name}</p>
        <p>{track.album_name}</p>
        <button
          on:click={() => {
            AudioHandler.setQueue(data.album.tracks, index);
            AudioHandler.play();
          }}
          >Play
        </button>
      </div>
    {/each}
  </div>
</div>
