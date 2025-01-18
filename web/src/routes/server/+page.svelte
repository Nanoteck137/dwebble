<script lang="ts">
  import { PUBLIC_COMMIT, PUBLIC_VERSION } from "$env/static/public";
  import { getApiClient, handleApiError } from "$lib";
  import { Button } from "@nanoteck137/nano-ui";

  const apiClient = getApiClient();

  let refillSearch = $state(false);
</script>

<p>Server Page (W.I.P)</p>

<p>Version: {PUBLIC_VERSION}</p>
<p>Commit: {PUBLIC_COMMIT}</p>

<Button
  onclick={async () => {
    refillSearch = true;
    const res = await apiClient.refillSearch();
    if (!res.success) {
      handleApiError(res.error);
    }
    refillSearch = false;
  }}
  disabled={refillSearch}
>
  Refill Search
</Button>
