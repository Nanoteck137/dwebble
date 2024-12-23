<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

<Button href="/taglists/new">New Taglist</Button>

<div class="flex flex-col gap-2">
  {#each data.taglists as taglist}
    <a href="/taglists/{taglist.id}">{taglist.name}</a>
  {/each}
</div>
