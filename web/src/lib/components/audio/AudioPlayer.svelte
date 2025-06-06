<script lang="ts">
  import { onMount } from "svelte";
  import LargePlayer from "$lib/components/audio/LargePlayer.svelte";
  import SmallPlayer from "$lib/components/audio/SmallPlayer.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import { browser } from "$app/environment";
  import type { MediaItem } from "$lib/api/types";
  import toast from "svelte-5-french-toast";

  const musicManager = getMusicManager();

  let showPlayer = $state(false);

  let loading = $state(false);
  let playing = $state(false);

  let currentTime = $state(0);
  let duration = $state(0);

  let volume = $state(getVolume());
  let muted = $state(false);

  let currentMediaItem = $state<MediaItem | null>(null);

  let audio: HTMLAudioElement;

  function getVolume(): number {
    if (!browser) return 1.0;

    const volume = localStorage.getItem("player-volume");
    if (volume) {
      return parseFloat(volume);
    }

    return 1.0;
  }

  function getMuted(): boolean {
    const muted = localStorage.getItem("player-muted");
    if (muted) {
      return muted === "true";
    }

    return false;
  }

  function updateTrack() {
    const mediaItem = musicManager.queue.getCurrentMediaItem();

    if (mediaItem) {
      if (currentMediaItem?.track.id === mediaItem.track.id) return;

      currentMediaItem = mediaItem;
      audio.src = mediaItem.mediaUrl;
    } else {
      currentMediaItem = null;
      audio.removeAttribute("src");
      audio.load();
    }
  }

  onMount(() => {
    audio = new Audio();

    audio.addEventListener("canplay", () => {
      // console.log("canplay");
      loading = false;
    });

    audio.addEventListener("loadstart", () => {
      // console.log("loadstart");
      loading = true;
    });

    audio.addEventListener("loadedmetadata", () => {
      // console.log("loadedmetadata");
      currentTime = audio.currentTime;
      duration = audio.duration;
    });

    audio.addEventListener("progress", () => {
      // console.log("progress");
    });

    audio.addEventListener("timeupdate", () => {
      currentTime = audio.currentTime;
    });

    audio.addEventListener("loadeddata", () => {
      // console.log("loadeddata");
    });

    audio.addEventListener("playing", () => {
      // console.log("playing");
      playing = true;
    });

    audio.addEventListener("pause", () => {
      // console.log("pause");
      playing = false;
    });

    audio.addEventListener("load", () => {
      // console.log("load");
    });

    audio.addEventListener("ended", () => {
      musicManager.nextTrack();
      musicManager.requestPlay();
    });

    musicManager.emitter.on("requestPlay", () => {
      audio.play();
    });

    musicManager.emitter.on("requestPause", () => {
      audio.pause();
    });

    musicManager.emitter.on("requestPlayPause", () => {
      if (playing) {
        audio.pause();
      } else {
        audio.play();
      }
    });
  });

  onMount(() => {
    volume = getVolume();
    muted = getMuted();

    if (muted) {
      audio.volume = 0.0;
    } else {
      audio.volume = volume;
    }
  });

  let queue: MediaItem[] = $state([]);
  let currentQueueIndex = $state(0);

  onMount(() => {
    let unsub = musicManager.emitter.on("onQueueUpdated", () => {
      showPlayer = !musicManager.queue.isQueueEmpty();
      queue = musicManager.queue.items;
      currentQueueIndex = musicManager.queue.index;

      updateTrack();
    });

    return () => {
      unsub();
    };
  });

  $effect(() => {
    if (showPlayer) {
      document.body.setAttribute("data-player", "true");
    } else {
      document.body.setAttribute("data-player", "false");
    }
  });
</script>

<!-- TODO(patrik): Fix this because input fields need space -->
<!-- <svelte:window
  onkeypress={(e) => {
    if (e.key === " ") {
      e.preventDefault();

      musicManager.requestPlayPause();
    }
  }}
/> -->

{#if showPlayer}
  <LargePlayer
    {playing}
    {loading}
    mediaItem={currentMediaItem}
    {currentTime}
    {duration}
    {volume}
    {queue}
    {currentQueueIndex}
    audioMuted={muted}
    onPlay={() => {
      audio.play();
    }}
    onPause={() => {
      audio.pause();
    }}
    onNextTrack={() => {
      musicManager.nextTrack();
      audio.play();
    }}
    onPrevTrack={() => {
      musicManager.previousTrack();
      audio.play();
    }}
    onSeek={(e) => {
      audio.currentTime = e;
    }}
    onVolumeChanged={(e) => {
      if (!muted) {
        audio.volume = e;
      }

      volume = e;
      localStorage.setItem("player-volume", e.toString());
    }}
    onToggleMuted={() => {
      muted = !muted;
      localStorage.setItem("player-muted", muted ? "true" : "false");

      if (muted) {
        audio.volume = 0;
      } else {
        audio.volume = volume;
      }
    }}
  />

  <SmallPlayer
    {playing}
    {loading}
    {currentMediaItem}
    {currentTime}
    {duration}
    {volume}
    {queue}
    {currentQueueIndex}
    audioMuted={muted}
    onPlay={() => {
      audio.play();
    }}
    onPause={() => {
      audio.pause();
    }}
    onNextTrack={() => {
      musicManager.nextTrack();
      audio.play();
    }}
    onPrevTrack={() => {
      musicManager.previousTrack();
      audio.play();
    }}
    onSeek={(e) => {
      audio.currentTime = e;
    }}
    onVolumeChanged={(e) => {
      if (!muted) {
        audio.volume = e;
      }

      volume = e;
      localStorage.setItem("player-volume", e.toString());
    }}
    onToggleMuted={() => {
      muted = !muted;
      localStorage.setItem("player-muted", muted ? "true" : "false");

      if (muted) {
        audio.volume = 0;
      } else {
        audio.volume = volume;
      }
    }}
  />
{/if}
