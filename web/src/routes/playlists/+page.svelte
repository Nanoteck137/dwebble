<script lang="ts">
  import type { PageData } from "./$types";

  interface Props {
    data: PageData;
  }

  let { data }: Props = $props();
</script>

<form method="post">
  <div class="flex flex-col">
    <input
      class="border-1 h-8 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-sm text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
      type="text"
      name="name"
      placeholder="Playlist Name"
    />
  </div>
  <button>Create Playlist</button>
</form>

<div class="flex flex-col gap-2">
  {#each data.playlists as playlist}
    <a href={`/playlists/${playlist.id}`}>{playlist.name}</a>
  {/each}
</div>
