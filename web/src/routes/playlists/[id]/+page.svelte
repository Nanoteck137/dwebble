<script lang="ts">
  import { createApiClient } from "$lib";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte.js";

  const { data } = $props();
  const apiClient = createApiClient(data);
  const musicManager = getMusicManager();
</script>

<p>{data.playlist.name}</p>

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
