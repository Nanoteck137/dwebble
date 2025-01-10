<script lang="ts">
  import { createApiClient } from "$lib";
  import TrackList from "$lib/components/track-list/TrackList.svelte";

  const { data } = $props();
  const apiClient = createApiClient(data);
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
/>
