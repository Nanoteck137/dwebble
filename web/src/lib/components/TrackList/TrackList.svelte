<script lang="ts">
  import { type Track } from "$lib/api/types";
  import { musicManager } from "$lib/music-manager.svelte";
  import { formatTime, trackToMusicTrack } from "$lib/utils";
  import { EllipsisVertical, Play } from "lucide-svelte";
  import TrackListHeader from "./TrackListHeader.svelte";

  type Props = {
    name: string;
    cover: string;
    tracks: Track[];
  };

  const { name, cover, tracks }: Props = $props();
</script>

<div class="flex flex-col gap-2 px-5">
  <div class="flex justify-center md:fixed md:h-full">
    <TrackListHeader
      {name}
      {cover}
      onPlay={() => {
        musicManager.clearQueue();

        tracks.forEach((t) =>
          musicManager.addTrackToQueue(trackToMusicTrack(t)),
        );
      }}
    />
  </div>
  <div class="md:ml-5 md:pl-80">
    <div class="flex flex-col gap-2">
      {#each tracks as track}
        <div class="group flex items-center gap-2 border-b pb-1 pr-4">
          <p class="w-10 text-right font-medium group-hover:hidden">
            {#if track.number}
              <span>{track.number}.</span>
            {/if}
          </p>
          <button
            class="hidden h-10 w-10 items-center justify-end group-hover:flex"
            onclick={() => {
              musicManager.clearQueue();

              tracks.forEach((t) =>
                musicManager.addTrackToQueue(trackToMusicTrack(t), false),
              );

              const index = tracks.findIndex((t) => t.id == track.id);
              musicManager.setQueueIndex(index);

              musicManager.requestPlay();
            }}
          >
            <Play size="30" />
          </button>
          <div class="flex flex-grow flex-col">
            <p class="line-clamp-1 font-medium" title={track.name}>
              {track.name}
            </p>
            <a
              class="line-clamp-1 w-fit text-sm hover:underline"
              title={track.artistName}
              href={`/artists/${track.artistId}`}
            >
              {track.artistName}
            </a>
          </div>
          <div>
            <p class="group-hover:hidden">
              {formatTime(track.duration ?? 0)}
            </p>
            <button
              class="hidden rounded-full p-1 hover:bg-black/20 group-hover:block"
            >
              <EllipsisVertical size="30" />
            </button>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>
