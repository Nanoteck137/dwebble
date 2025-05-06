<script lang="ts">
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/taglists">Taglists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/taglists/{data.taglist.id}">
          {data.taglist.name}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Edit</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="h-4"></div>

<form method="post">
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Edit Taglist</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="taglistName">Taglist Name</Label>
        <Input
          id="taglistName"
          name="taglistName"
          autocomplete="off"
          value={data.taglist.name}
        />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="taglistFilter">Filter</Label>
        <Input
          id="taglistFilter"
          name="taglistFilter"
          autocomplete="off"
          value={data.taglist.filter}
        />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/taglists/{data.taglist.id}" variant="outline">Back</Button
      >
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
