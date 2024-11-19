<script lang="ts">
  const { data } = $props();
</script>

<a href="cover">Change Cover</a>

<form method="post">
  <div class="flex flex-col gap-1">
    <label for="name">Album Name</label>
    <input
      class="w-full rounded bg-[--bg-color] text-xs"
      id="name"
      name="name"
      value={data.album.name}
      type="text"
    />
  </div>

  <div class="flex flex-col gap-1">
    <label for="artist">Album Artist</label>
    <input
      class="w-full rounded bg-[--bg-color] text-xs"
      id="artist"
      name="artist"
      value={data.album.artistName}
      type="text"
    />
  </div>

  <button>Save</button>
  <a href="/albums/{data.album.id}/edit">Back</a>
</form>
