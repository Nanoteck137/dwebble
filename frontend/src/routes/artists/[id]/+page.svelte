<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<div class="p-2">
  <p>{data.artist.name}</p>

  <div class="flex flex-wrap">
    {#each data.artist.albums as album}
      <a
        href={`/albums/${album.id}`}
        class="border shadow-lg hover:cursor-pointer hover:shadow-xl active:scale-95"
      >
        <img
          class="pointer-events-none h-40 w-40"
          src="https://placehold.co/600x600.png?text=Placeholder"
          alt=""
        />
        <p class="py-2 text-center">{album.name}</p>
      </a>
    {/each}
  </div>
</div>
