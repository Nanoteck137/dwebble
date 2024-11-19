<script lang="ts">
  import TrackList from "$lib/components/TrackList/TrackList.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<a href="/albums/{data.album.id}/edit">Edit Album</a>

<TrackList
  cover={data.album.coverArt.large}
  name={data.album.name}
  tracks={data.tracks}
/>
