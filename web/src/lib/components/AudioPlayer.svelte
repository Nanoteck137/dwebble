<script lang="ts">
  import { onMount } from "svelte";
  import { musicManager } from "$lib/music-manager";
  import LargePlayer from "$lib/components/audio/LargePlayer.svelte";
  import SmallPlayer from "$lib/components/audio/SmallPlayer.svelte";

  let showPlayer = $state(false);

  let loading = $state(false);
  let playing = $state(false);
  let currentTime = $state(0);
  let duration = $state(0);
  let volume = $state(0);

  let trackName = $state("");
  let artistName = $state("");
  let coverArt = $state("");

  let muted = $state(false);

  let audio: HTMLAudioElement;

  function getVolume(): number {
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

  onMount(() => {
    audio = new Audio();

    audio.addEventListener("canplay", () => {
      console.log("canplay");
      loading = false;
    });

    audio.addEventListener("loadstart", () => {
      console.log("loadstart");
      loading = true;
    });

    audio.addEventListener("loadedmetadata", () => {
      console.log("loadedmetadata");
      currentTime = audio.currentTime;
      duration = audio.duration;
    });

    audio.addEventListener("progress", () => {
      console.log("progress");
    });

    audio.addEventListener("timeupdate", () => {
      currentTime = audio.currentTime;
    });

    audio.addEventListener("loadeddata", () => {
      console.log("loadeddata");
    });

    audio.addEventListener("playing", () => {
      console.log("playing");
      playing = true;
    });

    audio.addEventListener("pause", () => {
      console.log("pause");
      playing = false;
    });

    audio.addEventListener("load", () => {
      console.log("load");
    });

    audio.addEventListener("ended", () => {
      musicManager.nextTrack();
      musicManager.requestPlay();
    });

    musicManager.emitter.on("onTrackChanged", () => {
      const track = musicManager.getCurrentTrack();
      console.log(track);
      if (track) {
        trackName = track.name;
        artistName = track.artistName;
        coverArt = track.coverArt;

        audio.src = track.source;
      } else {
        trackName = "";
        artistName = "";
        coverArt = "";

        audio.removeAttribute("src");
        audio.load();
      }
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

  onMount(() => {
    let unsub = musicManager.emitter.on("onQueueUpdated", () => {
      showPlayer = !musicManager.isQueueEmpty();
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

<svelte:window
  onkeypress={(e) => {
    if (e.key === " ") {
      e.preventDefault();
      e.stopPropagation();
      e.stopImmediatePropagation();

      musicManager.requestPlayPause();
    }
  }}
/>

<LargePlayer
  {showPlayer}
  {playing}
  {loading}
  {trackName}
  {artistName}
  {coverArt}
  {currentTime}
  {duration}
  {volume}
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
    musicManager.prevTrack();
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
  {showPlayer}
  {playing}
  {loading}
  {trackName}
  {artistName}
  {coverArt}
  {currentTime}
  {duration}
  {volume}
  audioMuted={false}
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
    musicManager.prevTrack();
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

<!-- <div
  class="fixed bottom-0 left-0 right-0 flex h-14 items-center gap-4 bg-red-300 px-2"
>
  <div class="flex gap-2">
    <button
      onclick={() => {
        musicManager.prevTrack();
        musicManager.requestPlay();
      }}>Prev</button
    >

    {#if loading}
      <p>Loading...</p>
    {:else if playing}
      <button
        onclick={() => {
          musicManager.requestPause();
        }}>Pause</button
      >
    {:else}
      <button
        onclick={() => {
          musicManager.requestPlay();
        }}>Play</button
      >
    {/if}

    <button
      onclick={() => {
        musicManager.nextTrack();
        musicManager.requestPlay();
      }}>Next</button
    >
  </div>

  <div class="flex flex-col">
    <p class="text-sm">{trackName}</p>
    <p class="text-xs">{artistName}</p>
  </div>

  <p>{formatTime(currentTime)} / {formatTime(duration)}</p>
</div> -->
