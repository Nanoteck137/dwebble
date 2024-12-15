<script lang="ts">
  import type { ModalQueryArtist, ModalState } from "$lib/modal.svelte";
  import type { QueryArtist } from "$lib/types";
  import { Button, Input, ScrollArea } from "@nanoteck137/nano-ui";

  interface Props {
    data: ModalQueryArtist;
    modalState: ModalState;
  }

  const { data, modalState }: Props = $props();

  let artist: QueryArtist | undefined = $state();

  let currentQuery = $state("");
  let queryResults = $state<QueryArtist[]>([]);

  let timer: NodeJS.Timeout;
  function onInput(e: Event) {
    const target = e.target as HTMLInputElement;
    const current = target.value;

    queryResults = [];
    currentQuery = current;

    clearTimeout(timer);
    timer = setTimeout(async () => {
      const res = await data.apiClient.searchArtists({
        query: {
          query: current,
        },
      });

      if (res.success) {
        queryResults = res.data.artists.map((artist) => ({
          id: artist.id,
          name: artist.name,
        }));
      }
    }, 500);
  }
</script>

<!-- svelte-ignore a11y_consider_explicit_label -->
<button
  class="fixed inset-0 z-[1000] bg-black/70"
  onclick={() => {
    modalState.popModal();
  }}
></button>
<div
  class="fixed left-[50%] top-[50%] z-[1000] grid w-full max-w-lg translate-x-[-50%] translate-y-[-50%] gap-4 border bg-background p-6 shadow-lg duration-200 data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[state=closed]:slide-out-to-left-1/2 data-[state=closed]:slide-out-to-top-[48%] data-[state=open]:slide-in-from-left-1/2 data-[state=open]:slide-in-from-top-[48%] sm:rounded-lg"
>
  <p class="text-lg font-semibold">{data.title ?? "Search Artist"}</p>
  <Input oninput={onInput} placeholder="Search..." />
  {#if currentQuery.length > 0}
    <Button
      type="submit"
      variant="secondary"
      onclick={() => {
        modalState.popModal();
        // TODO(patrik): Create new artist here
      }}
    >
      New Artist: {currentQuery}
    </Button>
  {/if}

  <ScrollArea class="max-h-36 overflow-y-clip">
    {#each queryResults as result}
      <Button
        type="submit"
        variant="ghost"
        onclick={() => {
          modalState.popModal();
          data.onArtistSelected(result);
        }}
      >
        {result.name}
      </Button>
    {/each}
  </ScrollArea>
</div>
