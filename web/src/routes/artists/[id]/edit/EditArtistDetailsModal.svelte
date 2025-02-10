<script lang="ts">
  import type { Artist } from "$lib/api/types";
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/new-modals";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";

  export type Props = {
    artist: Artist;
    open: boolean;
  };

  export type Result = {
    name: string;
    otherName: string | null;
    tags: string;
  };

  let {
    artist,
    open = $bindable(),
    onResult,
  }: Props & Modal<Result> = $props();

  let result = $state<Result>({
    name: "",
    otherName: null,
    tags: "",
  });

  $effect(() => {
    if (open) {
      result = {
        name: artist.name.default,
        otherName: artist.name.other,
        tags: artist.tags.join(","),
      };
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();
        onResult(result);
        open = false;
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Edit artist details</Dialog.Title>
      </Dialog.Header>

      <div class="flex flex-col gap-4">
        <FormItem>
          <Label for="name">Name</Label>
          <Input
            id="name"
            type="text"
            autocomplete="off"
            bind:value={result.name}
          />
        </FormItem>

        <FormItem>
          <Label for="otherName">Other Name</Label>
          <Input
            id="otherName"
            type="text"
            autocomplete="off"
            bind:value={result.otherName}
          />
        </FormItem>

        <FormItem>
          <Label for="tags">Tags</Label>
          <Input
            id="tags"
            type="text"
            autocomplete="off"
            bind:value={result.tags}
          />
        </FormItem>
      </div>

      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button type="submit">Save</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
