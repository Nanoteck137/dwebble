<script lang="ts">
  import type { ModalConfirm, ModalState } from "$lib/modal.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { fade, slide } from "svelte/transition";

  interface Props {
    data: ModalConfirm;
    modalState: ModalState;
  }

  const { data, modalState }: Props = $props();
</script>

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
