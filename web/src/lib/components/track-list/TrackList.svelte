<script lang="ts">
  import {
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import TrackListItem from "./TrackListItem.svelte";
  import { DiscAlbum, EllipsisVertical, ListPlus, Users } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import type { Playlist, Track } from "$lib/api/types";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import { goto, invalidateAll } from "$app/navigation";

  type Props = {
    totalTracks: number;
    tracks: Track[];

    isAlbumShowcase?: boolean;

    userPlaylists?: Playlist[] | null;
    quickPlaylist?: string | null;

    onPlay: (shuffle: boolean) => void;
    onTrackPlay: (trackId: string) => void;
  };

  const {
    isAlbumShowcase,
    totalTracks,
    tracks,
    userPlaylists,
    quickPlaylist,
    onTrackPlay,
  }: Props = $props();
</script>

<div class="flex flex-col">
  <div class="flex items-center justify-between">
    <p class="text-bold text-xl">Tracks</p>
    <p class="text-sm">{totalTracks} track(s)</p>
  </div>

  {#each tracks as track}
    <TrackListItem
      showNumber={isAlbumShowcase}
      {track}
      onPlayClicked={() => {
        onTrackPlay(track.id);
      }}
    >
      <QuickAddButton show={!!quickPlaylist} trackId={track.id} />

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
                goto(`/artists/${track.artists[0].id}`);
              }}
            >
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

                // await openAddToPlaylist({
                //   playlists: userPlaylists,
                //   track,
                // });

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
