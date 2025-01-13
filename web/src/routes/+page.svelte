<script lang="ts">
  import { createApiClient } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = createApiClient(data);
</script>

<p class="p-4 text-xl">Home Page</p>

<Button
  onclick={async () => {
    const queue = await apiClient.getDefaultQueue("dwebble-web-app");
    if (!queue.success) {
      throw queue.error;
    }

    const queueId = queue.data.id;

    console.log("Clear Queue", await apiClient.clearQueue(queueId));
    console.log(
      "Add album items",
      await apiClient.addToQueueFromAlbum(queueId, "njovbcxihdejop5q"),
    );

    console.log("Get queue items", await apiClient.getQueueItems(queueId));

    // const res = await apiClient.getQueue();
    // console.log(res);
  }}>Test</Button
>
