<script lang="ts">
  import TrackList from "$lib/components/TrackList/TrackList.svelte";
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<TrackList
  cover={data.album.coverArt.large}
  name={data.album.name}
  tracks={data.tracks}
/>
