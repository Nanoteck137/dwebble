<script lang="ts">
  import { enhance } from "$app/forms";

  let files = $state<File[]>([]);
  let fileSelector = $state<HTMLInputElement>();

  function submit(formData: FormData) {}
</script>

<div class="px-4">
  <form
    action="?/import"
    method="post"
    enctype="multipart/form-data"
    use:enhance={({ formData }) => {
      files.forEach((f) => {
        formData.append("files", f);
      });
    }}
  >
    <div class="flex flex-col gap-1">
      <label class="text-lg tracking-wide" for="albumName">Album Name</label>
      <input
        class="border-1 h-12 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
        id="albumName"
        name="albumName"
        type="text"
        placeholder="Album Name"
      />
    </div>

    <div class="flex flex-col gap-1">
      <label class="text-lg tracking-wide" for="artistName">Artist Name</label>
      <input
        class="border-1 h-12 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
        id="artistName"
        name="artistName"
        type="text"
        placeholder="Artist Name"
      />
    </div>

    {#each files as file}
      <p>{file.name}</p>
    {/each}

    <button
      class="rounded bg-blue-400 px-4 py-2 text-black hover:bg-blue-500 active:scale-95"
      type="button"
      onclick={() => {
        fileSelector?.click();
      }}>Add Files</button
    >

    <button
      class="rounded bg-blue-400 px-4 py-2 text-black hover:bg-blue-500 active:scale-95"
      >Import</button
    >
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
