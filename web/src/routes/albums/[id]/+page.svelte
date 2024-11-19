<script lang="ts">
  import TrackList from "$lib/components/TrackList/TrackList.svelte";
  import type { PageData } from "./$types";

  interface Props {
    data: PageData;
  }

  let { data }: Props = $props();
</script>

<a href="/albums/{data.album.id}/edit">Edit Album</a>

<TrackList
  cover={data.album.coverArt.large}
  name={data.album.name}
  tracks={data.tracks}
/>
