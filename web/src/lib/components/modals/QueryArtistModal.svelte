<script lang="ts">
  import type { ModalQueryArtist, ModalState } from "$lib/modal.svelte";
  import type { QueryArtist } from "$lib/types";
  import { Button, Input, ScrollArea } from "@nanoteck137/nano-ui";

  interface Props {
    data: ModalQueryArtist;
    modalState: ModalState;
  }

  const { data, modalState }: Props = $props();

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
      } else {
        console.error(res.error.message);
      }
    }, 500);
  }
</script>

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
  <div class="flex flex-col">
    {#each queryResults as result}
      <Button
        type="submit"
        variant="ghost"
        title={result.id}
        onclick={() => {
          modalState.popModal();
          data.onArtistSelected(result);
        }}
      >
        {result.name}
      </Button>
    {/each}
  </div>
</ScrollArea>
<div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
  <Button
    variant="outline"
    onclick={() => {
      modalState.popModal();
    }}
  >
    Close
  </Button>
</div>
