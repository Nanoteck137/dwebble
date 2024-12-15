<script lang="ts">
  import { Breadcrumb, Button, Card } from "@nanoteck137/nano-ui";

  const { data } = $props();
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
    <Button variant="destructive">Delete</Button>
  </Card.Footer>
</Card.Root>
