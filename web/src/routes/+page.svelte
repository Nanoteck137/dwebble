<script lang="ts">
  import { createApiClient } from "$lib";
  import { getMusicManager } from "$lib/music-manager";
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = createApiClient(data);

  const musicManager = getMusicManager();
</script>

<p class="p-4 text-xl">Home Page</p>

<Button
  onclick={async () => {
    const queue = await apiClient.getDefaultQueue("dwebble-web-app");
    if (!queue.success) {
      throw queue.error;
    }

    const queueId = queue.data.id;

    // console.log("Clear Queue", await apiClient.clearQueue(queueId));
    // console.log(
    //   "Add album items",
    //   await apiClient.addToQueueFromAlbum(queueId, "dwijd8a0o114mgof"),
    // );

    const queueItems = await apiClient.getQueueItems(queueId);
    console.log("Get queue items", queueItems);

    if (queueItems.success) {
      musicManager.setQueueItems(
        queueItems.data.items.map((i) => i.track),
        queueItems.data.index,
      );
      musicManager.requestPlay();
    }
  }}
>
  Test
</Button>
