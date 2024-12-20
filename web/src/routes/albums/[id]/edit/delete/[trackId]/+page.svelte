<script lang="ts">
  import { goto } from "$app/navigation";
  import { createApiClient } from "$lib";
  import { Breadcrumb, Button, Card } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = createApiClient(data);
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
        <Breadcrumb.Page>Delete</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<Card.Root class="mx-auto max-w-[450px]">
  <Card.Header>
    <Card.Title>Are you sure?</Card.Title>
  </Card.Header>
  <Card.Content>
    <p>You are about to delete '{data.album.name}'</p>
  </Card.Content>
  <Card.Footer class="flex justify-end gap-2">
    <Button href="/albums/{data.album.id}/edit" variant="outline">Back</Button>
    <Button
      variant="destructive"
      onclick={async () => {
        // const res = await apiClient.deleteAlbum(data.album.id);
        // if (!res.success) {
        //   // TODO(patrik): Toast
        //   throw res.error;
        // }

        goto(`/albums/${data.album.id}/edit`, { invalidateAll: true });
      }}
    >
      Delete
    </Button>
  </Card.Footer>
</Card.Root>
