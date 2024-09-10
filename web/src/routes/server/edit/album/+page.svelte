<script lang="ts">
  const { data } = $props();
</script>

<div
  class="grid grid-cols-2 gap-4 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-7 2xl:grid-cols-8"
>
  {#each data.albums as album}
    <div class="flex flex-col items-center">
      <div class="group">
        <a class="relative" href={`/server/edit/album/${album.id}`}>
          <img
            class="aspect-square w-40 rounded object-cover group-hover:brightness-75"
            src={album.coverArt}
            alt={`${album.name} Cover Art`}
            loading="lazy"
          />

          <div
            class="absolute bottom-0 left-0 right-0 top-0 bg-purple-400/70"
          ></div>
        </a>
        <div class="h-2"></div>
        <a
          class="line-clamp-2 w-40 text-sm font-medium group-hover:underline"
          title={album.name}
          href={`/albums/${album.id}`}
        >
          {album.name}
        </a>
      </div>
      <a
        class="line-clamp-1 w-40 text-xs hover:underline"
        title={album.artistName}
        href={`/artist/${album.artistId}`}
      >
        {album.artistName}
      </a>
    </div>
  {/each}
</div>
