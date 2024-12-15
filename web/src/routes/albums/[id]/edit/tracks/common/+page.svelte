<script lang="ts">
  import { goto } from "$app/navigation";
  import { artistQuery } from "$lib";
  import { ApiClient } from "$lib/api/client.js";
  import ArtistQuery from "$lib/components/ArtistQuery.svelte";
  import {
    Breadcrumb,
    Button,
    Card,
    Checkbox,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";

  const { data } = $props();

  let useYear = $state(false);
  let year = $state("");

  let useTags = $state(false);
  let tags = $state("");

  // svelte-ignore non_reactive_update
  let { open, artist, currentQuery, queryResults, onInput } = artistQuery(
    () => {
      return new ApiClient(data.apiAddress);
    },
  );

  async function submit() {
    const apiClient = new ApiClient(data.apiAddress);
    apiClient.setToken(data.userToken);

    const tagsArr = tags == "" ? [] : tags.split(",");
    const yearNum = year != "" ? parseInt(year) : 0;
    const artistId = $artist ? $artist.id : undefined;

    for (const track of data.tracks) {
      const res = await apiClient.editTrack(track.id, {
        artistId,
        tags: useTags ? tagsArr : undefined,
        year: useYear ? yearNum : undefined,
      });
      if (!res.success) {
        throw res.error.message;
      }
    }

    goto(`/albums/${data.album.id}/edit`, { invalidateAll: true });
  }
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
        <Breadcrumb.Page>Set Common Values</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<form
  onsubmit={async (e) => {
    e.preventDefault();
    submit();
  }}
>
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Set Common Values</Card.Title>
    </Card.Header>
    <Card.Content>
      <div class="flex w-32 flex-col gap-2">
        <Label for="year">Year</Label>

        <div class="flex items-center gap-2">
          <Checkbox class="h-5 w-5" bind:checked={useYear} />
          <Input
            id="year"
            type="number"
            disabled={!useYear}
            bind:value={year}
          />
        </div>
      </div>

      <div class="h-2"></div>

      <div class="flex flex-col gap-2">
        <Label for="tags">Tags</Label>

        <div class="flex items-center gap-2">
          <Checkbox class="h-5 w-5" bind:checked={useTags} />
          <Input
            id="tags"
            type="text"
            autocomplete="off"
            disabled={!useTags}
            bind:value={tags}
          />
        </div>
      </div>

      {#if $artist}
        <p>Artist: {$artist.name}</p>
        <p>Artist Id: {$artist.id}</p>
        <Button
          variant="destructive"
          onclick={() => {
            $artist = undefined;
          }}
        >
          Remove
        </Button>
      {/if}

      <div class="h-4"></div>

      <ArtistQuery
        bind:open={$open}
        currentQuery={$currentQuery}
        queryResults={$queryResults}
        onArtistSelected={(a) => {
          $artist = a;
        }}
        {onInput}
      />
    </Card.Content>
    <Card.Footer class="flex justify-end gap-2">
      <Button variant="outline" href="/albums/{data.album.id}/edit">
        Back
      </Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
