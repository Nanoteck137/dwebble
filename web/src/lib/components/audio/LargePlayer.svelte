<script lang="ts">
  import Slider from "$lib/components/SeekSlider.svelte";
  import { Button, ScrollArea, Separator, Sheet } from "@nanoteck137/nano-ui";
  import { musicManager, type MusicTrack } from "$lib/music-manager";
  import { formatTime } from "$lib/utils";
  import {
    ListX,
    Logs,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
  } from "lucide-svelte";

  interface Props {
    loading: boolean;
    playing: boolean;

    currentTime: number;
    duration: number;

    volume: number;
    audioMuted: boolean;

    trackName: string;
    artistName: string;
    coverArt: string;

    queue: MusicTrack[];
    currentQueueIndex: number;

    onPlay: () => void;
    onPause: () => void;
    onNextTrack: () => void;
    onPrevTrack: () => void;
    onSeek: (e: number) => void;
    onVolumeChanged: (e: number) => void;
    onToggleMuted: () => void;
  }

  let {
    queue,
    currentQueueIndex,
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
    onToggleMuted,
  }: Props = $props();
</script>

{#snippet queueSheet()}
  <Sheet.Root>
    <Sheet.Trigger>
      <Logs size="24" />
    </Sheet.Trigger>
    <Sheet.Content side="right">
      <div class="flex items-center gap-2 pb-2">
        <p>Queue</p>
        <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          onclick={() => {
            musicManager.clearQueue();
          }}
        >
          <ListX />
        </Button>
      </div>

      <ScrollArea class="h-full pb-6">
        <div class="flex flex-col gap-2">
          {#each queue as track, i}
            <div
              class={`flex items-center gap-2 rounded p-1 ${currentQueueIndex === i ? "bg-accent text-accent-foreground" : ""}`}
            >
              <div class="group relative">
                <img
                  class="inline-flex aspect-square w-12 min-w-12 items-center justify-center rounded border object-cover text-xs"
                  src={track.coverArt}
                  alt="cover"
                />
                {#if i == currentQueueIndex}
                  <div
                    class="absolute bottom-0 left-0 right-0 top-0 flex items-center justify-center rounded border bg-muted/80"
                  >
                    <Play class="text-muted-foreground" size="25" />
                  </div>
                {:else}
                  <button
                    class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded border bg-muted/80 group-hover:flex`}
                    onclick={() => {
                      musicManager.setQueueIndex(i);
                      musicManager.requestPlay();
                    }}
                  >
                    <Play class="text-muted-foreground" size="25" />
                  </button>
                {/if}
              </div>
              <div class="flex flex-col">
                <p class="line-clamp-1 text-sm" title={track.name}>
                  {track.name}
                </p>
                <p
                  class="line-clamp-1 text-xs font-light"
                  title={track.artistName}
                >
                  {track.artistName}
                </p>
              </div>
            </div>
            <Separator />
          {/each}
        </div>
      </ScrollArea>
    </Sheet.Content>
  </Sheet.Root>
{/snippet}

<div
  class="container z-30 hidden h-16 bg-background text-foreground transition-transform duration-500 md:block"
>
  <div class="absolute -top-1.5 left-0 right-0">
    <Slider
      value={currentTime / duration}
      onValue={(p) => {
        onSeek(p * duration);
      }}
    />
  </div>

  <div class="grid h-full grid-cols-footer">
    <div class="flex items-center">
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

    <div class="flex items-center justify-center gap-2 align-middle">
      <img
        class="inline-flex aspect-square w-12 min-w-12 items-center justify-center rounded border object-cover text-xs"
        src={coverArt}
        alt="cover"
      />
      <div class="flex flex-col">
        <p class="line-clamp-1 text-ellipsis text-sm" title={trackName}>
          {trackName}
        </p>

        <p class="line-clamp-1 min-w-80 text-ellipsis text-xs">
          {artistName}
        </p>
      </div>
    </div>
    <div class="flex items-center justify-evenly">
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

        {@render queueSheet()}
        <!-- <button onclick={() => {}}>
          
        </button> -->
      </div>
    </div>
  </div>
</div>
