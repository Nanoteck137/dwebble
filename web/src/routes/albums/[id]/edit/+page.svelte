<script lang="ts">
  import { musicManager } from "$lib/music-manager";
  import { cn, formatTime, trackToMusicTrack } from "$lib/utils";
  import { Edit, EllipsisVertical, Pencil, Play, Trash } from "lucide-svelte";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { goto } from "$app/navigation";

  const { data } = $props();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums/{data.album.id}">
          {data.album.name}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Edit</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="flex gap-2">
  <div class="relative aspect-square w-48 min-w-48">
    <img
      class="inline-flex h-full w-full items-center justify-center rounded border object-cover text-xs"
      src={data.album.coverArt.medium}
      alt="cover"
    />
    <div class="absolute inset-0 flex justify-end rounded p-1">
      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={cn(
            buttonVariants({ variant: "ghost", size: "icon" }),
            "rounded-full",
          )}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="center">
          <DropdownMenu.Group>
            <DropdownMenu.Item onSelect={() => {}}>
              <Play size="16" />
              Play
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${data.album.id}/edit/delete`);
              }}
            >
              <Trash />
              Remove
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>

  <div class="flex flex-col py-2">
    <p class="font-bold">
      {data.album.name}
    </p>
    <p class="text-xs">
      Artist: <a class="hover:underline" href="/artists/{data.album.artistId}">
        {data.album.artistName}
      </a>
    </p>

    {#if data.album.otherName}
      <p class="text-xs">Other Name: {data.album.otherName}</p>
    {/if}

    {#if data.album.year}
      <p class="text-xs">Year: {data.album.year}</p>
    {/if}

    <div class="flex-grow"></div>

    <div class="flex gap-2">
      <a class="text-sm hover:underline" href="edit/details">
        Edit Album Details
      </a>
      <a class="text-sm hover:underline" href="edit/import">Import Tracks</a>
    </div>
  </div>
</div>

<div class="flex flex-col">
  <Button href="edit/tracks/common" class="w-full" variant="outline">
    Set Common Values
  </Button>

  {#each data.tracks as track (track.id)}
    <div class="flex items-center gap-2 py-2">
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-2">
          <p class="text-sm font-medium" title={track.name}>
            {#if track.number}
              <span>{track.number}.</span>
            {/if}
            {track.name}
          </p>
        </div>
        <div class="h-1"></div>
        <p class="text-xs" title={track.artistName}>
          Artist: <a class="hover:underline" href="/artists/{track.artistId}"
            >{track.artistName}</a
          >
        </p>

        {#if track.otherName}
          <p class="text-xs">Other Name: {track.otherName}</p>
        {/if}

        {#if track.year}
          <p class="text-xs">Year: {track.year}</p>
        {/if}

        {#if track.tags.length > 0}
          <p class="text-xs">Tags: {track.tags.join(", ")}</p>
        {/if}

        {#if track.duration}
          <p class="text-xs">Duration: {formatTime(track.duration ?? 0)}</p>
        {/if}
      </div>

      <div class="flex items-center gap-2">
        <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          onclick={() => {
            musicManager.clearQueue();
            musicManager.addTrackToQueue(trackToMusicTrack(track), true);
          }}
        >
          <Play />
        </Button>

        <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          href="edit/tracks/{track.id}"
        >
          <Pencil />
        </Button>

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
                    musicManager.addTrackToQueue(
                      trackToMusicTrack(track),
                      true,
                    );
                  }}
                >
                  <Play size="16" />
                  Play
                </button>
              </DropdownMenu.Item>

              <DropdownMenu.Item>
                <button
                  class="flex w-full items-center gap-2"
                  onclick={() => {}}
                >
                  <Trash size="16" />
                  Remove
                </button>
              </DropdownMenu.Item>

              <DropdownMenu.Item>
                <a
                  class="flex w-full items-center gap-2"
                  href="edit/tracks/{track.id}"
                >
                  <Edit size="16" />
                  Edit
                </a>
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
    <Separator />
  {/each}
</div>
