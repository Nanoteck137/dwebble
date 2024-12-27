<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { createApiClient } from "$lib";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { musicManager } from "$lib/music-manager";
  import { cn, formatError, shuffle, trackToMusicTrack } from "$lib/utils.js";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import {
    DiscAlbum,
    EllipsisVertical,
    FileHeart,
    Play,
    Shuffle,
    Users,
  } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = createApiClient(data);
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

    {#if data.user && data.user.quickPlaylist !== data.id}
      <Button
        size="sm"
        variant="ghost"
        onclick={async () => {
          const res = await apiClient.updateUserSettings({
            quickPlaylist: data.id,
          });
          if (!res.success) {
            toast.error("Unknown error");
            throw formatError(res.error);
          }

          await invalidateAll();
        }}
      >
        <FileHeart />
        Set as Quick Playlist
      </Button>
    {/if}
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
      <QuickAddButton
        show={!!(data.user && data.user.quickPlaylist)}
        {track}
        {apiClient}
        isInQuickPlaylist={(trackId) => {
          if (!data.quickPlaylistIds) return false;
          return !!data.quickPlaylistIds.find((v) => v === trackId);
        }}
      />

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
            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${track.albumId}`);
              }}
            >
              <DiscAlbum size="16" />
              Go to Album
            </DropdownMenu.Item>
            <DropdownMenu.Item
              onSelect={() => {
                goto(`/artists/${track.artistId}`);
              }}
            >
              <Users size="16" />
              Go to Artist
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </TrackListItem>
  {/each}
</div>
