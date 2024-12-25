<script lang="ts">
  import { goto } from "$app/navigation";
  import { createApiClient, openConfirm } from "$lib";
  import { cn } from "$lib/utils";
  import {
    Breadcrumb,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";
  import {
    Edit,
    EllipsisVertical,
    Import,
    Play,
    Trash,
    Wallpaper,
  } from "lucide-svelte";

  const { data } = $props();
  const apiClient = createApiClient(data);
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists">Artists</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/artists/{data.artist.id}">
          {data.artist.name.default}
        </Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Edit</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="flex gap-2">
  <div class="relative aspect-square w-48 min-w-48">
    <img
      class="inline-flex h-full w-full items-center justify-center rounded border object-cover text-xs"
      src={data.artist.picture.medium}
      alt="cover"
    />
    <div class="absolute inset-0 flex justify-end rounded p-1">
      <DropdownMenu.Root>
        <DropdownMenu.Trigger
          class={cn(
            buttonVariants({ variant: "ghost", size: "icon" }),
            "rounded-full",
          )}
        >
          <EllipsisVertical />
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="center">
          <DropdownMenu.Group>
            <DropdownMenu.Item
              onSelect={() => {
                goto(`/artists/${data.artist.id}/edit/details`);
              }}
            >
              <Edit />
              Edit Artist
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                goto(`/artists/${data.artist.id}/edit/picture`);
              }}
            >
              <Wallpaper />
              Change Picture
            </DropdownMenu.Item>
          </DropdownMenu.Group>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </div>
  </div>

  <div class="flex flex-col py-2">
    <p class="font-bold">
      {data.artist.name.default}
    </p>
    <!-- <p class="text-xs">
      Artist:
      <a class="hover:underline" href="/artists/{data.album.artistId}">
        {data.artist.artistName.default}
      </a>
    </p> -->

    {#if data.artist.name.other}
      <p class="text-xs">Other Name: {data.artist.name.other}</p>
    {/if}

    {#if data.artist.tags.length > 0}
      <p class="text-xs">Tags: {data.artist.tags.join(", ")}</p>
    {/if}

    <!-- {#if data.album.featuringArtists.length > 0}
      <p class="text-xs">Featuring Artists</p>
      {#each data.album.featuringArtists as artist}
        <p class="pl-2 text-xs">{artist.name.default}</p>
      {/each}
    {/if} -->
  </div>
</div>
