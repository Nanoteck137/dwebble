<script>
  const { data } = $props();
</script>

<p>Change Album Cover</p>

<form method="post" enctype="multipart/form-data">
  <input name="albumId" type="hidden" value={data.album.id} />

  <input name="file" type="file" accept="image/png, image/jpeg" />
  <button>Save</button>
</form>
