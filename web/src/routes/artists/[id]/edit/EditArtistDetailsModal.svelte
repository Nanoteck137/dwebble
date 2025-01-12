<script lang="ts">
  import type { Artist } from "$lib/api/types";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {
    artist: Artist;
  };

  export type Result = {
    name: string;
    otherName: string;
    tags: string;
  };

  const { artist, isOpen, close }: Props & ModalProps<Result | null> =
    $props();

  let result = $state<Result>({
    name: artist.name.default,
    otherName: artist.name.other ?? "",
    tags: artist.tags.join(","),
  });
</script>

<Dialog.Root
  controlledOpen
  open={isOpen}
  onOpenChange={(v) => {
    if (!v) {
      close(null);
    }
  }}
>
  <Dialog.Content>
    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();
        close(result);
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Edit Artist Details</Dialog.Title>
      </Dialog.Header>

      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" bind:value={result.name} />
      </FormItem>

      <FormItem>
        <Label for="otherName">Other Name</Label>
        <Input id="otherName" bind:value={result.otherName} />
      </FormItem>

      <FormItem>
        <Label for="tags">Tags</Label>
        <Input id="tags" bind:value={result.tags} />
      </FormItem>

      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            close(null);
          }}
        >
          Close
        </Button>
        <Button type="submit">Save</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
