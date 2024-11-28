<script lang="ts">
  const { data } = $props();
</script>

<div class="flex flex-col gap-2">
  {#each data.taglists as taglist}
    <a href="/taglists/{taglist.id}">{taglist.name}</a>
  {/each}
</div>
