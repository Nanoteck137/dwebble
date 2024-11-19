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
          <a class="relative" href={`/artists/${artist.id}`}>
            <img
              class="aspect-square w-40 rounded object-cover group-hover:brightness-75"
              src={artist.picture.medium}
              alt={`${artist.name} Picture`}
              loading="lazy"
            />

            <div
              class="absolute bottom-0 left-0 right-0 top-0 bg-purple-400/70"
            ></div>
          </a>
          <div class="h-2"></div>
          <a
            class="line-clamp-2 w-40 text-sm font-medium group-hover:underline"
            title={artist.name}
            href={`/artists/${artist.id}`}
          >
            {artist.name}
          </a>
        </div>
        <div class="h-2"></div>
      </div>
    {/each}
  </div>
</div>
