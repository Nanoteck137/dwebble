<script lang="ts">
  import type { ApiClient } from "$lib/api/client";
  import type { Playlist, Track } from "$lib/api/types";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import type { ModalProps } from "svelte-modals";

  export type Props = {
    apiClient: ApiClient;
    track: Track;
    playlists: Playlist[];
  };

  const { apiClient, track, playlists, isOpen, close }: Props & ModalProps =
    $props();
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
  <Dialog.Content class="flex flex-col gap-4">
    <Dialog.Header>
      <Dialog.Title>Save track to playlist</Dialog.Title>
    </Dialog.Header>

    <ScrollArea class="max-h-36 overflow-y-clip">
      <div class="flex flex-col">
        {#each playlists as playlist, i}
          <Button
            variant="ghost"
            onclick={async () => {
              console.log(track);
              const res = await apiClient.addItemToPlaylist(playlist.id, {
                trackId: track.id,
              });
              if (!res.success) {
                // TODO(patrik): Toast
                console.log(res.error.message);

                close();
                return;
              }

              // TODO(patrik): Toast
              close();
            }}
          >
            {playlist.name}
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
