<script lang="ts">
  import ProgressBar from "./ProgressBar.svelte";
  import {
    AudioHandler,
    currentPlayingSong,
    isPlaying,
    time,
    volume,
  } from "$lib/audio";

  import Play from "~icons/mdi/play";
  import Pause from "~icons/mdi/pause";

  $: progress = $time.duration > 0 ? $time.currentTime / $time.duration : 0;

  function nextSong() {
    AudioHandler.nextSong();
  }

  function prevSong() {
    // TODO(patrik): Move to audio.ts
    // if (player.currentTime >= 5) {
    //   player.currentTime = 0;
    //   return;
    // }

    AudioHandler.prevSong();
  }

  // TODO(patrik): Move utils file
  function formatTime(seconds: number) {
    const min = Math.floor(seconds / 60);
    const sec = Math.floor(seconds % 60);

    return `${min}:${sec.toString().padStart(2, "0")}`;
  }
</script>

<svelte:window
  on:keydown={(e) => {
    if (e.code === "Space") {
      e.preventDefault();
      e.stopPropagation();
      e.stopImmediatePropagation();

      AudioHandler.playPause();
    }
  }}
/>

<div class="flex h-full flex-col items-center">
  <ProgressBar
    {progress}
    on:progress={(e) => {
      // TODO(patrik): Implement
      // player.currentTime = e.detail * player.duration;
    }}
  />

  <div class="grid w-full grid-cols-5 items-center">
    <div class="flex items-center">
      <!-- {#if loading} -->
      <!--   <div role="status"> -->
      <!--     <svg -->
      <!--       aria-hidden="true" -->
      <!--       class="h-6 w-6 animate-spin fill-blue-600 text-gray-200 dark:text-gray-600" -->
      <!--       viewBox="0 0 100 101" -->
      <!--       fill="none" -->
      <!--       xmlns="http://www.w3.org/2000/svg" -->
      <!--     > -->
      <!--       <path -->
      <!--         d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z" -->
      <!--         fill="currentColor" -->
      <!--       /> -->
      <!--       <path -->
      <!--         d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z" -->
      <!--         fill="currentFill" -->
      <!--       /> -->
      <!--     </svg> -->
      <!--     <span class="sr-only">Loading...</span> -->
      <!--   </div> -->
      <!-- {/if} -->

      <button class="flex items-center justify-center" on:click={prevSong}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="h-6 w-6"
        >
          <path
            d="M9.195 18.44c1.25.713 2.805-.19 2.805-1.629v-2.34l6.945 3.968c1.25.714 2.805-.188 2.805-1.628V8.688c0-1.44-1.555-2.342-2.805-1.628L12 11.03v-2.34c0-1.44-1.555-2.343-2.805-1.629l-7.108 4.062c-1.26.72-1.26 2.536 0 3.256l7.108 4.061z"
          />
        </svg>
      </button>

      <button
        class="flex items-center justify-center"
        on:click={AudioHandler.playPause}
      >
        {#if $isPlaying == "playing"}
          <Pause />
        {:else if $isPlaying === "paused"}
          <Play />
        {/if}
      </button>

      <button class="flex items-center justify-center" on:click={nextSong}>
        <svg
          xmlns="http://www.w3.org/2000/svg"
          viewBox="0 0 24 24"
          fill="currentColor"
          class="h-6 w-6"
        >
          <path
            d="M5.055 7.06c-1.25-.714-2.805.189-2.805 1.628v8.123c0 1.44 1.555 2.342 2.805 1.628L12 14.471v2.34c0 1.44 1.555 2.342 2.805 1.628l7.108-4.061c1.26-.72 1.26-2.536 0-3.256L14.805 7.06C13.555 6.346 12 7.25 12 8.688v2.34L5.055 7.06z"
          />
        </svg>
      </button>

      <span
        >{formatTime($time.currentTime)} / {formatTime($time.duration)}</span
      >
    </div>

    {#if $currentPlayingSong}
      <p class="col-span-3">{$currentPlayingSong.name}</p>
    {:else}
      <p class="col-span-3">No song playing</p>
    {/if}

    <div class="w-full">
      <ProgressBar
        progress={$volume}
        on:progress={(e) => {
          volume.set(e.detail as number)
        }}
      />
    </div>
  </div>
</div>
