<script lang="ts">
  import { cn, convertValue, formatTime } from "$lib/utils";
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
  import { getApiClient, handleApiError } from "$lib";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import EditSingleModal from "./EditSingleModal.svelte";
  import Image from "$lib/components/Image.svelte";
  import EditAlbumDetailsModal from "./EditAlbumDetailsModal.svelte";
  import toast from "svelte-5-french-toast";
  import ChangeAlbumCoverModal from "./ChangeAlbumCoverModal.svelte";
  import SetCommonValuesModal from "./SetCommonValuesModal.svelte";
  import EditTrackDetailsModal from "./EditTrackDetailsModal.svelte";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import ConfirmModal from "$lib/components/new-modals/ConfirmModal.svelte";
  import ImportTracksModal from "./ImportTracksModal.svelte";
  import { UploadTrackBody } from "$lib/api/types";

  const { data } = $props();
  const apiClient = getApiClient();
  const musicManager = getMusicManager();

  let openEditAlbumDetails = $state(false);
  let openChangeAlbumCoverModal = $state(false);

  let openConfirmDeleteAlbum = $state(false);
  let openConfirmDeleteTrack = $state({ open: false, trackId: "" });
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
    <Image
      class="h-full w-full"
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
              onSelect={async () => {
                await musicManager.clearQueue();
                await musicManager.addFromAlbum(data.album.id);
                await musicManager.setQueueIndex(0);
                musicManager.requestPlay();
              }}
            >
              <Play />
              Play
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                openEditAlbumDetails = true;
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
                openChangeAlbumCoverModal = true;
              }}
            >
              <Wallpaper />
              Change Cover Art
            </DropdownMenu.Item>

            <DropdownMenu.Separator />

            <DropdownMenu.Item
              onSelect={() => {
                openConfirmDeleteAlbum = true;
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
  <div class="flex flex-col gap-2 md:flex-row">
    {#if data.tracks.length === 1}
      <EditSingleModal
        class={buttonVariants({ variant: "outline", class: "w-full" })}
        onResult={async (editData) => {
          const name = convertValue(editData.name);
          const otherName = convertValue(editData.otherName);
          const year = convertValue(editData.year);
          const tags = convertValue(editData.tags);
          const artist = convertValue(editData.artist);
          const featuringArtists = convertValue(editData.featuringArtists);

          const body = {
            name,
            otherName,
            year,
            tags: tags?.split(","),
            artistId: artist?.id,
            featuringArtists: featuringArtists?.map((a) => a.id),
          };

          {
            const res = await apiClient.editAlbum(data.album.id, body);
            if (!res.success) {
              handleApiError(res.error);
            }
          }

          {
            const res = await apiClient.editTrack(data.tracks[0].id, body);
            if (!res.success) {
              handleApiError(res.error);
            }
          }

          invalidateAll();
        }}
      >
        <FolderPen />
        Edit as Single
      </EditSingleModal>
    {/if}

    <SetCommonValuesModal
      class={buttonVariants({ variant: "outline", class: "w-full" })}
      onResult={async (resultData) => {
        const year = convertValue(resultData.year);
        const tags = convertValue(resultData.tags);
        const artist = convertValue(resultData.artist);
        const featuringArtists = convertValue(resultData.featuringArtists);

        let error = false;

        const body = {
          year,
          tags: tags?.split(","),
          artistId: artist?.id,
          featuringArtists: featuringArtists?.map((a) => a.id),
        };

        if (resultData.changeAlbum) {
          const res = await apiClient.editAlbum(data.album.id, body);
          if (!res.success) {
            handleApiError(res.error);
            error = true;
          }
        }

        for (const track of data.tracks) {
          const res = await apiClient.editTrack(track.id, body);
          if (!res.success) {
            handleApiError(res.error);
            error = true;
          }
        }

        if (!error) {
          toast.success("Successfully updated values");
        }
        invalidateAll();
      }}
    >
      <FolderPen />
      Set Common Values
    </SetCommonValuesModal>

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
        <QuickAddButton
          show={!!(data.user && data.user.quickPlaylist)}
          trackId={track.id}
          isInQuickPlaylist={(trackId) => {
            if (!data.quickPlaylistIds) return false;
            return !!data.quickPlaylistIds.find((v) => v === trackId);
          }}
        />

        <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          onclick={async () => {
            await musicManager.clearQueue();
            await musicManager.addFromIds([track.id]);
            musicManager.requestPlay();
          }}
        >
          <Play />
        </Button>

        <EditTrackDetailsModal
          class={cn(
            buttonVariants({
              variant: "ghost",
              size: "icon",
            }),
            "rounded-full",
          )}
          {track}
          onResult={async (resultData) => {
            const res = await apiClient.editTrack(track.id, {
              name: resultData.name,
              otherName: resultData.otherName,
              artistId: resultData.artist.id,
              number: resultData.number,
              year: resultData.year,
              tags: resultData.tags.split(","),
              featuringArtists: resultData.featuringArtists.map((a) => a.id),
            });
            if (!res.success) {
              handleApiError(res.error);
              invalidateAll();
              return;
            }

            toast.success("Successfully updated track");
            invalidateAll();
          }}
        >
          <Pencil />
        </EditTrackDetailsModal>

        <!-- <Button
          class="rounded-full"
          variant="ghost"
          size="icon"
          href="edit/tracks/{track.id}"
        >
          <Pencil />
        </Button> -->

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
              <DropdownMenu.Item
                onSelect={async () => {
                  openConfirmDeleteTrack = { open: true, trackId: track.id };
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

<EditAlbumDetailsModal
  album={data.album}
  bind:open={openEditAlbumDetails}
  onResult={async (resultData) => {
    const res = await apiClient.editAlbum(data.album.id, {
      name: resultData.name,
      otherName: resultData.otherName,
      year: resultData.year ?? 0,
      artistId: resultData.artist.id,
      tags: resultData.tags.split(","),
      featuringArtists: resultData.featuringArtists.map((a) => a.id),
    });
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Successfully updated album");
    invalidateAll();
  }}
/>

<ChangeAlbumCoverModal
  bind:open={openChangeAlbumCoverModal}
  onResult={async (formData) => {
    const toastId = toast.loading("Uploading image...");

    const res = await apiClient.changeAlbumCover(data.album.id, formData);
    if (!res.success) {
      toast.dismiss(toastId);

      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.dismiss(toastId);
    toast.success("Successfully uploaded image, force reload page to see it");
  }}
/>

<ConfirmModal
  bind:open={openConfirmDeleteAlbum}
  removeTrigger
  confirmDelete
  onResult={async () => {
    const res = await apiClient.deleteAlbum(data.album.id);
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Successfully deleted album");
    goto("/albums", { invalidateAll: true });
  }}
/>

<ConfirmModal
  bind:open={openConfirmDeleteTrack.open}
  removeTrigger
  confirmDelete
  onResult={async () => {
    const res = await apiClient.deleteTrack(openConfirmDeleteTrack.trackId);
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Successfully deleted track");
    invalidateAll();
  }}
/>

<ImportTracksModal
  open={true}
  onResult={async (resultData) => {
    for (const file of resultData.file) {
      const body: UploadTrackBody = {
        name: file.name,
        albumId: data.album.id,
        artistId: data.album.artistId,
        number: 0,
        year: 0,
        tags: [],
        featuringArtists: [],
        otherName: "",
      };

      const trackData = new FormData();
      trackData.set("body", JSON.stringify(body));
      trackData.set("track", file);

      const res = await apiClient.uploadTrack(trackData);
      if (!res.success) {
        handleApiError(res.error);
        invalidateAll();
        return;
      }
      toast.success("Successfully uploaded tracks");
    }

    toast.success("Successfully uploaded tracks");
    invalidateAll();
  }}
/>
