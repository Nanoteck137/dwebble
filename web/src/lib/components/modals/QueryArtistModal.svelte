<script lang="ts">
  import type { Artist } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import type { QueryArtist } from "$lib/types";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {
    title?: string;
    apiClient: ApiClient;
  };

  const {
    title,
    apiClient,
    isOpen,
    close,
  }: Props & ModalProps<Artist | null> = $props();

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
      const res = await apiClient.searchArtists({
        query: {
          query: current,
        },
      });

      if (res.success) {
        queryResults = res.data.artists.map((artist) => ({
          id: artist.id,
          name: artist.name.default,
        }));
      } else {
        console.error(res.error.message);
      }
    }, 500);
  }
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
    <Dialog.Header>
      <Dialog.Title>{title ?? "Search Artist"}</Dialog.Title>
    </Dialog.Header>
    <Input oninput={onInput} placeholder="Search..." />
    {#if currentQuery.length > 0}
      <Button
        type="submit"
        variant="secondary"
        onclick={() => {
          close(null);
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
              close(result);
              // data.onArtistSelected(result);
            }}
          >
            {result.name}
          </Button>
        {/each}
      </div>
    </ScrollArea>
    <Dialog.Footer>
      <Button
        variant="outline"
        onclick={() => {
          close(null);
        }}
      >
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
