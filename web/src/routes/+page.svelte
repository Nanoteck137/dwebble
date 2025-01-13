<script lang="ts">
  import { createApiClient } from "$lib";
  import { getMusicManager } from "$lib/music-manager.svelte.js";
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = createApiClient(data);

  const musicManager = getMusicManager();
</script>

<p class="p-4 text-xl">Home Page</p>

<Button
  onclick={async () => {
    // console.log("Clear Queue", await apiClient.clearQueue(queueId));
    console.log(
      "Add album items",
      await apiClient.addToQueueFromAlbum(data.queueId, "ze8yir8qw5e12lx0"),
    );
  }}
>
  Test
</Button>
