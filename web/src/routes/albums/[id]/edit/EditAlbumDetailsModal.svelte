<script lang="ts">
  import type { Album } from "$lib/api/types";
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/new-modals";
  import QueryArtistModal from "$lib/components/new-modals/QueryArtistModal.svelte";
  import type { UIArtist } from "$lib/types";
  import {
    Button,
    buttonVariants,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { Plus, RotateCcw, X } from "lucide-svelte";

  export type Props = {
    album: Album;
    open: boolean;
  };

  export type Result = {
    name: string;
    otherName: string | null;
    year: number | null;
    tags: string;
    artist: UIArtist;
    featuringArtists: UIArtist[];
  };

  let {
    album,
    open = $bindable(),
    onResult,
  }: Props & Modal<Result> = $props();

  let result = $state<Result>({
    name: "",
    otherName: null,
    year: null,
    tags: "",
    artist: {
      id: "",
      name: "",
    },
    featuringArtists: [],
  });

  $effect(() => {
    if (open) {
      result = {
        name: album.name.default,
        otherName: album.name.other,
        year: album.year,
        tags: album.tags.join(","),
        artist: {
          id: album.artistId,
          name: album.artistName.default,
        },
        featuringArtists: album.featuringArtists.map((a) => ({
          id: a.id,
          name: a.name.default,
        })),
      };
    }
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();
        onResult(result);
        open = false;
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Edit album details</Dialog.Title>
      </Dialog.Header>

      <div class="flex flex-col gap-4">
        <FormItem>
          <Label for="name">Name</Label>
          <Input
            id="name"
            type="text"
            autocomplete="off"
            bind:value={result.name}
          />
        </FormItem>

        <FormItem>
          <Label for="otherName">Other Name</Label>
          <Input
            id="otherName"
            type="text"
            autocomplete="off"
            bind:value={result.otherName}
          />
        </FormItem>

        <FormItem>
          <Label for="year">Year</Label>
          <Input
            class="w-24"
            id="year"
            type="number"
            autocomplete="off"
            bind:value={result.year}
          />
        </FormItem>

        <FormItem>
          <Label for="tags">Tags</Label>
          <Input
            id="tags"
            type="text"
            autocomplete="off"
            bind:value={result.tags}
          />
        </FormItem>

        <FormItem>
          <Label>Artist</Label>
          <QueryArtistModal
            class={buttonVariants({ variant: "ghost", class: "w-fit" })}
            onResult={(artist) => {
              result.artist = artist;
            }}
          >
            {result.artist?.name ? `${result.artist?.name}` : "Not Selected"}
          </QueryArtistModal>
        </FormItem>

        <Label class="flex items-center gap-2">
          Featuring Artists
          <QueryArtistModal
            onResult={(artist) => {
              const index = result.featuringArtists.findIndex(
                (a) => a.id === artist.id,
              );
              if (index === -1) {
                result.featuringArtists.push(artist);
              }
            }}
          >
            <Plus size="16" />
          </QueryArtistModal>
        </Label>

        <div class="flex flex-wrap gap-2">
          {#each result.featuringArtists as artist, i}
            <p
              class="flex w-fit items-center gap-1 rounded-full bg-white px-2 py-1 text-xs text-black"
              title={`${artist.id}: ${artist.name}`}
            >
              <button
                type="button"
                class="text-red-400 hover:cursor-pointer"
                onclick={() => {
                  result.featuringArtists.splice(i, 1);
                }}
              >
                <X size="16" />
              </button>
              {artist.name}
            </p>
          {/each}
        </div>
      </div>

      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button type="submit">Save</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
