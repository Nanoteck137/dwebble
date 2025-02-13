<script lang="ts">
  import { getQuickPlaylist } from "$lib/quick-playlist.svelte";
  import { Button } from "@nanoteck137/nano-ui";
  import { Star } from "lucide-svelte";

  type Props = {
    show: boolean;
    trackId: string;
  };

  const { show, trackId }: Props = $props();
  const quickPlaylist = getQuickPlaylist();
</script>

{#if show}
  <Button
    type="submit"
    class="rounded-full"
    variant="ghost"
    size="icon-lg"
    onclick={() => {
      quickPlaylist.toggleTrack(trackId);
    }}
  >
    {#if quickPlaylist.hasTrack(trackId)}
      <Star class="fill-primary" />
    {:else}
      <Star />
    {/if}
  </Button>
{/if}
