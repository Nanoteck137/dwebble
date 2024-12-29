<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { createApiClient } from "$lib";
  import QuickAddButton from "$lib/components/QuickAddButton.svelte";
  import TrackList from "$lib/components/track-list/TrackList.svelte";
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { musicManager } from "$lib/music-manager";
  import { cn, formatError, shuffle, trackToMusicTrack } from "$lib/utils.js";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import {
    DiscAlbum,
    EllipsisVertical,
    FileHeart,
    Play,
    Shuffle,
    Users,
  } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = createApiClient(data);
</script>

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
