<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import { formatError } from "$lib/utils.js";
  import { Breadcrumb, Button, Card, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();

  let fileSelector = $state<HTMLInputElement>();
  let files = $state<FileList>();

  async function submit() {
    console.log(files);
    const file = files?.item(0);
    if (file) {
      const formData = new FormData();
      formData.set("cover", file);

      const res = await apiClient.changeArtistPicture(
        data.artist.id,
        formData,
      );
      if (!res.success) {
        handleApiError(res.error);
        return;
      }
      // TODO(patrik): Not the best solution, but i need the browser to
      // refresh for the new image to show
      window.location.href = `/artists/${data.artist.id}/edit`;
      // goto(`/albums/${data.album.id}/edit`, { invalidateAll: true });
    }
  }
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists/{data.artist.id}">
          {data.artist.name.default}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists/{data.artist.id}/edit">
          Edit
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Change Picture</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<form
  onsubmit={(e) => {
    e.preventDefault();
    submit();
  }}
>
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Change Artist Picture</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <Label>Picture</Label>
      {#if files && files.length >= 1}
        <p>{files!.item(0)!.name}</p>
      {/if}

      <Button
        variant="outline"
        onclick={() => {
          fileSelector?.click();
        }}
      >
        Select Picture
      </Button>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/artists/{data.artist.id}/edit" variant="outline">
        Back
      </Button>
      <Button type="submit" disabled={files ? files.length <= 0 : true}>
        Save
      </Button>
    </Card.Footer>
  </Card.Root>
</form>

<input
  bind:this={fileSelector}
  class="hidden"
  type="file"
  accept="image/png,image/jpeg"
  bind:files
/>
