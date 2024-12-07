<script lang="ts">
  import { enhance } from "$app/forms";
  import { goto, invalidateAll, onNavigate } from "$app/navigation";
  import { Artist } from "$lib/api/types";
  import AlbumListItem from "$lib/components/AlbumListItem.svelte";
  import ArtistListItem from "$lib/components/ArtistListItem.svelte";
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { musicManager } from "$lib/music-manager";
  import { cn, trackToMusicTrack } from "$lib/utils";
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
    <ArtistListItem {artist} link />
  {/each}
{/if}

{#if data.albums.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Albums</p>
    <p class="text-xs">{data.albums.length} album(s)</p>
  </div>

  {#each data.albums as album}
    <AlbumListItem {album} link />
  {/each}
{/if}

{#if data.tracks.length > 0}
  <div class="flex items-center justify-between">
    <p class="text-bold">Tracks</p>
    <p class="text-xs">{data.tracks.length} track(s)</p>
  </div>

  {#each data.tracks as track}
    <TrackListItem
      {track}
      onPlayClicked={() => {
        musicManager.clearQueue();
        musicManager.addTrackToQueue(trackToMusicTrack(track));
      }}
    />
  {/each}
{/if}
