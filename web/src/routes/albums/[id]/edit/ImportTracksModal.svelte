<script lang="ts">
  import type { Modal } from "$lib/components/new-modals";
  import type { UIArtist } from "$lib/types";
  import { Button, Dialog, ScrollArea, Separator } from "@nanoteck137/nano-ui";
  import { X } from "lucide-svelte";

  export type Props = {
    open: boolean;
  };

  export type Result = {
    file: File[];
  };

  let { open = $bindable(), onResult }: Props & Modal<Result> = $props();

  let fileSelector = $state<HTMLInputElement>();

  let result = $state<Result>({
    file: [],
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

      <ScrollArea class="max-h-[280px]">
        {#each result.file as file, i}
          <div class="flex items-center gap-2 py-2">
            <p>{i}. {file.name}</p>
            <button
              type="button"
              onclick={() => {
                result.file.splice(i, 1);
              }}
            >
              <X size={20} />
            </button>
          </div>
          <Separator />
        {/each}
      </ScrollArea>

      <Button
        variant="outline"
        onclick={() => {
          fileSelector?.click();
        }}
      >
        Open File Selector
      </Button>

      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button type="submit">Upload</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>

<input
  class="hidden"
  bind:this={fileSelector}
  type="file"
  multiple
  accept="audio/*"
  onchange={(e) => {
    const input = e.target as HTMLInputElement;
    if (!input.files) return;

    const files = input.files;
    for (let i = 0; i < files.length; i++) {
      const file = files.item(i);
      if (file) {
        result.file.push(file);
      }
    }
  }}
/>
