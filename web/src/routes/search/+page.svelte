<script lang="ts">
  import { goto, invalidateAll, onNavigate } from "$app/navigation";
  import { createApiClient } from "$lib";
  import AlbumListItem from "$lib/components/AlbumListItem.svelte";
  import ArtistListItem from "$lib/components/ArtistListItem.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { formatError } from "$lib/utils";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = createApiClient(data);

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

<div class="h-4"></div>

{#if data.artistError}
  <p class="text-red-400">
    Artist Query Error: {formatError(data.artistError)}
  </p>
{/if}

{#if data.albumError}
  <p class="text-red-400">
    Album Query Error: {formatError(data.albumError)}
  </p>
{/if}

{#if data.trackError}
  <p class="text-red-400">
    Track Query Error: {formatError(data.trackError)}
  </p>
{/if}

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
  <TrackList
    {apiClient}
    hideHeader={true}
    totalTracks={data.tracks.length}
    tracks={data.tracks}
    userPlaylists={data.userPlaylists}
    quickPlaylist={data.user?.quickPlaylist}
    isInQuickPlaylist={(trackId) => {
      if (!data.quickPlaylistIds) return false;
      return !!data.quickPlaylistIds.find((v) => v === trackId);
    }}
  />
{/if}
