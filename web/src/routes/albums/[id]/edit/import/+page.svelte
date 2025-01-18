<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { UploadTrackBody } from "$lib/api/types.js";
  import EditTrackItem from "$lib/components/EditTrackItem.svelte";
  import type { EditTrackData, UIArtist } from "$lib/types.js";
  import { formatError } from "$lib/utils.js";
  import { Breadcrumb, Button, Card, Separator } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();

  const albumArtist = $state<UIArtist>({
    name: data.album.artistName.default,
    id: data.album.artistId,
  });

  type Track = EditTrackData & {
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
    uploadState.uploading = true;
    uploadState.numTracks = tracks.length;
    uploadState.currentTrack = 1;
    for (const track of tracks) {
      console.log(track);
      const body: UploadTrackBody = {
        name: track.name,
        otherName: track.otherName,
        albumId: data.album.id,
        artistId: track.artist.id,
        number: track.num ?? 0,
        year: track.year ?? 0,
        featuringArtists: track.featuringArtists.map((a) => a.id),
        tags: track.tags.split(","),
      };

      const trackData = new FormData();
      trackData.set("body", JSON.stringify(body));
      trackData.set("track", track.file);

      const res = await apiClient.uploadTrack(trackData);
      if (!res.success) {
        handleApiError(res.error);
        return;
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
      {#each tracks as track, i}
        <p>{track.file.name}</p>
        <EditTrackItem {apiClient} bind:track={tracks[i]} />
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
        disabled={tracks.length <= 0 || uploadState.uploading}>Upload</Button
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
          year: data.album.year ?? undefined,
          otherName: "",
          featuringArtists: [],
          file,
          artist: albumArtist,
          tags: "",
        });
      }
    }
  }}
/>
