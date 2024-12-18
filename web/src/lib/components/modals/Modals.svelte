<script lang="ts">
  import { browser } from "$app/environment";
  import { navigating } from "$app/stores";
  import ConfirmModal from "$lib/components/modals/ConfirmModal.svelte";
  import QueryArtistModal from "$lib/components/modals/QueryArtistModal.svelte";
  import { getModalState } from "$lib/modal.svelte";
  import { fade } from "svelte/transition";

  const modalState = getModalState();

  $effect(() => {
    if (!browser) return;

    if (modalState.modals.length > 0) {
      document.body.classList.add("modal-open");
    } else {
      document.body.classList.remove("modal-open");
    }
  });

  $effect(() => {
    if ($navigating !== null) {
      modalState.modals.forEach((modal) => {
        modalState.removeModal(modal.id);
      });
    }
  });
</script>

{#each modalState.modals as modal (modal.id)}
  <!-- svelte-ignore a11y_consider_explicit_label -->
  <button
    class="fixed inset-0 z-[1000] bg-black/70"
    onclick={() => {
      modalState.popModal();
    }}
  ></button>
  <div
    class="fixed left-[50%] top-[50%] z-[1000] grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg transition ease-in-out sm:rounded-lg"
    transition:fade={{ duration: 150 }}
  >
    {#if modal.data.type === "modal-confirm"}
      <ConfirmModal {modalState} data={modal.data} />
    {:else if modal.data.type === "modal-query-artist"}
      <QueryArtistModal {modalState} data={modal.data} />
    {/if}
  </div>
{/each}
