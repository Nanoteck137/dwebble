<script lang="ts">
  import { writable } from "svelte/store";
  import "../app.css";
  import { setContext } from "svelte";
  import AudioPlayer from "$lib/components/audioplayer/AudioPlayer.svelte";
  import Header from "$lib/components/Header.svelte";

  const playing = writable("http://10.28.28.6:3000/songs/song.m4a");
  setContext("playing", playing);

  let queue = [
    { name: "Song #1", url: "http://10.28.28.6:3000/songs/song.m4a" },
    { name: "Song #2", url: "http://10.28.28.6:3000/songs/test.mp3" },
  ];
</script>

<div class="flex h-screen">
  <Header />

  <main class="z-40 my-16 flex-grow">
    <slot />
  </main>

  <footer class="fixed bottom-0 left-0 right-0 z-50 h-16 bg-white shadow-lg">
    <AudioPlayer {queue} />
  </footer>
</div>
