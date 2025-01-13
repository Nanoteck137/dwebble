<script lang="ts">
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, ListPlus, Pencil, Play } from "lucide-svelte";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import { createApiClient } from "$lib";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { goto } from "$app/navigation";

  let { data } = $props();
  const apiClient = createApiClient(data);
  const musicManager = getMusicManager();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.album.name.default}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="flex gap-2">
  <img
    class="inline-flex aspect-square w-48 min-w-48 items-center justify-center rounded border object-cover text-xs"
    src={data.album.coverArt.medium}
    alt="cover"
  />

  <div class="flex flex-col py-2">
    <div class="flex flex-col">
      <p class="font-bold">
        {data.album.name.default}
        {#if data.album.year}
          ({data.album.year})
        {/if}
      </p>
      <ArtistList artists={data.album.allArtists} />
      {#if data.album.tags}
        <p class="text-xs">{data.album.tags.join(", ")}</p>
      {/if}
      <!-- <a class="text-xs hover:underline" href="/artists/{data.album.artistId}">
        {data.album.artistName.default}
      </a> -->
    </div>

    <div class="flex-grow"></div>

    <div>
      <Button
        variant="outline"
        onclick={async () => {
          await musicManager.clearQueue();
          await musicManager.addFromAlbum(data.album.id);
          musicManager.requestPlay();
        }}
      >
        <Play />
        Play
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={buttonVariants({ variant: "outline", size: "icon" })}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
          <DropdownMenu.Group>
            <DropdownMenu.Item
              onSelect={async () => {
                await musicManager.addFromAlbum(data.album.id);
                musicManager.requestPlay();
              }}
            >
              <ListPlus />
              Append to Queue
            </DropdownMenu.Item>
            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${data.album.id}/edit`);
              }}
            >
              <Pencil />
              Edit Album
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>

<div class="h-4"></div>

<TrackList
  {apiClient}
  isAlbumShowcase={true}
  totalTracks={0}
  tracks={data.tracks}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  isInQuickPlaylist={(trackId) => {
    if (!data.quickPlaylistIds) return false;
    return !!data.quickPlaylistIds.find((v) => v === trackId);
  }}
/>
