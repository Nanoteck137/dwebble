<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

{#each data.playlists as playlist}
  <p>{playlist.name}</p>
{/each}

