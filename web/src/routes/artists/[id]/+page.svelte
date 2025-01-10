<script lang="ts">
  import { Breadcrumb, Button, Separator } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>{data.artist.name.default}</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<Button href="/artists/{data.artist.id}/edit">Edit</Button>
<p>Artist: {data.artist.name.default}</p>

<p>Num Albums: {data.albums.length}</p>
{#each data.albums as album}
  <p>{album.name.default}</p>
{/each}

<Separator />

<p>Num Tracks: {data.trackPage.totalItems}</p>
{#each data.tracks as track}
  <p>{track.name.default}</p>
{/each}
