<script lang="ts">
  import { run } from 'svelte/legacy';

  import { browser } from "$app/environment";
  import Slider from "$lib/components/Slider.svelte";
  import { musicManager, type MusicTrack } from "$lib/music-manager";
  import { formatTime } from "$lib/utils";
  import {
    ListMusic,
    ListX,
    Logs,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
  } from "lucide-svelte";
  import { onMount } from "svelte";



  interface Props {
    showPlayer: boolean;
    loading: boolean;
    playing: boolean;
    currentTime: number;
    duration: number;
    volume: number;
    audioMuted: boolean;
    trackName: string;
    artistName: string;
    coverArt: string;
    onPlay: () => void;
    onPause: () => void;
    onNextTrack: () => void;
    onPrevTrack: () => void;
    onSeek: (e: number) => void;
    onVolumeChanged: (e: number) => void;
    onToggleMuted: () => void;
  }

  let {
    showPlayer,
    loading,
    playing,
    currentTime,
    duration,
    volume,
    audioMuted,
    trackName,
    artistName,
    coverArt,
    onPlay,
    onPause,
    onNextTrack,
    onPrevTrack,
    onSeek,
    onVolumeChanged,
    onToggleMuted
  }: Props = $props();

  let open = $state(false);
  let tracks: MusicTrack[] = $state([]);
  let currentTrack = $state(0);

  onMount(() => {
    let unsub = musicManager.emitter.on("onQueueUpdated", () => {
      tracks = musicManager.queue.items;
      currentTrack = musicManager.queue.index;
    });

    return () => {
      unsub();
    };
  });

  run(() => {
    if (open) {
      if (browser) document.body.style.overflow = "hidden";
    } else {
      if (browser) document.body.style.overflow = "";
    }
  });
</script>

<div
  class={`fixed bottom-0 left-0 right-0 z-30 hidden h-20 bg-purple-400 transition-transform duration-500 md:block ${showPlayer ? "translate-y-0" : "translate-y-[110%]"}`}
>
  <div class="absolute -top-2 left-0 right-0">
    <Slider
      value={currentTime / duration}
      onValue={(p) => {
        onSeek(p * duration);
      }}
    />
  </div>

  <div class="grid h-full grid-cols-footer">
    <div class="flex items-center bg-cyan-600">
      <div class="flex items-center">
        <button
          onclick={() => {
            onPrevTrack();
          }}
        >
          <SkipBack size="30" />
        </button>

        {#if loading}
          <p>Loading...</p>
        {:else if playing}
          <button onclick={onPause}>
            <Pause size={38} />
          </button>
        {:else}
          <button onclick={onPlay}>
            <Play size={38} />
          </button>
        {/if}

        <button
          onclick={() => {
            onNextTrack();
          }}
        >
          <SkipForward size="30" />
        </button>
      </div>

      <p class="hidden min-w-20 text-xs font-medium lg:block">
        {formatTime(currentTime)} /{" "}
        {formatTime(Number.isNaN(duration) ? 0 : duration)}
      </p>
    </div>

    <div
      class="flex items-center justify-center gap-2 bg-amber-600 align-middle"
    >
      <img
        class="aspect-square h-12 rounded object-cover"
        src={coverArt}
        alt="Cover Art"
      />
      <div class="flex flex-col">
        <p class="line-clamp-1 text-ellipsis" title={trackName}>
          {trackName}
        </p>

        <p class="line-clamp-1 min-w-80 text-ellipsis text-sm text-gray-800">
          {artistName}
        </p>
      </div>
    </div>
    <div class="flex items-center justify-evenly bg-violet-500">
      <div class="flex w-full items-center gap-4 p-4">
        <Slider
          value={volume}
          onValue={(p) => {
            onVolumeChanged(p);
          }}
        />
        <button
          onclick={() => {
            onToggleMuted();
          }}
        >
          {#if audioMuted}
            <VolumeX size="25" />
          {:else}
            <Volume2 size="25" />
          {/if}
        </button>
        <button
          onclick={() => {
            open = true;
          }}
        >
          <Logs size="24" />
        </button>
      </div>
    </div>
  </div>
</div>

{#if open}
  <button
    class="fixed inset-0 z-50 bg-[--modal-overlay-bg]"
    onclick={() => {
      open = false;
    }}
  ></button>
{/if}

<aside
  class={`fixed bottom-0 right-0 top-0 z-50 flex w-96 flex-col gap-2 overflow-y-scroll bg-[--bg-color] p-4 text-[--fg-color] transition-transform duration-300 ${open ? "translate-x-0" : "translate-x-[100%]"}`}
>
  <div class="flex items-center gap-2">
    <ListMusic size="32" />
    <p class="text-3xl font-semibold">Queue</p>
    <div class="flex-grow"></div>

    <button
      title="Clear Queue"
      onclick={() => {
        musicManager.clearQueue();
        open = false;
      }}
    >
      <ListX size="32" />
    </button>
  </div>
  <div class="flex flex-col gap-2">
    {#each tracks as track, i}
      <div class="flex gap-2">
        <div class="group relative">
          <img
            class="aspect-square w-12 rounded object-cover"
            src={track.coverArt}
            alt={`${track.name} Cover Art`}
          />
          {#if i == currentTrack}
            <div
              class="absolute bottom-0 left-0 right-0 top-0 flex items-center justify-center bg-[--overlay-bg]"
            >
              <Play size="25" />
            </div>
          {:else}
            <button
              class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center bg-[--overlay-bg] group-hover:flex`}
              onclick={() => {
                musicManager.setQueueIndex(i);
                musicManager.requestPlay();
              }}
            >
              <Play size="25" />
            </button>
          {/if}
        </div>
        <div class="flex flex-col">
          <p>{track.name}</p>
          <p>{track.artistName}</p>
        </div>
      </div>
    {/each}
  </div>
</aside>
