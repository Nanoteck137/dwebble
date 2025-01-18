<script lang="ts">
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { Check, EllipsisVertical, FileHeart } from "lucide-svelte";
  import { cn } from "$lib/utils";
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";

  let { data } = $props();
  const apiClient = getApiClient();
</script>

<div class="flex flex-col gap-4">
  <div>
    <Button href="/playlists/new">New Playlist</Button>
  </div>

  <div class="flex flex-col gap-2">
    {#each data.playlists as playlist}
      <div class="flex items-center justify-between px-2">
        <a
          class="flex flex-grow items-center gap-2 py-2 hover:underline"
          href="/playlists/{playlist.id}"
        >
          {#if data.user?.quickPlaylist === playlist.id}
            <Check class="min-h-4 min-w-4" size={16} />
          {:else}
            <div class="min-h-4 min-w-4"></div>
          {/if}
          {playlist.name}
        </a>

        <div>
          <DropdownMenu.Root>
            <DropdownMenu.Trigger
              class={cn(
                buttonVariants({ variant: "ghost", size: "icon-lg" }),
                "rounded-full",
              )}
            >
              <EllipsisVertical />
            </DropdownMenu.Trigger>
            <DropdownMenu.Content align="end">
              <DropdownMenu.Group>
                {#if data.user?.quickPlaylist !== playlist.id}
                  <DropdownMenu.Item
                    onSelect={async () => {
                      const res = await apiClient.updateUserSettings({
                        quickPlaylist: playlist.id,
                      });
                      if (!res.success) {
                        handleApiError(res.error);
                        return;
                      }

                      await invalidateAll();
                    }}
                  >
                    <FileHeart />
                    Set as Quick Playlist
                  </DropdownMenu.Item>
                {/if}
              </DropdownMenu.Group>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        </div>
      </div>
      <Separator />
    {/each}
  </div>
</div>
