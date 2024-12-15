<script lang="ts">
  import type { Artist } from "$lib";
  import {
    Button,
    buttonVariants,
    Dialog,
    Input,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import type { FormEventHandler } from "svelte/elements";

  interface Props {
    open: boolean;
    currentQuery: string;
    queryResults: Artist[];

    onArtistSelected: (artist: Artist) => void;
    onInput: FormEventHandler<HTMLInputElement>;
  }

  let {
    open = $bindable(),
    currentQuery,
    queryResults,
    onArtistSelected,
    onInput,
  }: Props = $props();
</script>

<Dialog.Root bind:open>
  <Dialog.Trigger type="button" class={buttonVariants({ variant: "outline" })}>
    Change Artist
  </Dialog.Trigger>
  <Dialog.Content class="sm:max-w-[425px]">
    <Dialog.Header>
      <Dialog.Title>Change artist</Dialog.Title>
    </Dialog.Header>

    <Input oninput={onInput} />
    {#if currentQuery.length > 0}
      <!-- <form
        action="?/createArtist"
        method="post"
        use:enhance={() => {
          return async ({ update }) => {
            $open = false;
            $artist = undefined;
            await update();
          };
        }}
      >
        <input name="name" value={$currentQuery} type="hidden" />
        <Button type="submit" variant="secondary">
          New Artist: {$currentQuery}
        </Button>
      </form> -->
      <Button type="submit" variant="secondary">
        New Artist: {currentQuery}
      </Button>
    {/if}

    <ScrollArea class="max-h-36 overflow-y-clip">
      {#each queryResults as a}
        <Dialog.Close
          type="submit"
          class={buttonVariants({ variant: "ghost" })}
          onclick={() => {
            // $artist = a;
            onArtistSelected(a);
          }}
        >
          {a.name}
        </Dialog.Close>
      {/each}
    </ScrollArea>
  </Dialog.Content>
</Dialog.Root>
