<script lang="ts">
  import { Button, Select } from "@nanoteck137/nano-ui";
  import type { PageData } from "./$types";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import { isRoleAdmin } from "$lib/utils";

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

{#if isRoleAdmin(data.user?.role || "")}
  <Button href="/albums/new">New Album</Button>
{/if}

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
          <a href="/albums/{album.id}">
            <img
              class="inline-flex aspect-square w-40 min-w-40 items-center justify-center rounded border object-cover text-xs group-hover:brightness-75"
              src={album.coverArt.medium}
              alt="cover"
              loading="lazy"
            />
          </a>
          <div class="h-2"></div>
          <a
            class="line-clamp-2 w-40 text-sm font-medium group-hover:underline"
            title={album.name.default}
            href="/albums/{album.id}"
          >
            {album.name.default}
          </a>
        </div>
        <ArtistList artists={album.allArtists} />
        <div class="h-2"></div>
      </div>
    {/each}
  </div>
</div>
