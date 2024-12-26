<script lang="ts">
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { musicManager } from "$lib/music-manager";
  import { cn, formatTime, trackToMusicTrack } from "$lib/utils";
  import { EllipsisVertical, Pencil, Play } from "lucide-svelte";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import TrackListItem from "$lib/components/TrackListItem.svelte";

  let { data } = $props();
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
      <p class="font-bold">{data.album.name.default}</p>
      <ArtistList artists={data.album.allArtists} />
      <!-- <a class="text-xs hover:underline" href="/artists/{data.album.artistId}">
        {data.album.artistName.default}
      </a> -->
    </div>

    <div class="flex-grow"></div>

    <div>
      <Button
        variant="outline"
        onclick={() => {
          musicManager.clearQueue();

          for (const track of data.tracks) {
            musicManager.addTrackToQueue(trackToMusicTrack(track), false);
          }

          musicManager.setQueueIndex(0);
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
            <DropdownMenu.Item>
              <a
                class="flex w-full items-center gap-2"
                href="/albums/{data.album.id}/edit"
              >
                <Pencil size="16" />
                Edit Album
              </a>
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>
</div>

<div class="h-4"></div>

<div class="flex flex-col gap-2">
  {#each data.tracks as track, i}
    <TrackListItem
      {track}
      showNumber={true}
      onPlayClicked={() => {
        musicManager.clearQueue();

        for (const track of data.tracks) {
          musicManager.addTrackToQueue(trackToMusicTrack(track), false);
        }

        musicManager.setQueueIndex(i);
        musicManager.requestPlay();
      }}
    />
    <Separator />
  {/each}
</div>
