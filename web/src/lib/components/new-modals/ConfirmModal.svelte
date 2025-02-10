<script lang="ts">
  import type { Modal } from "$lib/components/new-modals";
  import { Button, Dialog, Input } from "@nanoteck137/nano-ui";

  export type Props = {
    open?: boolean;

    removeTrigger?: boolean;

    title?: string;
    description?: string;
    confirmDelete?: boolean;
  };

  let {
    open = $bindable(false),

    removeTrigger,

    title,
    description,
    confirmDelete,

    class: className,
    children,
    onResult,
  }: Props & Modal<unknown> = $props();
</script>

<Dialog.Root bind:open>
  {#if !removeTrigger}
    <Dialog.Trigger class={className}>
      {@render children?.()}
    </Dialog.Trigger>
  {/if}

  <Dialog.Content>
    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();
        onResult(true);
        open = false;
      }}
    >
      <Dialog.Header>
        <Dialog.Title>{title ?? "Are you sure?"}</Dialog.Title>
        {#if description}
          <Dialog.Description>{description}</Dialog.Description>
        {/if}
      </Dialog.Header>
      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button
          variant={confirmDelete ? "destructive" : "default"}
          type="submit"
        >
          {confirmDelete ? "Delete" : "Confirm"}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
