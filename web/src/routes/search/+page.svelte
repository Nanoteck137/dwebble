<script lang="ts">
  import { enhance } from "$app/forms";
  import { goto, invalidateAll, onNavigate } from "$app/navigation";
  import { cn } from "$lib/utils";
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Pencil, Play, Plus } from "lucide-svelte";

  const { data } = $props();

  async function search(query: string) {
    await goto(`?query=${query}`, {
      invalidateAll: true,
      keepFocus: true,
      replaceState: true,
    });
  }

  let value = "";

  let timer: NodeJS.Timeout;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;
    value = current;

    clearTimeout(timer);
    timer = setTimeout(async () => {
      search(current);
    }, 500);
  }

  // NOTE(patrik): Fix for clicking the search button
  onNavigate((e) => {
    if (e.type === "link" && e.from?.url.pathname === "/search") {
      invalidateAll();
    }
  });
</script>

<form
  action=""
  method="get"
  onsubmit={(e) => {
    e.preventDefault();
    clearTimeout(timer);
    search(value);
  }}
>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-2">
      <Label for="query">Search Query</Label>
      <Input
        id="query"
        name="query"
        autocomplete="off"
        value={data.query}
        oninput={onInput}
      />
    </div>
    <Button type="submit">Search</Button>
  </div>
</form>

{#if data.artists.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Artists</p>
    <p class="text-xs">{data.artists.length} artist(s)</p>
  </div>

  {#each data.artists as artist}
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <img
        class="inline-flex aspect-square min-w-14 max-w-14 items-center justify-center rounded object-cover text-xs"
        src={artist.picture.small}
        alt="cover"
        loading="lazy"
      />

      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <a
            class="line-clamp-1 w-fit font-medium hover:underline"
            href="/artists/{artist.id}"
            title={artist.name}
          >
            {artist.name}
          </a>
        </div>
      </div>

      <div class="flex items-center">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon-lg" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <!-- <DropdownMenu.Item>
              <a
                class="flex h-full w-full items-center gap-2 py-1"
                href="/albums/{album.id}"
              >
                <Pencil size="16" />
                Go to Album
              </a>
            </DropdownMenu.Item> -->
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  {/each}
{/if}

{#if data.albums.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Albums</p>
    <p class="text-xs">{data.albums.length} album(s)</p>
  </div>

  {#each data.albums as album}
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <img
        class="inline-flex aspect-square min-w-14 max-w-14 items-center justify-center rounded object-cover text-xs"
        src={album.coverArt.small}
        alt="cover"
        loading="lazy"
      />

      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <a
            class="line-clamp-1 w-fit font-medium hover:underline"
            href="/albums/{album.id}"
            title={album.name}
          >
            {album.name}
          </a>
        </div>

        <a
          class="line-clamp-1 w-fit text-sm font-light hover:underline"
          title={album.artistName}
          href={`/artists/${album.artistId}`}
        >
          {album.artistName}
        </a>
      </div>

      <div class="flex items-center">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon-lg" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item>
                <a
                  class="flex h-full w-full items-center gap-2 py-1"
                  href="/albums/{album.id}"
                >
                  <Pencil size="16" />
                  Go to Album
                </a>
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  {/each}
{/if}

{#if data.tracks.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Tracks</p>
    <p class="text-xs">{data.tracks.length} track(s)</p>
  </div>

  {#each data.tracks as track, i}
    <div class="flex items-center gap-2 border-b py-2 pr-2">
      <div class="group relative">
        <img
          class="inline-flex aspect-square min-w-14 max-w-14 items-center justify-center rounded object-cover text-xs"
          src={track.coverArt.small}
          alt="cover"
          loading="lazy"
        />

        <button
          class={`absolute bottom-0 left-0 right-0 top-0 hidden items-center justify-center rounded bg-black/80 group-hover:flex`}
          onclick={() => {
            // musicManager.clearQueue();
            // for (const track of tracks) {
            //   musicManager.addTrackToQueue(trackToMusicTrack(track), false);
            // }
            // musicManager.setQueueIndex(i);
            // musicManager.requestPlay();
          }}
        >
          <Play size="25" />
        </button>
      </div>
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <p class="line-clamp-1 w-fit font-medium" title={track.name}>
            {track.name}
          </p>
        </div>

        <a
          class="line-clamp-1 w-fit text-sm font-light hover:underline"
          title={track.artistName}
          href={`/artists/${track.artistId}`}
        >
          {track.artistName}
        </a>

        <p class="line-clamp-1 text-xs font-light">
          {#if track.tags.length > 0}
            {track.tags.join(", ")}
          {:else}
            No Tags
          {/if}
        </p>
      </div>
      <div class="flex items-center">
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon-lg" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item>
                <form
                  class="jusitfy-center flex items-center"
                  action="?/quickAddToPlaylist"
                  method="post"
                  use:enhance
                >
                  <input type="hidden" name="trackId" value={track.id} />
                  <button
                    class="flex w-full items-center gap-2 py-1"
                    title="Quick Add"
                  >
                    <Plus size="16" />
                    Quick add to Playlist
                  </button>
                </form>
              </DropdownMenu.Item>
              <DropdownMenu.Item>
                <a
                  class="flex h-full w-full items-center gap-2 py-1"
                  href="/albums/{track.albumId}"
                >
                  <Pencil size="16" />
                  Go to Album
                </a>
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  {/each}
{/if}
