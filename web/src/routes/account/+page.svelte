<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  let { data } = $props();
</script>

<p>{data.user?.username}</p>
<p>{data.user?.quickPlaylist}</p>

<Button href="/account/password">Change Password</Button>
<Button href="/account/quickplaylist">Change Quick Playlist</Button>
