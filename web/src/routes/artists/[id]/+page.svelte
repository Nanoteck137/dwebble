<script lang="ts">
  import AlbumListItem from "$lib/components/AlbumListItem.svelte";
  import Image from "$lib/components/Image.svelte";
  import TrackListItem from "$lib/components/track-list/TrackListItem.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { isRoleAdmin } from "$lib/utils.js";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import {
    EllipsisVertical,
    ListPlus,
    Pencil,
    Play,
    Shuffle,
  } from "lucide-svelte";

  const { data } = $props();
  const musicManager = getMusicManager();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.artist.name.default}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="flex h-48">
  <Image class="w-48 min-w-48" src={data.artist.picture.medium} alt="cover" />
  <div class="w-2"></div>

  <div class="flex flex-col">
    <div class="flex flex-col">
      <p class="font-bold">
        {data.artist.name.default}
        {#if data.artist.name.other}
          - {data.artist.name.other}
        {/if}
      </p>

      {#if data.artist.tags.length > 0}
        <p class="text-xs text-muted-foreground">
          {data.artist.tags.join(", ")}
        </p>
      {/if}
    </div>

    <div class="flex-grow"></div>

    <div class="flex gap-2">
      <Button variant="outline" onclick={() => {}}>
        <Play />
        Play
      </Button>

      <Button variant="outline" onclick={() => {}}>
        <Shuffle />
        Shuffle
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={buttonVariants({ variant: "outline", size: "icon" })}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="start">
          <DropdownMenu.Group>
            <DropdownMenu.Item onSelect={async () => {}}>
              <ListPlus />
              Append to Queue
            </DropdownMenu.Item>
            {#if isRoleAdmin(data.user?.role || "")}
              <DropdownMenu.Link href="/artists/{data.artist.id}/edit">
                <Pencil />
                Edit Artist
              </DropdownMenu.Link>
            {/if}
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>

<div class="h-8"></div>

<a href="/artists/{data.artist.id}/tracks" class="text-xl font-bold">Tracks</a>

{#each data.tracks as track}
  <TrackListItem {track} />
  <Separator />
{/each}

<div class="h-2"></div>

<Button href="/artists/{data.artist.id}/tracks" variant="outline">
  Show More
</Button>

<div class="h-8"></div>

<a href="/artists/{data.artist.id}/albums" class="text-xl font-bold">Albums</a>

{#each data.albums as album}
  <AlbumListItem {album} link />
{/each}

<div class="h-2"></div>

<Button href="/artists/{data.artist.id}/albums" variant="outline">
  Show More
</Button>

{#if false}
  <Button href="/artists/{data.artist.id}/edit">Edit</Button>

  <Button
    onclick={async () => {
      await musicManager.clearQueue();
      await musicManager.addFromArtist(data.artist.id);
      musicManager.requestPlay();
    }}
  >
    Play
  </Button>

  <p>Artist: {data.artist.name.default}</p>

  <p>Num Albums: {data.albums.length}</p>
  <div class="flex flex-col">
    {#each data.albums as album}
      <a href="/albums/{album.id}">{album.name.default}</a>
    {/each}
  </div>

  <Separator />

  <p>Num Tracks: {data.trackPage.totalItems}</p>
  <div class="flex flex-col">
    {#each data.tracks as track}
      <a href="/albums/{track.albumId}">{track.name.default}</a>
    {/each}
  </div>
{/if}
