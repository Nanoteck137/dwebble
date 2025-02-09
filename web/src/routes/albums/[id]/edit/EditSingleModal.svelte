<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/new-modals";
  import QueryArtistModal from "$lib/components/new-modals/QueryArtistModal.svelte";
  import type { CheckedValue, UIArtist } from "$lib/types";
  import {
    Button,
    buttonVariants,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { Plus, X } from "lucide-svelte";

  export type Props = {};

  export type Result = {
    name: CheckedValue<string>;
    otherName: CheckedValue<string>;
    year: CheckedValue<number>;
    tags: CheckedValue<string>;
    artist: CheckedValue<UIArtist | null>;
    featuringArtists: CheckedValue<UIArtist[]>;
  };

  const {
    class: className,
    children,
    onResult,
  }: Props & Modal<Result> = $props();

  let open = $state(false);
  let result = $state<Result>({
    name: {
      checked: false,
      value: "",
    },
    otherName: {
      checked: false,
      value: "",
    },
    year: {
      checked: false,
      value: 0,
    },
    tags: {
      checked: false,
      value: "",
    },
    artist: {
      checked: false,
      value: null,
    },
    featuringArtists: {
      checked: false,
      value: [],
    },
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger class={className}>
    {@render children?.()}
  </Dialog.Trigger>

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
        <Dialog.Title>Edit as Single</Dialog.Title>
      </Dialog.Header>

      <div class="flex flex-col gap-2">
        <FormItem>
          <Label for="name">Name</Label>

          <div class="flex items-center gap-2">
            <Checkbox class="h-5 w-5" bind:checked={result.name.checked} />
            <Input
              id="name"
              type="text"
              autocomplete="off"
              disabled={!result.name.checked}
              bind:value={result.name.value}
            />
          </div>
        </FormItem>

        <FormItem>
          <Label for="otherName">Other Name</Label>

          <div class="flex items-center gap-2">
            <Checkbox
              class="h-5 w-5"
              bind:checked={result.otherName.checked}
            />
            <Input
              id="otherName"
              type="text"
              autocomplete="off"
              disabled={!result.otherName.checked}
              bind:value={result.otherName.value}
            />
          </div>
        </FormItem>

        <FormItem>
          <Label for="year">Year</Label>

          <div class="flex max-w-32 items-center gap-2">
            <Checkbox class="h-5 w-5" bind:checked={result.year.checked} />
            <Input
              id="year"
              type="number"
              autocomplete="off"
              disabled={!result.year.checked}
              bind:value={result.year.value}
            />
          </div>
        </FormItem>

        <FormItem>
          <Label for="tags">Tags</Label>

          <div class="flex items-center gap-2">
            <Checkbox class="h-5 w-5" bind:checked={result.tags.checked} />
            <Input
              id="tags"
              type="text"
              autocomplete="off"
              disabled={!result.tags.checked}
              bind:value={result.tags.value}
            />
          </div>
        </FormItem>

        <FormItem>
          <Label>Artist</Label>
          <div class="flex items-center gap-2">
            <Checkbox class="h-5 w-5" bind:checked={result.artist.checked} />
            {#if result.artist.checked}
              <QueryArtistModal
                class={buttonVariants({ variant: "ghost", class: "w-fit" })}
                onResult={(artist) => {
                  result.artist.value = artist;
                }}
              >
                {result.artist.value?.name
                  ? `${result.artist.value?.name}`
                  : "Not Selected"}
              </QueryArtistModal>
            {:else}
              <Button class="w-fit" variant="ghost" disabled>
                {result.artist.value?.name
                  ? `${result.artist.value?.name}`
                  : "Not Selected"}
              </Button>
            {/if}
          </div>
        </FormItem>

        <Label class="flex items-center gap-2">
          Featuring Artists
          <!-- <button
            type="button"
            class="hover:cursor-pointer"
            tabindex={7}
            onclick={async () => {
              const res = await openArtistQuery({});
              if (res) {
                const index = featuringArtists.findIndex(
                  (a) => a.id === res.id,
                );
                if (index === -1) {
                  featuringArtists.push(res);
                }
              }
            }}
          >
          </button> -->

          {#if result.featuringArtists.checked}
            <QueryArtistModal
              onResult={(artist) => {
                const index = result.featuringArtists.value.findIndex(
                  (a) => a.id === artist.id,
                );
                if (index === -1) {
                  result.featuringArtists.value.push(artist);
                }
              }}
            >
              <Plus size="16" />
            </QueryArtistModal>
          {/if}
        </Label>

        <div class="flex flex-wrap gap-2">
          <Checkbox
            class="h-5 w-5"
            bind:checked={result.featuringArtists.checked}
          />
          {#each result.featuringArtists.value as artist, i}
            <p
              class="flex w-fit items-center gap-1 rounded-full bg-white px-2 py-1 text-xs text-black"
              title={`${artist.id}: ${artist.name}`}
            >
              <button
                type="button"
                class="text-red-400 hover:cursor-pointer"
                onclick={() => {
                  result.featuringArtists.value.splice(i, 1);
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
