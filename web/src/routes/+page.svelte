<script lang="ts">
  import { flip } from "svelte/animate";
  import { fade, fly, slide } from "svelte/transition";

  let popups = $state<string[]>([]);
</script>

<p class="p-4 text-xl">Home Page</p>

<button
  onclick={() => {
    popups.push("Hello " + Date.now());
  }}>Test</button
>

{#each popups as popup (popup)}
  <div animate:flip={{ duration: 2000 }}>
    <p>{popup}</p>
  </div>
{/each}
