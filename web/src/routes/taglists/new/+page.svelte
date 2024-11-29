<script lang="ts">
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/taglists">Taglists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>New Taglist</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="h-4"></div>

<form method="post">
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Create Taglist</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="taglistName">Taglist Name</Label>
        <Input
          id="taglistName"
          name="taglistName"
          type="text"
          autocomplete="off"
        />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="taglistFilter">Taglist Filter</Label>
        <Input
          id="taglistFilter"
          name="taglistFilter"
          type="text"
          autocomplete="off"
        />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums" variant="outline">Back</Button>
      <Button type="submit">Create</Button>
    </Card.Footer>
  </Card.Root>
</form>
