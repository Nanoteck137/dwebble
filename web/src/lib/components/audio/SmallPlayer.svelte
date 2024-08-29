<script lang="ts">
  import { browser } from "$app/environment";
  import Slider from "$lib/components/Slider.svelte";
  import { formatTime } from "$lib/utils";
  import {
    ChevronDown,
    ChevronUp,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
  } from "lucide-svelte";

  let open = false;

  export let showPlayer: boolean;

  export let loading: boolean;
  export let playing: boolean;
  export let currentTime: number;
  export let duration: number;
  export let volume: number;
  export let audioMuted: boolean;
  export let trackName: string;
  export let artistName: string;
  export let coverArt: string;

  export let onPlay: () => void;
  export let onPause: () => void;
  export let onNextTrack: () => void;
  export let onPrevTrack: () => void;
  export let onSeek: (e: number) => void;
  export let onVolumeChanged: (e: number) => void;
  export let onToggleMuted: () => void;

  $: {
    if (open) {
      if (browser) document.body.style.overflow = "hidden";
    } else {
      if (browser) document.body.style.overflow = "";
    }
  }
</script>

<div
  class={`fixed bottom-0 left-0 right-0 z-30 h-16 bg-[--bg-color-alt] text-[--fg-color] transition-transform duration-300 md:hidden ${open || !showPlayer ? "translate-y-[100%]" : "translate-y-0"}`}
>
  <div class="flex items-center">
    {#if playing}
      <button class="p-4" onclick={() => onPause()}>
        <Pause size="28" />
      </button>
    {:else}
      <button class="p-4" onclick={() => onPlay()}>
        <Play size="28" />
      </button>
    {/if}
    <button
      class="flex grow items-center"
      onclick={() => {
        open = true;
      }}
    >
      <img
        class="aspect-square h-16 rounded object-cover p-1"
        src={coverArt}
        alt="Cover Art"
      />
      <div class="flex flex-col items-start justify-center px-2">
        <p class="line-clamp-1">{trackName}</p>
        <p class="line-clamp-1">{artistName}</p>
      </div>
      <div class="flex-grow"></div>
      <div class="flex h-16 min-w-16 items-center justify-center">
        <ChevronUp size="30" />
      </div>
    </button>
  </div>
</div>

<div
  class={`fixed left-0 right-0 top-0 z-50 h-screen bg-[--bg-color] transition-transform duration-300 md:hidden ${open ? "" : "translate-y-[100%]"}`}
>
  <div class="relative flex flex-col items-center justify-center gap-2">
    <div class="flex h-12 w-full items-center">
      <div class="w-2"></div>
      <button
        class="rounded-full p-2 hover:bg-black/70"
        onclick={() => {
          open = false;
        }}
      >
        <ChevronDown size="30" />
      </button>
    </div>
    <img
      class="aspect-square w-80 rounded object-cover"
      src={coverArt}
      alt="Track Cover Art"
    />
    <p class="text-lg font-medium">{trackName}</p>
    <p class="">{artistName}</p>

    <div class="flex w-full flex-col gap-1 px-4">
      <Slider
        value={currentTime / duration}
        onValue={(e) => onSeek(e * duration)}
      />

      <div class="flex justify-between">
        <p class="text-sm">
          {formatTime(currentTime)}
        </p>

        <p class="text-sm">
          {formatTime(Number.isNaN(duration) ? 0 : duration)}
        </p>
      </div>
    </div>

    <div class="flex w-full items-center gap-4 px-4">
      <div class="flex gap-4">
        <button
          onclick={() => {
            onPrevTrack();
          }}
        >
          <SkipBack size="38" />
        </button>

        {#if loading}
          <p>Loading...</p>
        {:else if playing}
          <button onclick={onPause}>
            <Pause size={46} />
          </button>
        {:else}
          <button onclick={onPlay}>
            <Play size={46} />
          </button>
        {/if}

        <button
          onclick={() => {
            onNextTrack();
          }}
        >
          <SkipForward size="38" />
        </button>
      </div>

      <div class="flex-grow"></div>

      <div class="flex w-full max-w-56 items-center gap-4">
        <Slider value={volume} onValue={(e) => onVolumeChanged(e)} />
        <button
          onclick={() => {
            onToggleMuted();
          }}
        >
          {#if audioMuted}
            <VolumeX size="30" />
          {:else}
            <Volume2 size="30" />
          {/if}
        </button>
      </div>
    </div>
  </div>
</div>
