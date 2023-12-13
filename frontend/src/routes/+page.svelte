<script lang="ts">
  import type { PageData } from "./$types";

  export let data: PageData;
</script>

<!-- <div class="flex flex-wrap">
  {#each data.artists as artist}
    <div class="h-64 w-64 bg-red-400">
      <p>{artist.name}</p>
    </div>
  {/each}
</div> -->

<div class="grid grid-cols-2 gap-2 p-2">
  {#each data.artists as artist (artist.id)}
    <a
      class="relative flex aspect-square flex-col overflow-clip rounded border-2 bg-white shadow-md transition-transform hover:shadow-xl active:scale-95"
      href={`/artists/${artist.id}`}
    >
      <div class="relative h-full w-full">
        <img
          class="pointer-events-none absolute flex border-b-2 object-cover"
          src="https://placehold.co/600x600.png?text=Placeholder"
          crossorigin="anonymous"
          alt=""
        />
      </div>

      <p class="w-full p-4 text-center">{artist.name}</p>
    </a>
  {/each}
</div>
