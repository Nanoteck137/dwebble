<script lang="ts">
  import { artistQuery, createApiClient, openArtistQuery } from "$lib";
  import { ApiClient } from "$lib/api/client.js";
  import ArtistQuery from "$lib/components/ArtistQuery.svelte";
  import Errors from "$lib/components/Errors.svelte";
  import type { UIArtist } from "$lib/types.js";
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import SuperDebug, { superForm } from "sveltekit-superforms";

  const { data } = $props();
  const apiClient = createApiClient(data);

  const { form, errors, enhance } = superForm(data.form, { onError: "apply" });

  let artist = $state<UIArtist>();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/albums">Albums</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>New Album</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<form method="post" use:enhance>
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Create Album</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </div>

      {#if artist}
        <input name="artistId" value={artist.id} type="hidden" />
        <p>Artist: {artist.name}</p>
        <p>Artist Id: {artist.id}</p>
      {/if}

      <Button
        variant="outline"
        onclick={async () => {
          const res = await openArtistQuery({ apiClient });
          if (res) {
            artist = res;
          }
        }}
      >
        Set Artist
      </Button>

      <Errors errors={$errors.artistId} />
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button type="submit">Create Album</Button>
    </Card.Footer>
  </Card.Root>
</form>

<SuperDebug data={$form} />
