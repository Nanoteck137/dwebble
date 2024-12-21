<script lang="ts">
  import type { ApiClient } from "$lib/api/client";
  import type { QueryArtist, UIArtist } from "$lib/types";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {
    title?: string;
    placeholder?: string;
  };

  const {
    title,
    placeholder,
    isOpen,
    close,
  }: Props & ModalProps<string | null> = $props();

  let input = $state("");
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
        close(input);
      }}
    >
      <Dialog.Header>
        <Dialog.Title>{title ?? "Enter input"}</Dialog.Title>
      </Dialog.Header>
      <Input {placeholder} bind:value={input} />
      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            close(null);
          }}
        >
          Close
        </Button>
        <Button type="submit">Ok</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
