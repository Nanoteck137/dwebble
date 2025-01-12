<script lang="ts">
  import { createApiClient } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();
  const apiClient = createApiClient(data);
</script>

<p class="p-4 text-xl">Home Page</p>

<Button
  onclick={async () => {
    const res = await apiClient.getQueue("dwebble-web-app");
    console.log(res);
  }}>Test</Button
>
