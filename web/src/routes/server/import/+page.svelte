<script lang="ts">
  import { enhance } from "$app/forms";
  import { X } from "lucide-svelte";
  import type { ActionData, PageData } from "./$types";

  const { data, form }: { data: PageData; form: ActionData } = $props();

  $effect(() => {
    console.log("Form", form);
  });

  let files = $state<File[]>([]);
  let fileSelector = $state<HTMLInputElement>();

  let coverArt = $state<File>();
  let coverArtSelector = $state<HTMLInputElement>();

  let tags = $state<string[]>([]);
</script>

<div class="px-4">
  <form
    class="flex flex-col gap-2"
    action="?/import"
    method="post"
    enctype="multipart/form-data"
    use:enhance={({ formData }) => {
      if (coverArt) {
        formData.set("coverArt", coverArt);
      }

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
      {#if form?.errors?.albumName}
        {#each form?.errors?.albumName as error}
          <p class="text-red-400">{error}</p>
        {/each}
      {/if}
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
      {#if form?.errors?.artistName}
        {#each form?.errors?.artistName as error}
          <p class="text-red-400">{error}</p>
        {/each}
      {/if}
    </div>

    <p>Cover Art: {coverArt?.name}</p>

    <button
      class="rounded bg-blue-400 px-4 py-2 text-black hover:bg-blue-500 active:scale-95"
      type="button"
      onclick={() => {
        coverArtSelector?.click();
      }}>Set Cover Art</button
    >

    <p>Global Tags</p>
    <div class="flex flex-wrap">
      {#each tags as tag, i}
        <div
          class="flex items-center gap-1 overflow-hidden rounded bg-purple-400 px-2 py-1"
        >
          <button
            type="button"
            onclick={() => {
              tags.splice(i, 1);
            }}
          >
            <X size="25" />
          </button>
          <p class="text-ellipsis text-sm">{tag}</p>
        </div>
      {/each}
    </div>

    <button
      type="button"
      onclick={() => {
        const newTag = prompt("Tag Name");
        if (!newTag) {
          return;
        }

        for (let i = 0; i < tags.length; i++) {
          if (tags[i].toLowerCase() === newTag.toLowerCase()) {
            return;
          }
        }

        tags.push(newTag);
      }}>Add Tag</button
    >

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
  bind:this={coverArtSelector}
  type="file"
  accept="image/png, image/jpeg"
  onchange={(e) => {
    const target = e.target as HTMLInputElement;

    if (!target.files) {
      return;
    }

    coverArt = target.files[0];
  }}
/>

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
