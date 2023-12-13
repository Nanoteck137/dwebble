<script lang="ts">
  import "../app.css";
  import AudioPlayer from "$lib/components/audioplayer/AudioPlayer.svelte";
  import Header from "$lib/components/Header.svelte";
  import { AudioHandler, playQueue } from "$lib/audio";
  import { onMount } from "svelte";

  onMount(() => {
    AudioHandler.init();
  });
</script>

<div class="flex h-screen">
  <Header />

  <main class="z-40 my-16 flex-grow">
    <slot />
    <!-- <div>
      {#each $playQueue.songs as queueItem, index}
        <p class={index === $playQueue.index ? "text-red-400" : ""}>
          {queueItem.name}
        </p>
      {/each}
    </div> -->
  </main>

  <footer class="fixed bottom-0 left-0 right-0 z-50 h-16 bg-white shadow-lg">
    <AudioPlayer />
  </footer>
</div>
