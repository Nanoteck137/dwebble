<script lang="ts">
  import type { Track } from "$lib/api/types";
  import { Input, Label } from "@nanoteck137/nano-ui";
  import { musicManager } from "$lib/music-manager";
  import { trackToMusicTrack } from "$lib/utils";
  import { Play } from "lucide-svelte";

  type Props = {
    track: Track;
  };

  const { track }: Props = $props();

  function id(s: string) {
    return `${track.id}-${s}`;
  }
</script>

<input name="trackId" value={track.id} type="hidden" />

<div id="track-{track.id}" class="flex flex-col gap-2 py-2">
  <div class="flex gap-2">
    <button
      type="button"
      onclick={() => {
        musicManager.clearQueue();
        musicManager.addTrackToQueue(trackToMusicTrack(track));
      }}
    >
      <Play size="18" />
    </button>
    <p>{track.name}</p>
  </div>

  <div class="flex flex-col gap-1">
    <div class="flex items-center gap-2">
      <Label class="w-24" for={id("trackNumber")}>Number</Label>
      <Label for={id("trackName")}>Track Name</Label>
    </div>
    <div class="flex items-center gap-2">
      <Input
        class="w-24"
        id={id("trackNumber")}
        name="trackNumber"
        value={track.number}
        type="number"
        tabindex={1}
      />
      <Input
        class="w-full"
        id={id("trackName")}
        name="trackName"
        value={track.name}
        type="text"
        autocomplete="off"
        tabindex={2}
      />
    </div>
  </div>

  <div class="flex flex-col gap-2">
    <div class="flex items-center gap-2">
      <Label class="w-24" for={id("trackYear")}>Year</Label>
      <Label for={id("trackTags")}>Tags</Label>
    </div>
    <div class="flex items-center gap-2">
      <Input
        class="w-24"
        id={id("trackYear")}
        name="trackYear"
        value={track.year}
        type="number"
        tabindex={3}
      />
      <Input
        class="w-full"
        id={id("trackTags")}
        name="trackTags"
        value={track.tags.join(",")}
        type="text"
        tabindex={4}
      />
    </div>

    <div class="flex flex-col gap-1">
      <Label for={id("trackArtist")}>Artist</Label>
      <Input
        class="w-full"
        id={id("trackArtist")}
        name="trackArtist"
        value={track.artistName}
        type="text"
        tabindex={5}
      />
    </div>
  </div>
</div>
