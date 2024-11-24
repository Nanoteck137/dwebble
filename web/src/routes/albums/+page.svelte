<script lang="ts">
  import { Button, Select } from "@nanoteck137/nano-ui";
  import type { PageData } from "./$types";
  import { enhance } from "$app/forms";

  interface Props {
    data: PageData;
  }

  let { data }: Props = $props();

  const sorts = [
    { value: "sort=album", label: "Name (A-Z)" },
    { value: "sort=-album", label: "Name (Z-A)" },
    { value: "sort=-created", label: "Newest First" },
    { value: "sort=created", label: "Newest Last" },
  ];

  let value = $state("");

  const triggerContent = $derived(
    sorts.find((f) => f.value === value)?.label ?? "Sort",
  );

  let form: HTMLFormElement | undefined = $state();
</script>

<div class="flex flex-col gap-2">
  <form bind:this={form} method="get">
    <Select.Root type="single" name="sort" bind:value>
      <Select.Trigger class="w-[180px]">
        {triggerContent}
      </Select.Trigger>
      <Select.Content>
        <Select.Group>
          <Select.GroupHeading>Sort</Select.GroupHeading>
          {#each sorts as sort}
            <Select.Item value={sort.value} label={sort.label} />
          {/each}
        </Select.Group>
      </Select.Content>
    </Select.Root>

    <Button type="submit">Filter</Button>
  </form>

  <div class="flex items-center justify-between">
    <p class="text-bold text-xl">Albums</p>
    <p class="text-sm">{data.albums.length} albums(s)</p>
  </div>
  <div
    class="grid grid-cols-2 gap-2 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 2xl:grid-cols-7"
  >
    {#each data.albums as album}
      <div class="flex flex-col items-center">
        <div class="group">
          <a class="relative" href={`/albums/${album.id}`}>
            <img
              class="aspect-square w-40 rounded object-cover group-hover:brightness-75"
              src={album.coverArt.medium}
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
        <div class="h-2"></div>
      </div>
    {/each}
  </div>
</div>
