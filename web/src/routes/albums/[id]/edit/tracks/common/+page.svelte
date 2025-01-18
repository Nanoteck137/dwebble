<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient, handleApiError, openArtistQuery } from "$lib";
  import type { UIArtist } from "$lib/types.js";
  import {
    Breadcrumb,
    Button,
    Card,
    Checkbox,
    Input,
    Label,
    Switch,
  } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = getApiClient();

  let useYear = $state(false);
  let year = $state("");

  let useTags = $state(false);
  let tags = $state("");

  let changeAlbum = $state(true);

  let artist = $state<UIArtist>();

  async function submit() {
    const tagsArr = tags == "" ? [] : tags.split(",");
    const yearNum = year != "" ? parseInt(year) : 0;
    const artistId = artist ? artist.id : undefined;

    if (changeAlbum) {
      const res = await apiClient.editAlbum(data.album.id, {
        artistId,
        tags: useTags ? tagsArr : undefined,
        year: useYear ? yearNum : undefined,
      });

      if (!res.success) {
        handleApiError(res.error);
        return;
      }
    }

    for (const track of data.tracks) {
      const res = await apiClient.editTrack(track.id, {
        artistId,
        tags: useTags ? tagsArr : undefined,
        year: useYear ? yearNum : undefined,
      });
      if (!res.success) {
        handleApiError(res.error);
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
          {data.album.name.default}
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
    <Card.Content class="flex flex-col gap-2">
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

      <div class="flex items-center gap-2">
        <Switch bind:checked={changeAlbum} />
        <p class="text-sm">Set album values</p>
      </div>

      {#if artist}
        <p>Artist: {artist.name}</p>
        <p>Artist Id: {artist.id}</p>
        <Button
          variant="destructive"
          onclick={() => {
            artist = undefined;
          }}
        >
          Remove
        </Button>
      {/if}

      <Button
        variant="outline"
        onclick={async () => {
          const res = await openArtistQuery({});
          if (res) {
            artist = res;
          }
        }}
      >
        Set Artist
      </Button>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-2">
      <Button variant="outline" href="/albums/{data.album.id}/edit">
        Back
      </Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
