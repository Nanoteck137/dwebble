<script lang="ts">
  import { goto } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import { formatError } from "$lib/utils.js";
  import { Breadcrumb, Button, Card } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  const { data } = $props();
  const apiClient = getApiClient();
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
    <p>You are about to delete '{data.album.name.default}'</p>
  </Card.Content>
  <Card.Footer class="flex justify-end gap-2">
    <Button href="/albums/{data.album.id}/edit" variant="outline">Back</Button>
    <Button
      variant="destructive"
      onclick={async () => {
        const res = await apiClient.deleteAlbum(data.album.id);
        if (!res.success) {
          handleApiError(res.error);
          return;
        }

        goto("/albums", { invalidateAll: true });
      }}
    >
      Delete
    </Button>
  </Card.Footer>
</Card.Root>
