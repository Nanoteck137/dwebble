<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { Track } from "$lib/api/types";
  import { Button } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";
  import { toast } from "svelte-5-french-toast";

  type Props = {
    show: boolean;
    track: Track;
    isInQuickPlaylist: (trackId: string) => boolean;
  };

  const { show, track, isInQuickPlaylist }: Props = $props();
  const apiClient = getApiClient();
</script>

{#if show}
  <Button
    type="submit"
    class="rounded-full"
    variant="ghost"
    size="icon-lg"
    onclick={async () => {
      if (isInQuickPlaylist(track.id)) {
        const res = await apiClient.removeItemFromUserQuickPlaylist({
          trackId: track.id,
        });

        if (!res.success) {
          handleApiError(res.error);
          return;
        }
      } else {
        const res = await apiClient.addToUserQuickPlaylist({
          trackId: track.id,
        });

        if (!res.success) {
          handleApiError(res.error);
          return;
        }
      }

      await invalidateAll();
    }}
  >
    {#if isInQuickPlaylist(track.id)}
      <Star class="fill-primary" />
    {:else}
      <Star />
    {/if}
  </Button>
{/if}
