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
      <a class="text-xs hover:underline" href="/artists/{data.album.artistId}">
        {data.album.artistName.default}
      </a>
    </div>

    <div class="flex-grow"></div>

    <div>
      <Button variant="outline">
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
    <div class="group flex items-center gap-2 py-1 pr-4">
      <p class="min-w-10 text-right text-sm font-medium group-hover:hidden">
        {#if track.number}
          <span>{track.number}.</span>
        {/if}
      </p>
      <button
        class="hidden h-10 w-10 items-center justify-end group-hover:flex"
        onclick={() => {
          musicManager.clearQueue();

          data.tracks.forEach((t) =>
            musicManager.addTrackToQueue(trackToMusicTrack(t), false),
          );

          musicManager.setQueueIndex(i);

          musicManager.requestPlay();
        }}
      >
        <Play size="20" />
      </button>
      <div class="flex flex-grow flex-col py-1">
        <p class="line-clamp-1 text-sm font-medium" title={track.name.default}>
          {track.name.default}
        </p>
        <a
          class="line-clamp-1 w-fit text-xs hover:underline"
          title={track.artistName.default}
          href={`/artists/${track.artistId}`}
        >
          {track.artistName.default}
        </a>
      </div>
      <div class="flex items-center gap-2">
        <p class="text-xs">
          {formatTime(track.duration ?? 0)}
        </p>
        <!-- <button
              class="hidden rounded-full p-1 hover:bg-black/20 group-hover:block"
            >
            </button> -->

        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item>
                <button
                  class="flex w-full items-center gap-2"
                  onclick={() => {
                    musicManager.clearQueue();

                    data.tracks.forEach((t) =>
                      musicManager.addTrackToQueue(
                        trackToMusicTrack(t),
                        false,
                      ),
                    );

                    musicManager.setQueueIndex(i);
                    musicManager.requestPlay();
                  }}
                >
                  <Play size="16" />
                  Play
                </button>
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
    <Separator />
  {/each}
</div>
