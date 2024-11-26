<script lang="ts">
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import type { ActionData, PageData } from "./$types";

  const { form }: { data: PageData; form: ActionData } = $props();
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

<div class="h-4"></div>

<form method="post">
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Create album</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="albumName">Album Name</Label>
        <Input
          id="albumName"
          name="albumName"
          type="text"
          autocomplete="off"
        />

        {#if form?.errors?.albumName}
          {#each form?.errors?.albumName as error}
            <p class="text-xs text-red-400">{error}</p>
          {/each}
        {/if}
      </div>
      <div class="flex flex-col gap-2">
        <Label for="artistName">Artist Name</Label>
        <Input
          id="artistName"
          name="artistName"
          type="text"
          autocomplete="off"
        />
        {#if form?.errors?.artistName}
          {#each form?.errors?.artistName as error}
            <p class="text-xs text-red-400">{error}</p>
          {/each}
        {/if}
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/albums" variant="outline">Back</Button>
      <Button type="submit">Create</Button>
    </Card.Footer>
  </Card.Root>
</form>
