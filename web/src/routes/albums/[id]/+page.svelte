<script lang="ts">
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, ListPlus, Pencil, Play } from "lucide-svelte";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { goto } from "$app/navigation";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";

  let { data } = $props();
  const musicManager = getMusicManager();

  function getName() {
    const name = data.album.name.default;

    if (data.album.year) {
      return `${name} (${data.album.year})`;
    }

    return name;
  }
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

<TrackListHeader
  name={getName()}
  image={data.album.coverArt.medium}
  artists={data.album.allArtists}
  tags={data.album.tags}
>
  {#snippet more()}
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
      <DropdownMenu.Link href="/albums/{data.album.id}/edit">
        <Pencil />
        Edit Album
      </DropdownMenu.Link>
    </DropdownMenu.Group>
  {/snippet}
</TrackListHeader>

{#if 0}
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
          <DropdownMenu.Content></DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  </div>
{/if}

<div class="h-4"></div>

<TrackList
  isAlbumShowcase={true}
  totalTracks={0}
  tracks={data.tracks}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  isInQuickPlaylist={(trackId) => {
    if (!data.quickPlaylistIds) return false;
    return !!data.quickPlaylistIds.find((v) => v === trackId);
  }}
  onPlay={async () => {
    await musicManager.clearQueue();
    await musicManager.addFromAlbum(data.album.id);
    musicManager.requestPlay();
  }}
  onTrackPlay={async (trackId) => {
    await musicManager.clearQueue();
    await musicManager.addFromAlbum(data.album.id);
    await musicManager.setQueueIndex(
      musicManager.queue.items.findIndex((t) => t.track.id === trackId),
    );
    musicManager.requestPlay();
  }}
/>
