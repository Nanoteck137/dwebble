<script lang="ts">
  const { data } = $props();
</script>

<p>Edit Album</p>

<p>{data.album.name}</p>
<p>{data.album.artistName}</p>
<p>{data.album.coverArt}</p>

{#each data.tracks as track (track.id)}
  <p>{track.name}</p>
{/each}
