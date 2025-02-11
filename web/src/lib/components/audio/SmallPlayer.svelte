<script lang="ts">
  import {
    ChevronUp,
    Pause,
    Play,
    SkipBack,
    SkipForward,
    Volume2,
    VolumeX,
  } from "lucide-svelte";
  import { formatTime } from "$lib/utils";
  import {
    ScrollArea,
    Slider,
    Sheet,
    buttonVariants,
  } from "@nanoteck137/nano-ui";
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import { fly } from "svelte/transition";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import type { MediaItem } from "$lib/api/types";
  import Image from "$lib/components/Image.svelte";

  const musicManager = getMusicManager();

  interface Props {
    loading: boolean;
    playing: boolean;

    currentTime: number;
    duration: number;

    volume: number;
    audioMuted: boolean;

    currentMediaItem: MediaItem | null;

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
    loading,
    playing,

    currentTime,
    duration,

    volume,
    audioMuted,

    currentMediaItem,
    queue,
    currentQueueIndex,

    onPlay,
    onPause,
    onNextTrack,
    onPrevTrack,
    onSeek,
    onVolumeChanged,
    onToggleMuted,
  }: Props = $props();

  let vol = $state([0]);

  $effect(() => {
    vol = [volume * 100];
  });
</script>

{#snippet queueSheet()}
  <Sheet.Root>
    <Sheet.Trigger class={buttonVariants({ variant: "outline" })}>
      Queue
    </Sheet.Trigger>
    <Sheet.Content side="bottom">
      <p class="pb-2">Queue</p>
      <ScrollArea class="h-96">
        <div class="flex flex-col gap-2">
          {#each queue as mediaItem, i}
            <div class="flex items-center gap-2">
              <div class="group relative">
                <Image
                  class="w-12 min-w-12"
                  src={mediaItem.coverArt.small}
                  alt="cover"
                />
                {#if i == currentQueueIndex}
                  <div
                    class="absolute bottom-0 left-0 right-0 top-0 flex items-center justify-center border bg-black/80"
                  >
                    <Play size="25" />
                  </div>
                {:else}
                  <button
                    class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center border bg-black/80 group-hover:flex`}
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
                <p class="line-clamp-1 text-sm" title={mediaItem.track.name}>
                  {mediaItem.track.name}
                </p>
                <p
                  class="line-clamp-1 text-xs"
                  title={mediaItem.artists[0].name}
                >
                  {mediaItem.artists[0].name}
                </p>
              </div>
            </div>
          {/each}
        </div>
      </ScrollArea>
    </Sheet.Content>
  </Sheet.Root>
{/snippet}

<div
  class="z-30 h-16 border-t bg-background text-foreground md:hidden"
  transition:fly={{ y: 200 }}
>
  <div class="flex items-center">
    {#if playing}
      <button class="p-4" onclick={() => onPause()}>
        <Pause size="24" />
      </button>
    {:else}
      <button class="p-4" onclick={() => onPlay()}>
        <Play size="24" />
      </button>
    {/if}

    <Sheet.Root>
      <Sheet.Trigger class="flex grow items-center">
        <Image
          class="w-12 min-w-12"
          src={currentMediaItem?.coverArt.small}
          alt="cover"
        />

        <div class="flex flex-col items-start justify-center px-2">
          <p class="line-clamp-1 text-sm">{currentMediaItem?.track.name}</p>
          <p class="line-clamp-1 text-xs">
            {currentMediaItem?.artists[0].name}
          </p>
        </div>

        <div class="flex-grow"></div>
        <div class="flex h-16 min-w-16 items-center justify-center">
          <ChevronUp size="30" />
        </div>
      </Sheet.Trigger>
      <Sheet.Content side="bottom">
        <div class="relative flex flex-col items-center justify-center gap-2">
          {@render queueSheet()}

          <Image
            class="w-64"
            src={currentMediaItem?.coverArt.medium}
            alt="Track Cover Art"
          />

          <div class="flex flex-col items-center">
            <p class="font-medium">{currentMediaItem?.track.name}</p>
            <p class="text-xs">{currentMediaItem?.artists[0].name}</p>
          </div>

          <div class="flex w-full flex-col gap-1 px-4 py-2">
            <SeekSlider
              value={currentTime / duration}
              onValue={(p) => {
                onSeek(p * duration);
              }}
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
              <Slider
                bind:value={vol}
                onValueChange={(e) => onVolumeChanged(e[0] / 100)}
              />
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
      </Sheet.Content>
    </Sheet.Root>
  </div>
</div>
