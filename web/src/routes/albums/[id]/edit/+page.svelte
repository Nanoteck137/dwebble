<script lang="ts">
  import { cn, formatTime } from "$lib/utils";
  import {
    Edit,
    EllipsisVertical,
    FolderPen,
    Import,
    Pencil,
    Play,
    Trash,
    Wallpaper,
  } from "lucide-svelte";
  import {
    Breadcrumb,
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { goto, invalidateAll } from "$app/navigation";
  import { modals } from "svelte-modals";
  import ConfirmModal from "$lib/components/modals/ConfirmModal.svelte";
  import { getApiClient, handleApiError } from "$lib";

  const { data } = $props();
  const apiClient = getApiClient();
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
          {data.album.name.default}
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
            <DropdownMenu.Item
              onSelect={() => {
                // musicManager.clearQueue();
                // for (const track of data.tracks) {
                //   musicManager.addTrackToQueue(
                //     trackToMusicTrack(track),
                //     false,
                //   );
                // }
                // musicManager.setQueueIndex(0);
                // musicManager.requestPlay();
              }}
            >
              <Play />
              Play
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${data.album.id}/edit/details`);
              }}
            >
              <Edit />
              Edit Album
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${data.album.id}/edit/import`);
              }}
            >
              <Import />
              Import Tracks
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${data.album.id}/edit/cover`);
              }}
            >
              <Wallpaper />
              Change Cover Art
            </DropdownMenu.Item>

            <DropdownMenu.Separator />

            <DropdownMenu.Item
              onSelect={() => {
                goto(`/albums/${data.album.id}/edit/delete`);
              }}
            >
              <Trash />
              Delete Album
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>

  <div class="flex flex-col py-2">
    <p class="font-bold">
      {data.album.name.default}
    </p>
    <p class="text-xs">
      Artist:
      <a class="hover:underline" href="/artists/{data.album.artistId}">
        {data.album.artistName.default}
      </a>
    </p>

    {#if data.album.name.other}
      <p class="text-xs">Other Name: {data.album.name.other}</p>
    {/if}

    {#if data.album.year}
      <p class="text-xs">Year: {data.album.year}</p>
    {/if}

    {#if data.album.tags.length > 0}
      <p class="text-xs">Tags: {data.album.tags.join(", ")}</p>
    {/if}

    {#if data.album.featuringArtists.length > 0}
      <p class="text-xs">Featuring Artists</p>
      {#each data.album.featuringArtists as artist}
        <p class="pl-2 text-xs">{artist.name.default}</p>
      {/each}
    {/if}
  </div>
</div>

<div class="py-4">
  <Separator />
</div>

<div class="flex flex-col">
  <div class="flex gap-2">
    <Button href="edit/tracks/common" class="w-full" variant="outline">
      <FolderPen />
      Set Common Values
    </Button>
    <Button href="edit/import" class="w-full" variant="outline">
      <Import />
      Import Tracks
    </Button>
  </div>

  {#each data.tracks as track (track.id)}
    <div class="flex items-center gap-2 py-2">
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-2">
          <p class="text-sm font-medium" title={track.name.default}>
            {#if track.number}
              <span>{track.number}.</span>
            {/if}
            {track.name.default}
          </p>
        </div>
        <div class="h-1"></div>
        <p class="text-xs" title={track.artistName.default}>
          Artist:
          <a class="hover:underline" href="/artists/{track.artistId}">
            {track.artistName.default}
          </a>
        </p>

        {#if track.name.other}
          <p class="text-xs">Other Name: {track.name.other}</p>
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

        {#if track.featuringArtists.length > 0}
          <p class="text-xs">Featuring Artists</p>
          {#each track.featuringArtists as artist}
            <p class="pl-2 text-xs">{artist.name.default}</p>
          {/each}
        {/if}
      </div>

      <div class="flex items-center gap-2">
        <!-- <QuickAddButton
          show={!!(data.user && data.user.quickPlaylist)}
          {track}
          {apiClient}
          isInQuickPlaylist={(trackId) => {
            if (!data.quickPlaylistIds) return false;
            return !!data.quickPlaylistIds.find((v) => v === trackId);
          }}
        /> -->

        <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          onclick={() => {
            // musicManager.clearQueue();
            // musicManager.addTrackToQueue(trackToMusicTrack(track), true);
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
                    // musicManager.clearQueue();
                    // musicManager.addTrackToQueue(
                    //   trackToMusicTrack(track),
                    //   true,
                    // );
                  }}
                >
                  <Play size="16" />
                  Play
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

              <DropdownMenu.Separator />

              <DropdownMenu.Item
                onSelect={async () => {
                  const confirmed = await modals.open(ConfirmModal, {
                    title: "Are you sure?",
                    description: "You are about to delete this track",
                    confirmDelete: true,
                  });

                  if (confirmed) {
                    const res = await apiClient.deleteTrack(track.id);
                    if (!res.success) {
                      handleApiError(res.error);
                      return;
                    }

                    await invalidateAll();
                  }
                }}
              >
                <Trash />
                Delete Track
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
    <Separator />
  {/each}
</div>
