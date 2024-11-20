<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Label } from "$lib/components/ui/label";

  const { data } = $props();
</script>

<form class="flex flex-col gap-2" method="post">
  <div class="flex flex-col gap-1">
    <Label for="name">Album Name</Label>
    <Input
      class="w-full rounded bg-[--bg-color] text-xs"
      id="name"
      name="name"
      value={data.album.name}
      type="text"
    />
  </div>

  <div class="flex flex-col gap-1">
    <Label for="artist">Album Artist</Label>
    <Input
      class="w-full rounded bg-[--bg-color] text-xs"
      id="artist"
      name="artist"
      value={data.album.artistName}
      type="text"
    />
  </div>

  <div class="flex justify-end gap-2">
    <Button variant="ghost" href="/albums/{data.album.id}/edit">Back</Button>
    <Button type="submit">Save</Button>
  </div>
</form>
