<script lang="ts">
  import { openArtistQuery } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import type { EditTrackData } from "$lib/types";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
  import { Plus, X } from "lucide-svelte";

  type Props = {
    apiClient: ApiClient;
    track: EditTrackData;
  };

  const { apiClient, track = $bindable() }: Props = $props();
</script>

<div class="flex flex-col gap-2">
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
      <Input id="trackName" type="text" tabindex={2} bind:value={track.name} />
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
    <Input id="trackTags" type="text" tabindex={5} bind:value={track.tags} />
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
          const index = track.extraArtists.findIndex((a) => a.id === res.id);
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
</div>
