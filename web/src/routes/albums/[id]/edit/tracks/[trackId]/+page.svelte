<script lang="ts">
  import { goto } from "$app/navigation";
  import { artistQuery, type Artist } from "$lib";
  import { ApiClient } from "$lib/api/client.js";
  import ArtistQuery from "$lib/components/ArtistQuery.svelte";
  import { musicManager } from "$lib/music-manager";
  import { trackToMusicTrack } from "$lib/utils";
  import { Button, Card, Input, Label } from "@nanoteck137/nano-ui";
  import { Play } from "lucide-svelte";

  const { data } = $props();

  const apiClient = new ApiClient(data.apiAddress);
  apiClient.setToken(data.userToken);

  let trackNumber = $state(data.track.number?.toString() ?? "");
  let trackName = $state(data.track.name.default);
  let trackOtherName = $state(data.track.name.other ?? "");
  let trackYear = $state(data.track.year?.toString() ?? "");
  let trackTags = $state(data.track.tags.join(","));

  let currentArtist: Artist = $state({
    id: data.track.artistId,
    name: data.track.artistName.default,
  });

  let { open, artist, currentQuery, queryResults, onInput } = artistQuery(
    () => {
      return apiClient;
    },
  );

  $effect(() => {
    if ($artist) {
      currentArtist = $artist;
    } else {
      currentArtist = {
        id: data.track.artistId,
        name: data.track.artistName.default,
      };
    }
  });

  async function submit() {
    const num = trackNumber !== "" ? parseInt(trackNumber) : 0;
    console.log(num);

    const year = trackYear !== "" ? parseInt(trackYear) : 0;
    const tags = trackTags !== "" ? trackTags.split(",") : [];

    const artistId = currentArtist.id;

    const res = await apiClient.editTrack(data.track.id, {
      number: num,
      name: trackName,
      otherName: trackOtherName,
      year,
      tags,
      artistId,
    });
    if (!res.success) {
      // TODO(patrik): Toast
      throw res.error.message;
    }

    goto(`/albums/${data.album.id}/edit`, { invalidateAll: true });
  }
</script>

<form
  onsubmit={(e) => {
    e.preventDefault();
    submit();
  }}
>
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Edit Track</Card.Title>
      <Card.Description>{data.track.name.default}</Card.Description>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2 py-2">
        <div class="flex gap-2">
          <button
            type="button"
            onclick={() => {
              musicManager.clearQueue();
              musicManager.addTrackToQueue(trackToMusicTrack(data.track));
            }}
          >
            <Play size="18" />
          </button>
          <p>{data.track.name.default}</p>
        </div>

        <div class="flex flex-col gap-2">
          <div class="flex items-center gap-2">
            <Label class="w-24" for="trackNumber">Number</Label>
            <Label for="trackName">Track Name</Label>
          </div>
          <div class="flex items-center gap-2">
            <Input
              class="w-24"
              id="trackNumber"
              bind:value={trackNumber}
              type="number"
            />
            <Input
              class="w-full"
              id="trackName"
              bind:value={trackName}
              type="text"
              autocomplete="off"
            />
          </div>
        </div>

        <div class="flex flex-col gap-2">
          <Label for="trackOtherName">Other Name</Label>
          <Input
            class="w-full"
            id="trackOtherName"
            bind:value={trackOtherName}
            type="text"
            autocomplete="off"
          />
        </div>

        <div class="flex flex-col gap-2">
          <div class="flex items-center gap-2">
            <Label class="w-24" for="trackYear">Year</Label>
            <Label for="trackTags">Tags</Label>
          </div>
          <div class="flex items-center gap-2">
            <Input
              class="w-24"
              id="trackYear"
              bind:value={trackYear}
              type="number"
            />
            <Input
              class="w-full"
              id="trackTags"
              bind:value={trackTags}
              type="text"
            />
          </div>

          {#if currentArtist}
            <p>Artist: {currentArtist.name}</p>
            <p>Artist Id: {currentArtist.id}</p>
          {/if}

          <ArtistQuery
            bind:open={$open}
            currentQuery={$currentQuery}
            queryResults={$queryResults}
            onArtistSelected={(a) => {
              $artist = a;
            }}
            {onInput}
          />
        </div>
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline"
        >Back</Button
      >
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
