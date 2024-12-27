<script lang="ts">
  import { goto } from "$app/navigation";
  import { createApiClient } from "$lib";
  import EditTrackItem from "$lib/components/EditTrackItem.svelte";
  import { musicManager } from "$lib/music-manager";
  import { type EditTrackData } from "$lib/types.js";
  import { formatError, trackToMusicTrack } from "$lib/utils";
  import { Button, Card } from "@nanoteck137/nano-ui";
  import { Play } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = createApiClient(data);

  let track = $state<EditTrackData>({
    name: data.track.name.default,
    otherName: data.track.name.other ?? "",
    artist: {
      id: data.track.artistId,
      name: data.track.artistName.default,
    },
    num: data.track.number ?? 0,
    year: data.track.year ?? 0,
    tags: data.track.tags.join(","),
    featuringArtists: data.track.featuringArtists.map((a) => ({
      id: a.id,
      name: a.name.default,
    })),
  });

  async function submit() {
    const res = await apiClient.editTrack(data.track.id, {
      name: track.name,
      otherName: track.otherName,
      artistId: track.artist.id,
      number: track.num,
      year: track.year,
      tags: track.tags.split(","),
      featuringArtists: track.featuringArtists.map((a) => a.id),
    });
    if (!res.success) {
      toast.error("Unknown error");
      throw formatError(res.error);
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
        <EditTrackItem {apiClient} bind:track />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline">
        Back
      </Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
