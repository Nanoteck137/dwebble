<script lang="ts">
  import { getContext } from "svelte";
  import type { Writable } from "svelte/store";

  const playing = getContext<Writable<string>>("playing");
</script>

<p>Test Page</p>

<button on:click={() => {
  playing.set("http://127.0.0.1:3000/songs/test.mp3");
}}>Update</button>
