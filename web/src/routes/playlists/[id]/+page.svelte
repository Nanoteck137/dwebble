<script lang="ts">
  import { createApiClient } from "$lib";
  import Spacer from "$lib/components/Spacer.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListHeader from "$lib/components/track-list/TrackListHeader.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { Breadcrumb, DropdownMenu } from "@nanoteck137/nano-ui";
  import { Pencil } from "lucide-svelte";

  const { data } = $props();
  const apiClient = createApiClient(data);
  const musicManager = getMusicManager();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/playlists">Playlists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.playlist.name}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<TrackListHeader
  name={data.playlist.name}
  onPlay={async () => {
    await musicManager.clearQueue();
    await musicManager.addFromPlaylist(data.playlist.id);
    musicManager.requestPlay();
  }}
>
  {#snippet more()}
    <DropdownMenu.Group>
      <DropdownMenu.Item onSelect={() => {}}>
        <Pencil />
        Edit Playlist
      </DropdownMenu.Item>
    </DropdownMenu.Group>
  {/snippet}
</TrackListHeader>

<Spacer size="md" />

<TrackList
  {apiClient}
  totalTracks={data.items.length}
  tracks={data.items}
  userPlaylists={data.userPlaylists}
  quickPlaylist={data.user?.quickPlaylist}
  isInQuickPlaylist={(trackId) => {
    if (!data.quickPlaylistIds) return false;
    return !!data.quickPlaylistIds.find((v) => v === trackId);
  }}
  onPlay={async () => {
    await musicManager.clearQueue();
    await musicManager.addFromPlaylist(data.playlist.id);
    musicManager.requestPlay();
  }}
  onTrackPlay={async (trackId) => {
    await musicManager.clearQueue();
    await musicManager.addFromPlaylist(data.playlist.id);
    await musicManager.setQueueIndex(
      musicManager.queue.items.findIndex((t) => t.id === trackId),
    );
    musicManager.requestPlay();
  }}
/>
