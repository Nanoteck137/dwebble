<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import type { Playlist } from "$lib/api/types";
  import type { Modal } from "$lib/components/new-modals";
  import Spinner from "$lib/components/Spinner.svelte";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import { Check } from "lucide-svelte";
  import { onMount } from "svelte";

  export type Props = {
    playlists: Playlist[];
    currentQuickPlaylistId: string | undefined;
  };

  const {
    playlists,
    currentQuickPlaylistId,

    class: className,
    children,
    onResult,
  }: Props & Modal<string> = $props();

  let open = $state(false);
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger class={className}>
    {@render children?.()}
  </Dialog.Trigger>

  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Select Playlist</Dialog.Title>
    </Dialog.Header>

    <ScrollArea class="max-h-[280px]">
      <div class="flex flex-col gap-2">
        {#each playlists as playlist}
          <Button
            class="justify-start"
            variant="ghost"
            disabled={currentQuickPlaylistId === playlist.id}
            onclick={() => {
              onResult(playlist.id);
              open = false;
            }}
          >
            {#if currentQuickPlaylistId === playlist.id}
              <Check class="min-h-4 min-w-4" size={16} />
            {:else}
              <div class="min-h-4 min-w-4"></div>
            {/if}
            {playlist.name}
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
