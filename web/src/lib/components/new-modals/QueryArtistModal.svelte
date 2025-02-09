<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import type { Modal } from "$lib/components/new-modals";
  import InputModal from "$lib/components/new-modals/InputModal.svelte";
  import type { QueryArtist, UIArtist } from "$lib/types";
  import {
    Button,
    buttonVariants,
    Dialog,
    Input,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";

  export type Props = {
    title?: string;
  };

  const {
    title,

    class: className,
    children,
    onResult,
  }: Props & Modal<UIArtist> = $props();

  const apiClient = getApiClient();

  let open = $state(false);

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
        handleApiError(res.error);
      }
    }, 500);
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger class={className}>
    {@render children?.()}
  </Dialog.Trigger>

  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>{title ?? "Search Artist"}</Dialog.Title>
    </Dialog.Header>

    <Input oninput={onInput} placeholder="Search..." />

    <InputModal
      class={buttonVariants({ variant: "secondary" })}
      onResult={async (value) => {
        const artist = await apiClient.getArtistById(value);
        if (!artist.success) {
          if (artist.error.type === "ARTIST_NOT_FOUND") {
            toast.error("No artist with id");
          } else {
            handleApiError(artist.error);
          }

          return;
        }

        onResult({
          id: artist.data.id,
          name: artist.data.name.default,
        });
        open = false;
      }}
    >
      Use ID
    </InputModal>

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
            handleApiError(res.error);
            return;
          }

          const artist = await apiClient.getArtistById(res.data.id);
          if (!artist.success) {
            handleApiError(artist.error);
            return;
          }

          onResult({
            id: artist.data.id,
            name: artist.data.name.default,
          });
          open = false;
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
              onResult(result);
              open = false;
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
          open = false;
        }}
      >
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
