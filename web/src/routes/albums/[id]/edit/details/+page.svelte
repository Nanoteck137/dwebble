<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient, handleApiError, openArtistQuery } from "$lib";
  import type { UIArtist } from "$lib/types.js";
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { Plus, X } from "lucide-svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let name = $state(data.album.name.default);
  let otherName = $state(data.album.name.other);
  let year = $state(data.album.year);
  let tags = $state(data.album.tags.join(","));
  let artist: UIArtist = $state({
    id: data.album.artistId,
    name: data.album.artistName.default,
  });
  let featuringArtists = $state<UIArtist[]>(
    data.album.featuringArtists.map((a) => ({
      id: a.id,
      name: a.name.default,
    })),
  );

  async function submit() {
    const res = await apiClient.editAlbum(data.album.id, {
      name: name,
      otherName: otherName,
      year: year ?? 0,
      artistId: artist.id,
      tags: tags.split(","),
      featuringArtists: featuringArtists.map((a) => a.id),
    });
    if (!res.success) {
      handleApiError(res.error);
      return;
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
        <Breadcrumb.Page>Edit Details</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

{#snippet cardContent()}
  <div class="flex flex-col gap-2">
    <div class="flex flex-col gap-2 sm:flex-row">
      <div class="flex flex-col gap-2 sm:max-w-24">
        <Label for="year">Year</Label>
        <Input id="year" type="number" bind:value={year} />
      </div>

      <div class="flex w-full flex-col gap-2">
        <Label for="name">Name</Label>
        <Input id="name" type="text" bind:value={name} />
      </div>
    </div>

    <div class="flex w-full flex-col gap-2">
      <Label for="otherName">Other Name</Label>
      <Input id="otherName" type="text" bind:value={otherName} />
    </div>

    <div class="flex w-full flex-col gap-2">
      <Label for="tags">Tags</Label>
      <Input id="tags" type="text" bind:value={tags} />
    </div>

    <div class="flex flex-col gap-2">
      <p>Artist: {artist?.name}</p>
      <p>Artist Id: {artist?.id}</p>
    </div>

    <Button
      variant="outline"
      onclick={async () => {
        const res = await openArtistQuery({});
        if (res) {
          artist = res;
        }
      }}
    >
      Change Artist
    </Button>

    <Label class="flex items-center gap-2">
      Featuring Artists
      <button
        type="button"
        class="hover:cursor-pointer"
        tabindex={7}
        onclick={async () => {
          const res = await openArtistQuery({});
          if (res) {
            const index = featuringArtists.findIndex((a) => a.id === res.id);
            if (index === -1) {
              featuringArtists.push(res);
            }
          }
        }}
      >
        <Plus size="16" />
      </button>
    </Label>

    <div class="flex flex-wrap gap-2">
      {#each featuringArtists as artist, i}
        <p
          class="flex w-fit items-center gap-1 rounded-full bg-white px-2 py-1 text-xs text-black"
          title={`${artist.id}: ${artist.name}`}
        >
          <button
            type="button"
            class="text-red-400 hover:cursor-pointer"
            onclick={() => {
              featuringArtists.splice(i, 1);
            }}
          >
            <X size="16" />
          </button>
          {artist.name}
        </p>
      {/each}
    </div>
  </div>
{/snippet}

<form
  onsubmit={(e) => {
    e.preventDefault();
    submit();
  }}
>
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Edit Album Details</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      {@render cardContent()}
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline">
        Back
      </Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
