<script lang="ts">
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { musicManager } from "$lib/music-manager";
  import { cn, shuffle, trackToMusicTrack } from "$lib/utils.js";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import {
    DiscAlbum,
    EllipsisVertical,
    Play,
    Shuffle,
    Users,
  } from "lucide-svelte";

  const { data } = $props();
</script>

<div class="flex flex-col">
  <div class="flex gap-2">
    <Button
      size="sm"
      variant="ghost"
      onclick={() => {
        musicManager.clearQueue();
        for (const track of data.items) {
          musicManager.addTrackToQueue(trackToMusicTrack(track));
        }
      }}
    >
      <Play />
      Play
    </Button>

    <Button
      size="sm"
      variant="ghost"
      onclick={() => {
        let tracks = shuffle([...data.items]);

        musicManager.clearQueue();
        for (const track of tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track));
        }
      }}
    >
      <Shuffle />
      Shuffle Play
    </Button>
  </div>

  <div class="flex items-center justify-between">
    <p class="text-bold text-xl">Tracks</p>
    <p class="text-sm">{data.items.length} track(s)</p>
  </div>

  {#each data.items as track, i}
    <TrackListItem
      {track}
      onPlayClicked={() => {
        musicManager.clearQueue();
        for (const track of data.items) {
          musicManager.addTrackToQueue(trackToMusicTrack(track), false);
        }
        musicManager.setQueueIndex(i);
        musicManager.requestPlay();
      }}
    >
      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={cn(
            buttonVariants({ variant: "ghost", size: "icon-lg" }),
            "rounded-full",
          )}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="end">
          <DropdownMenu.Group>
            <DropdownMenu.Item class="p-0">
              <!-- svelte-ignore a11y_invalid_attribute -->
              <a
                class="flex h-full w-full items-center gap-2 px-3 py-2"
                href="#"
              >
                <DiscAlbum size="16" />
                Go to Album
              </a>
            </DropdownMenu.Item>
            <DropdownMenu.Item class="p-0">
              <!-- svelte-ignore a11y_invalid_attribute -->
              <a
                class="flex h-full w-full items-center gap-2 px-3 py-2"
                href="#"
              >
                <Users size="16" />
                Go to Artist
              </a>
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </TrackListItem>
  {/each}
</div>
