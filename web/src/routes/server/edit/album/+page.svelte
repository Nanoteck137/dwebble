<script lang="ts">
  const { data } = $props();
</script>

<div class="flex flex-col">
  {#each data.albums as album}
    <a href={`/server/edit/album/${album.id}`}>{album.name}</a>
  {/each}
</div>
