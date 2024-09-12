<script lang="ts">
  import { enhance } from "$app/forms";

  const { data } = $props();

  let files = $state<File[]>([]);
  let fileSelector = $state<HTMLInputElement>();
</script>

<div class="px-4">
  <form
    class="flex flex-col gap-2"
    method="post"
    enctype="multipart/form-data"
    use:enhance={({ formData }) => {
      files.forEach((f) => {
        formData.append("files", f);
      });
    }}
  >
    <input name="albumId" value={data.album.id} type="hidden" />

    {#each files as file}
      <p>{file.name}</p>
    {/each}

    <button
      class="rounded bg-blue-400 px-4 py-2 text-black hover:bg-blue-500 active:scale-95"
      type="button"
      onclick={() => {
        fileSelector?.click();
      }}
    >
      Add Files
    </button>

    <button
      class="rounded bg-blue-400 px-4 py-2 text-black hover:bg-blue-500 active:scale-95"
    >
      Import
    </button>
  </form>
</div>

<input
  class="hidden"
  bind:this={fileSelector}
  type="file"
  multiple
  accept="audio/*"
  onchange={(e) => {
    console.log(e);
    const target = e.target as HTMLInputElement;

    if (!target.files) {
      return;
    }

    for (let i = 0; i < target.files.length; i++) {
      const item = target.files.item(i);
      if (!item) {
        continue;
      }

      files.push(item);
    }

    files.sort((a, b) => a.name.localeCompare(b.name));
  }}
/>
