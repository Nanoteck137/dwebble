<script lang="ts">
  import Slider from "$lib/components/SeekSlider.svelte";
  import { Button, ScrollArea, Separator, Sheet } from "@nanoteck137/nano-ui";
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
  import { getMusicManager } from "$lib/music-manager.svelte";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { MediaItem } from "$lib/api/types";

  const musicManager = getMusicManager();

  interface Props {
    loading: boolean;
    playing: boolean;

    currentTime: number;
    duration: number;

    volume: number;
    audioMuted: boolean;

    mediaItem: MediaItem | null;

    queue: MediaItem[];
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
    mediaItem,
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
        <div class="mr-3 flex flex-col gap-2">
          {#each queue as mediaItem, i}
            <div
              class={`flex items-center gap-2 rounded p-1 ${currentQueueIndex === i ? "bg-accent text-accent-foreground" : ""}`}
            >
              <div class="group relative">
                <img
                  class="inline-flex aspect-square w-12 min-w-12 items-center justify-center rounded border object-cover text-xs"
                  src={mediaItem?.coverArt.small}
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
                    onclick={async () => {
                      await musicManager.setQueueIndex(i);
                      musicManager.requestPlay();
                    }}
                  >
                    <Play class="text-muted-foreground" size="25" />
                  </button>
                {/if}
              </div>
              <div class="flex flex-col">
                <p class="line-clamp-1 text-sm" title={mediaItem?.track.name}>
                  {mediaItem?.track.name}
                </p>
                <p
                  class="line-clamp-1 text-xs font-light"
                  title={mediaItem?.artists[0].name}
                >
                  {mediaItem?.artists[0].name}
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
    <div class="flex items-center gap-2">
      <div class="flex items-center gap-2">
        <button
          onclick={() => {
            onPrevTrack();
          }}
        >
          <SkipBack size={32} />
        </button>

        {#if loading}
          <Spinner class="h-8 w-8" />
        {:else if playing}
          <button onclick={onPause}>
            <Pause size={32} />
          </button>
        {:else}
          <button onclick={onPlay}>
            <Play size={32} />
          </button>
        {/if}

        <button
          onclick={() => {
            onNextTrack();
          }}
        >
          <SkipForward size={32} />
        </button>
      </div>

      <p class="hidden min-w-20 text-xs font-medium lg:block">
        {formatTime(currentTime)} /{" "}
        {formatTime(Number.isNaN(duration) ? 0 : duration)}
      </p>

      <div class="flex items-center justify-center gap-2 align-middle">
        <img
          class="inline-flex aspect-square w-12 min-w-12 items-center justify-center rounded border object-cover text-xs"
          src={mediaItem?.coverArt.small}
          alt="cover"
        />
        <div class="flex flex-col">
          <p
            class="line-clamp-1 text-ellipsis text-sm"
            title={mediaItem?.track.name}
          >
            {mediaItem?.track.name}
          </p>

          <p class="line-clamp-1 min-w-80 text-ellipsis text-xs">
            {mediaItem?.artists[0].name}
          </p>
        </div>
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
      </div>
    </div>
  </div>
</div>
