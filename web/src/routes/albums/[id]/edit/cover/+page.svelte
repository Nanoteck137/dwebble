<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = getApiClient();

  async function submit(formData: FormData) {
    const res = await apiClient.changeAlbumCover(data.album.id, formData);
    if (!res.success) {
      handleApiError(res.error);
      return;
    }

    // TODO(patrik): Not the best solution, but i need the browser to
    // refresh for the new image to show
    window.location.href = `/albums/${data.album.id}/edit`;
    // goto(`/albums/${data.album.id}/edit`, { invalidateAll: true });
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
        <Breadcrumb.Page>Change Cover</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<form
  onsubmit={(e) => {
    e.preventDefault();
    const formData = new FormData(e.target as HTMLFormElement);
    submit(formData);
  }}
>
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Change Album Cover</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="coverArt">Cover Art</Label>
        <Input
          id="coverArt"
          name="cover"
          type="file"
          accept="image/png,image/jpeg"
        />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline"
        >Back</Button
      >
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
