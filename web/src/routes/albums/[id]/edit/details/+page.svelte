<script lang="ts">
  import { artistQuery, type Artist } from "$lib";
  import { ApiClient } from "$lib/api/client.js";
  import ArtistQuery from "$lib/components/ArtistQuery.svelte";
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";

  const { data } = $props();

  let currentArtist: Artist = $state({
    id: data.album.artistId,
    name: data.album.artistName.default,
  });

  let { open, artist, currentQuery, queryResults, onInput } = artistQuery(
    () => {
      return new ApiClient(data.apiAddress);
    },
  );

  $effect(() => {
    if ($artist) {
      currentArtist = $artist;
    } else {
      currentArtist = {
        id: data.album.artistId,
        name: data.album.artistName.default,
      };
    }
  });
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums/{data.album.id}">
          {data.album.name}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums/{data.album.id}/edit">
          Edit
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Edit Details</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<form action="?/submitEdit" class="flex flex-col gap-2" method="post">
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Edit Album Details</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-1">
        <div class="flex items-center gap-2">
          <Label class="w-24" for="year">Year</Label>
          <Label for="name">Name</Label>
        </div>
        <div class="flex items-center gap-2">
          <Input
            class="w-24"
            id="year"
            name="year"
            value={data.album.year ?? 0}
            type="number"
            autocomplete="off"
          />
          <Input
            class="w-full"
            id="name"
            name="name"
            value={data.album.name}
            type="text"
            autocomplete="off"
          />
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Label for="other-name">Other Name</Label>
        <Input
          id="other-name"
          name="otherName"
          value={data.album.name.other ?? ""}
          type="text"
        />
      </div>

      <div class="flex flex-col gap-2">
        <p>Artist: {currentArtist?.name}</p>
        <p>Artist Id: {currentArtist?.id}</p>
        <input name="artistId" value={currentArtist?.id} type="hidden" />
        <ArtistQuery
          bind:open={$open}
          currentQuery={$currentQuery}
          queryResults={$queryResults}
          onArtistSelected={(a) => {
            $artist = a;
          }}
          {onInput}
        />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline">
        Back
      </Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
