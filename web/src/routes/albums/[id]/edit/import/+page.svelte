<script lang="ts">
  import { goto } from "$app/navigation";
  import { openArtistQuery, type Artist } from "$lib";
  import { ApiClient } from "$lib/api/client";
  import type { UploadTrackBody } from "$lib/api/types.js";
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { Plus, X } from "lucide-svelte";

  const { data } = $props();

  const apiClient = new ApiClient(data.apiAddress);
  apiClient.setToken(data.userToken);

  const albumArtist = $state<Artist>({
    name: data.album.artistName.default,
    id: data.album.artistId,
  });

  type Track = {
    name: string;
    otherName: string;
    artist: Artist;
    num?: number;
    year?: number;
    tags: string;
    extraArtists: Artist[];
    file: File;
  };

  let uploadState = $state({
    uploading: false,
    currentTrack: 0,
    numTracks: 0,
  });

  let fileSelector = $state<HTMLInputElement>();
  let tracks = $state<Track[]>([]);

  async function submit() {
    const apiClient = new ApiClient(data.apiAddress);
    apiClient.setToken(data.userToken);

    uploadState.uploading = true;
    uploadState.numTracks = tracks.length;
    uploadState.currentTrack = 1;
    for (const track of tracks) {
      const body: UploadTrackBody = {
        name: track.name,
        otherName: track.otherName,
        albumId: data.album.id,
        artistId: track.artist.id,
        number: track.num ?? 0,
        year: track.year ?? 0,
        extraArtists: track.extraArtists.map((a) => a.id),
        tags: track.tags.split(","),
      };

      const trackData = new FormData();
      trackData.set("body", JSON.stringify(body));
      trackData.set("track", track.file);

      const res = await apiClient.uploadTrack(trackData);
      if (!res.success) {
        // TODO(patrik): Toast
        console.log(res.error.message);
      }
      uploadState.currentTrack += 1;
    }
    uploadState.uploading = false;

    goto(`/albums/${data.album.id}/edit`, { invalidateAll: true });
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
        <Breadcrumb.Link href="/albums/{data.album.id}">
          {data.album.name.default}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums/{data.album.id}/edit">
          Edit
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Import Tracks</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<form
  onsubmit={(e) => {
    e.preventDefault();
    submit();
  }}
>
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Import Tracks</Card.Title>
      <Card.Description>{data.album.name.default}</Card.Description>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      {#each tracks as track}
        <p>{track.file.name}</p>
        <div class="flex flex-col gap-2 sm:flex-row">
          <div class="flex flex-col gap-2 sm:max-w-24">
            <Label for="trackNumber">Number</Label>
            <Input
              id="trackNumber"
              type="number"
              bind:value={track.num}
              tabindex={1}
            />
          </div>

          <div class="w-full">
            <Label for="trackName">Name</Label>
            <Input
              id="trackName"
              type="text"
              tabindex={2}
              bind:value={track.name}
            />
          </div>
        </div>

        <div class="flex flex-col gap-2 sm:flex-row">
          <div class="flex flex-col gap-2 sm:max-w-24">
            <Label for="trackYear">Year</Label>
            <Input
              id="trackYear"
              type="number"
              bind:value={track.year}
              tabindex={3}
            />
          </div>

          <div class="flex w-full flex-col gap-2">
            <Label for="trackOtherName">Other Name</Label>
            <Input
              id="trackOtherName"
              type="text"
              tabindex={4}
              bind:value={track.otherName}
            />
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <Label for="trackTags">Tags</Label>
          <Input
            id="trackTags"
            type="text"
            tabindex={5}
            bind:value={track.tags}
          />
        </div>

        <div class="flex flex-1 flex-col gap-2">
          <Label>Artist</Label>
          <Button
            class="justify-start"
            variant="ghost"
            onclick={async () => {
              const res = await openArtistQuery({ apiClient });
              if (res) {
                track.artist = res;
              }
            }}
            tabindex={6}
          >
            {track.artist.name}
          </Button>
        </div>

        <Label class="flex items-center gap-2">
          Extra Artists
          <button
            type="button"
            class="hover:cursor-pointer"
            tabindex={7}
            onclick={async () => {
              const res = await openArtistQuery({ apiClient });
              if (res) {
                const index = track.extraArtists.findIndex(
                  (a) => a.id === res.id,
                );
                if (index === -1) {
                  track.extraArtists.push(res);
                }
              }
            }}
          >
            <Plus size="16" />
          </button>
        </Label>

        <div class="flex flex-wrap gap-2">
          {#each track.extraArtists as artist, i}
            <p
              class="flex w-fit items-center gap-1 rounded-full bg-white px-2 py-1 text-xs text-black"
              title={`${artist.id}: ${artist.name}`}
            >
              <button
                type="button"
                class="text-red-400 hover:cursor-pointer"
                onclick={() => {
                  track.extraArtists.splice(i, 1);
                }}
              >
                <X size="16" />
              </button>
              {artist.name}
            </p>
          {/each}
        </div>

        <Separator />
      {/each}

      <Button
        variant="outline"
        onclick={() => {
          fileSelector?.click();
        }}
      >
        Open File Selector
      </Button>

      {#if uploadState.uploading}
        <p>Uploading: {uploadState.currentTrack} / {uploadState.numTracks}</p>
      {/if}
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline">
        Back
      </Button>
      <Button
        type="submit"
        disabled={tracks.length <= 0 || uploadState.uploading}>Save</Button
      >
    </Card.Footer>
  </Card.Root>
</form>

<input
  class="hidden"
  bind:this={fileSelector}
  type="file"
  multiple
  accept="audio/*"
  onchange={(e) => {
    const input = e.target as HTMLInputElement;
    if (!input.files) return;

    const files = input.files;
    for (let i = 0; i < files.length; i++) {
      const file = files.item(i);
      if (file) {
        tracks.push({
          name: file.name,
          otherName: "",
          extraArtists: [],
          file,
          artist: albumArtist,
          tags: "",
        });
      }
    }
  }}
/>
