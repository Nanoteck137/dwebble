<script>
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
        <Breadcrumb.Link href="/playlists">Playlists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>New Playlist</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="h-4"></div>

<Card.Root class="mx-auto max-w-[450px]">
  <Card.Header>
    <Card.Title>Create playlist</Card.Title>
  </Card.Header>
  <Card.Content>
    <form method="post">
      <div class="flex flex-col gap-4">
        <div class="flex flex-col gap-2">
          <Label for="playlistName">Name</Label>
          <Input
            id="playlistName"
            name="playlistName"
            type="text"
            autocomplete="off"
          />
        </div>
      </div>
    </form>
  </Card.Content>
  <Card.Footer class="flex justify-end gap-4">
    <Button href="/playlists" variant="outline">Back</Button>
    <Button type="submit">Create</Button>
  </Card.Footer>
</Card.Root>
