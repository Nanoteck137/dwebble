<script lang="ts">
  import { openInput } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import { GetArtistById } from "$lib/api/types";
  import type { QueryArtist, UIArtist } from "$lib/types";
  import { formatError } from "$lib/utils";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
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
  }: Props & ModalProps<UIArtist | null> = $props();

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

    <Button
      class="flex-1"
      variant="secondary"
      onclick={async () => {
        const res = await openInput({});
        if (res) {
          const artist = await apiClient.getArtistById(res);
          if (!artist.success) {
            if (artist.error.type === "ARTIST_NOT_FOUND") {
              toast.error("Failed to find an artist with that id");
            } else {
              toast.error("Unknown error");
            }

            console.error(formatError(artist.error));
            return;
          }

          close({
            id: artist.data.id,
            name: artist.data.name.default,
          });
        }
      }}
    >
      Use ID
    </Button>

    {#if currentQuery.length > 0}
      <Button
        class="line-clamp-1 flex-1 overflow-ellipsis"
        type="submit"
        variant="secondary"
        onclick={async () => {
          const res = await apiClient.createArtist({
            name: currentQuery,
            otherName: "",
          });
          if (!res.success) {
            // TODO(patrik): Should the toast include the error message?
            toast.error("Unknown error");
            throw formatError(res.error);
          }

          const artist = await apiClient.getArtistById(res.data.id);
          if (!artist.success) {
            // TODO(patrik): Should the toast include the error message?
            toast.error("Unknown error");
            throw formatError(artist.error);
          }

          close({ id: artist.data.id, name: artist.data.name.default });
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
