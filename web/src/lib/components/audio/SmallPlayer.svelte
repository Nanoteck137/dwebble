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
  import * as Drawer from "$lib/components/ui/drawer";
  import { formatTime } from "$lib/utils";
  import { Slider } from "$lib/components/ui/slider";
  import SeekSlider from "$lib/components/SeekSlider.svelte";

  // let open = $state(false);

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
    onToggleMuted,
  }: Props = $props();

  let vol = $state([0]);

  $effect(() => {
    vol = [volume * 100];
  });

  // $effect(() => {
  //   if (open) {
  //     if (browser) document.body.style.overflow = "hidden";
  //   } else {
  //     if (browser) document.body.style.overflow = "";
  //   }
  // });
</script>

<div
  class="fixed bottom-0 left-0 right-0 z-30 h-16 border-t bg-background text-foreground md:hidden"
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

    <Drawer.Root>
      <Drawer.Trigger class="flex grow items-center">
        <img
          class="aspect-square h-16 rounded object-cover p-1"
          src={coverArt}
          alt="Cover Art"
        />

        <div class="flex flex-col items-start justify-center px-2">
          <p class="line-clamp-1 text-sm">{trackName}</p>
          <p class="line-clamp-1 text-xs">{artistName}</p>
        </div>

        <div class="flex-grow"></div>
        <div class="flex h-16 min-w-16 items-center justify-center">
          <ChevronUp size="30" />
        </div>
      </Drawer.Trigger>
      <Drawer.Content>
        <!-- <Drawer.Header>
          <Drawer.Title>Are you sure absolutely sure?</Drawer.Title>
          <Drawer.Description>This action cannot be undone.</Drawer.Description
          >
        </Drawer.Header> -->

        <div class="relative flex flex-col items-center justify-center gap-2">
          <img
            class="aspect-square w-80 rounded object-cover"
            src={coverArt}
            alt="Track Cover Art"
          />
          <p class="text-lg font-medium">{trackName}</p>
          <p class="">{artistName}</p>

          <div class="flex w-full flex-col gap-1 px-4">
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

        <Drawer.Footer>
          <Drawer.Close>Close</Drawer.Close>
        </Drawer.Footer>
      </Drawer.Content>
    </Drawer.Root>
  </div>
</div>
