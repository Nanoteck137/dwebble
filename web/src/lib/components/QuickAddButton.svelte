<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import type { ApiClient } from "$lib/api/client";
  import type { Track } from "$lib/api/types";
  import { formatError } from "$lib/utils";
  import { Button } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";
  import { toast } from "svelte-5-french-toast";

  type Props = {
    show: boolean;
    track: Track;
    apiClient: ApiClient;
    isInQuickPlaylist: (trackId: string) => boolean;
  };

  const { show, track, apiClient, isInQuickPlaylist }: Props = $props();
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
          toast.error("Unknown error");
          throw formatError(res.error);
        }
      } else {
        const res = await apiClient.addToUserQuickPlaylist({
          trackId: track.id,
        });

        if (!res.success) {
          toast.error("Unknown error");
          throw formatError(res.error);
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
