<script lang="ts">
  import { goto } from "$app/navigation";
  import type { ArtistInfo } from "$lib/api/types";
  import ArtistList from "$lib/components/ArtistList.svelte";
  import Image from "$lib/components/Image.svelte";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Play, Shuffle } from "lucide-svelte";
  import type { Snippet } from "svelte";

  export type Props = {
    name: string;

    image?: string;
    tags?: string[];
    artists?: ArtistInfo[];

    more?: Snippet<[]>;

    onPlay?: (shuffle: boolean) => void;
  };

  const { name, image, tags, artists, more, onPlay }: Props = $props();
</script>

<div class="flex h-48">
  {#if image}
    <Image class="w-48 min-w-48" src={image} alt="cover" />
    <div class="w-2"></div>
  {/if}

  <div class="flex flex-col">
    <div class="flex flex-col">
      <p class="font-bold">
        {name}
      </p>

      {#if artists}
        <ArtistList {artists} />
      {/if}

      {#if tags}
        <p class="text-xs text-muted-foreground">{tags.join(", ")}</p>
      {/if}
    </div>

    <div class="flex-grow"></div>

    <div class="flex gap-2">
      <Button
        variant="outline"
        onclick={() => {
          onPlay?.(false);
        }}
      >
        <Play />
        Play
      </Button>

      <Button
        variant="outline"
        onclick={() => {
          onPlay?.(true);
        }}
      >
        <Shuffle />
        Shuffle
      </Button>

      {#if more}
        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={buttonVariants({ variant: "outline", size: "icon" })}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content align="start">
            {@render more()}
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      {/if}
    </div>
  </div>
</div>
