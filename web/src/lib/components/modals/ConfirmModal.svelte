<script lang="ts">
  import type { ModalConfirm, ModalState } from "$lib/modal.svelte";
  import { Button } from "@nanoteck137/nano-ui";

  interface Props {
    data: ModalConfirm;
    modalState: ModalState;
  }

  const { data, modalState }: Props = $props();
</script>

<div class="fixed inset-0 z-[1000] bg-black/70"></div>
<div
  class="fixed left-[50%] top-[50%] z-[1000] grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg"
>
  <div class="flex flex-col gap-2">
    <p class="text-lg font-semibold">{data.title}</p>
    {#if data.description}
      <p class="text-sm text-muted-foreground">{data.description}</p>
    {/if}
  </div>
  <div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
    <Button
      variant="outline"
      onclick={() => {
        modalState.popModal();
      }}
    >
      Close
    </Button>

    {#if data.confirmDelete}
      <Button
        variant="destructive"
        onclick={() => {
          modalState.popModal();
          data.onConfirm?.();
        }}
      >
        Delete
      </Button>
    {:else}
      <Button
        onclick={() => {
          modalState.popModal();
          data.onConfirm?.();
        }}
      >
        Confirm
      </Button>
    {/if}
  </div>
</div>
