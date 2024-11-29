<script lang="ts">
  import type { PageData } from "./$types";

  interface Props {
    data: PageData;
  }

  let { data }: Props = $props();
</script>

<div class="flex flex-col">
  <p>Num Artists: {data.artists.length}</p>
  <div
    class="grid grid-cols-2 gap-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
  >
    {#each data.artists as artist}
      <div class="flex flex-col items-center">
        <div class="group">
          <a href="/artists/${artist.id}">
            <!-- svelte-ignore a11y_img_redundant_alt -->
            <img
              class="inline-flex aspect-square w-40 min-w-40 items-center justify-center rounded border object-cover text-xs group-hover:brightness-75"
              src={artist.picture.medium}
              alt="picture"
              loading="lazy"
            />
          </a>
          <div class="h-2"></div>
          <a
            class="line-clamp-2 w-40 text-sm font-medium group-hover:underline"
            title={artist.name}
            href="/artists/{artist.id}"
          >
            {artist.name}
          </a>
        </div>
        <div class="h-2"></div>
      </div>
    {/each}
  </div>
</div>
