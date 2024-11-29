<script lang="ts">
  import { Button, Card, Input, Label } from "@nanoteck137/nano-ui";

  const { data } = $props();
</script>

<form class="flex flex-col gap-2" method="post">
  <Card.Root class="mx-auto w-full max-w-[560px]">
    <Card.Header>
      <Card.Title>Edit Album Details</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-1">
        <div class="flex items-center gap-2">
          <Label class="w-24" for="year">Number</Label>
          <Label for="name">Name</Label>
        </div>
        <div class="flex items-center gap-2">
          <Input
            class="w-24"
            id="year"
            name="year"
            value={data.album.year ?? 0}
            type="number"
            autocomplete="off"
          />
          <Input
            class="w-full"
            id="name"
            name="name"
            value={data.album.name}
            type="text"
            autocomplete="off"
          />
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Label for="artist">Artist</Label>
        <Input
          id="artist"
          name="artist"
          value={data.album.artistName}
          type="text"
        />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums/{data.album.id}/edit" variant="outline">
        Back
      </Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
