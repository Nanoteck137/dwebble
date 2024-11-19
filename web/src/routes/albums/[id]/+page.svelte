<script lang="ts">
  import { Button, buttonVariants } from "$lib/components/ui/button";
  import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
  import Separator from "$lib/components/ui/separator/separator.svelte";
  import { musicManager } from "$lib/music-manager";
  import { cn, formatTime, trackToMusicTrack } from "$lib/utils";
  import { EllipsisVertical, ListPlus, Pencil, Play } from "lucide-svelte";

  let { data } = $props();
</script>

{#snippet header()}
  <div class="w-64">
    <img
      class="aspect-square w-64 rounded object-cover"
      src={data.album.coverArt.large}
      alt=""
    />
    <div class="h-2"></div>

    <p class="line-clamp-2 text-center font-bold" title={data.album.name}>
      {data.album.name}
    </p>
    <p class="line-clamp-1 text-center text-xs" title={data.album.artistName}>
      {data.album.artistName}
    </p>

    <div class="h-2"></div>
    <div class="flex items-center justify-center gap-5">
      <Button class="rounded-full" size="icon" variant="outline">
        <ListPlus />
      </Button>

      <Button
        class="rounded-full p-2"
        size="icon-lg"
        variant="default"
        onclick={() => {
          musicManager.clearQueue();

          data.tracks.forEach((t) =>
            musicManager.addTrackToQueue(trackToMusicTrack(t)),
          );
        }}
      >
        <Play />
      </Button>

      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={cn(
            buttonVariants({ variant: "outline", size: "icon" }),
            "rounded-full",
          )}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
          <DropdownMenu.Group>
            <DropdownMenu.Item>
              <a
                class="flex w-full items-center gap-2"
                href="/albums/{data.album.id}/edit"
              >
                <Pencil size="16" />
                Edit Album
              </a>
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
    <div class="h-2"></div>
  </div>
{/snippet}

<div class="flex flex-col gap-2">
  <div class="flex justify-center md:fixed md:h-full">
    {@render header()}
  </div>
  <div class="md:ml-5 md:pl-64">
    <div class="flex flex-col gap-2">
      {#each data.tracks as track, i}
        <div class="group flex items-center gap-2 py-1 pr-4">
          <p class="w-10 text-right text-sm font-medium group-hover:hidden">
            {#if track.number}
              <span>{track.number}.</span>
            {/if}
          </p>
          <button
            class="hidden h-10 w-10 items-center justify-end group-hover:flex"
            onclick={() => {
              musicManager.clearQueue();

              data.tracks.forEach((t) =>
                musicManager.addTrackToQueue(trackToMusicTrack(t), false),
              );

              const index = data.tracks.findIndex((t) => t.id == track.id);
              musicManager.setQueueIndex(index);

              musicManager.requestPlay();
            }}
          >
            <Play size="20" />
          </button>
          <div class="flex flex-grow flex-col py-1">
            <p class="line-clamp-1 text-sm font-medium" title={track.name}>
              {track.name}
            </p>
            <a
              class="line-clamp-1 w-fit text-xs hover:underline"
              title={track.artistName}
              href={`/artists/${track.artistId}`}
            >
              {track.artistName}
            </a>
          </div>
          <div class="flex items-center gap-2">
            <p class="text-xs">
              {formatTime(track.duration ?? 0)}
            </p>
            <!-- <button
              class="hidden rounded-full p-1 hover:bg-black/20 group-hover:block"
            >
            </button> -->

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
                        musicManager.clearQueue();

                        data.tracks.forEach((t) =>
                          musicManager.addTrackToQueue(
                            trackToMusicTrack(t),
                            false,
                          ),
                        );

                        musicManager.setQueueIndex(i);
                        musicManager.requestPlay();
                      }}
                    >
                      <Play size="16" />
                      Play
                    </button>
                  </DropdownMenu.Item>
                </DropdownMenu.Group>
              </DropdownMenu.Content>
            </DropdownMenu.Root>
          </div>
        </div>
        <Separator />
      {/each}
    </div>
  </div>
</div>
