<script lang="ts">
  import { AlertDialog, Button } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";
  import { fade } from "svelte/transition";

  export type Props = {
    title: string;
    description?: string;
    confirmDelete?: boolean;
  };

  const {
    title,
    description,
    isOpen,
    confirmDelete,
    close,
  }: Props & ModalProps<boolean> = $props();
</script>

<AlertDialog.Root
  controlledOpen
  open={isOpen}
  onOpenChange={(v) => {
    if (!v) {
      close(false);
    }
  }}
>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>{title}</AlertDialog.Title>
      {#if description}
        <AlertDialog.Description>
          {description}
        </AlertDialog.Description>
      {/if}
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
      <Button
        variant={confirmDelete ? "destructive" : "default"}
        onclick={() => {
          close(true);
        }}
      >
        {confirmDelete ? "Delete" : "Confirm"}
      </Button>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
