<script lang="ts">
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import TrackListItem from "./TrackListItem.svelte";
  import { musicManager } from "$lib/music-manager";
  import {
    DiscAlbum,
    EllipsisVertical,
    ListPlus,
    Play,
    Shuffle,
    Users,
  } from "lucide-svelte";
  import { cn, shuffle, trackToMusicTrack } from "$lib/utils";
  import type { Playlist, Track } from "$lib/api/types";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import type { ApiClient } from "$lib/api/client";
  import { openAddToPlaylist } from "$lib";
  import { goto, invalidateAll } from "$app/navigation";

  type Props = {
    apiClient: ApiClient;
    isAlbumShowcase?: boolean;
    totalTracks: number;
    tracks: Track[];
    userPlaylists?: Playlist[] | null;
    quickPlaylist?: string | null;
    isInQuickPlaylist: (trackId: string) => boolean;
  };

  const {
    apiClient,
    isAlbumShowcase,
    totalTracks,
    tracks,
    userPlaylists,
    quickPlaylist,
    isInQuickPlaylist,
  }: Props = $props();
</script>

<div class="flex flex-col">
  {#if !isAlbumShowcase}
    <div class="flex gap-2">
      <Button
        size="sm"
        variant="ghost"
        onclick={() => {
          musicManager.clearQueue();
          for (const track of tracks) {
            musicManager.addTrackToQueue(trackToMusicTrack(track));
          }
        }}
      >
        <Play />
        Play
      </Button>

      <!-- <form action="?/newPlaylist" method="post">
      <input name="filter" value={data.filter} type="hidden" />
      <input name="sort" value={data.sort} type="hidden" />
      <Button type="submit" size="sm" variant="ghost">
        <ListPlus />
        Create Playlist
      </Button>
    </form> -->

      <Button
        size="sm"
        variant="ghost"
        onclick={() => {
          const newTracks = shuffle([...tracks]);

          musicManager.clearQueue();
          for (const track of newTracks) {
            musicManager.addTrackToQueue(trackToMusicTrack(track));
          }
        }}
      >
        <Shuffle />
        Shuffle Play
      </Button>
    </div>
  {/if}

  <div class="flex items-center justify-between">
    <p class="text-bold text-xl">Tracks</p>
    <p class="text-sm">{totalTracks} track(s)</p>
  </div>

  {#each tracks as track, i}
    <TrackListItem
      showNumber={isAlbumShowcase}
      {track}
      onPlayClicked={() => {
        musicManager.clearQueue();
        for (const track of tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track), false);
        }
        musicManager.setQueueIndex(i);
        musicManager.requestPlay();
      }}
    >
      <QuickAddButton
        show={!!quickPlaylist}
        {track}
        {apiClient}
        {isInQuickPlaylist}
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
            <DropdownMenu.Item>
              <Users />
              Go to Artist
            </DropdownMenu.Item>
            {#if !isAlbumShowcase}
              <DropdownMenu.Item
                onSelect={() => {
                  goto(`/albums/${track.albumId}`);
                }}
              >
                <DiscAlbum />
                Go to Album
              </DropdownMenu.Item>
            {/if}
            <DropdownMenu.Item
              onSelect={async () => {
                if (!userPlaylists) return;

                await openAddToPlaylist({
                  apiClient,
                  playlists: userPlaylists,
                  track,
                });

                await invalidateAll();
              }}
            >
              <ListPlus />
              Save to Playlist
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </TrackListItem>

    <Separator />
  {/each}
</div>
