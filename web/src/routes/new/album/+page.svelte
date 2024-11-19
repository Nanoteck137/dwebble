<script lang="ts">
  import type { ActionData, PageData } from "./$types";

  const { form }: { data: PageData; form: ActionData } = $props();
</script>

<div class="px-4">
  <form
    class="flex flex-col gap-2"
    method="post"
    enctype="multipart/form-data"
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

    <div class="h-2"></div>

    <button
      class="rounded bg-blue-400 px-4 py-2 text-black hover:bg-blue-500 active:scale-95"
    >
      Create Album
    </button>
  </form>
</div>
