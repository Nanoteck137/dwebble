<script lang="ts">
  import type { Track } from "$lib/api/types";
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

<div class="flex flex-col gap-2 border-b pb-2">
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
    <p>{track.name}:</p>
  </div>

  <div class="flex flex-col gap-1">
    <div class="flex items-center gap-2">
      <label class="w-24 text-sm" for={id("trackNumber")}>Number</label>
      <label class="text-sm" for={id("trackName")}>Track Name</label>
    </div>
    <div class="flex items-center gap-2">
      <input
        class="w-24 rounded bg-[--bg-color] text-xs"
        id={id("trackNumber")}
        name="trackNumber"
        value={track.number}
        type="number"
      />
      <input
        class="w-full rounded bg-[--bg-color] text-xs"
        id={id("trackName")}
        name="trackName"
        value={track.name}
        type="text"
        autocomplete="off"
      />
    </div>
  </div>

  <div class="flex flex-col gap-1">
    <div class="flex items-center gap-2">
      <label class="w-24 text-sm" for={id("trackYear")}>Year</label>
      <label class="text-sm" for={id("trackTags")}>Tags</label>
    </div>
    <div class="flex items-center gap-2">
      <input
        class="w-24 rounded bg-[--bg-color] text-xs"
        id={id("trackYear")}
        name="trackYear"
        value={track.year}
        type="number"
      />
      <input
        class="w-full rounded bg-[--bg-color] text-xs"
        id={id("trackTags")}
        name="trackTags"
        value={track.tags.join(",")}
        type="text"
      />
    </div>
  </div>
</div>
