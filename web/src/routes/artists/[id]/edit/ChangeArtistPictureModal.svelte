<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/new-modals";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";

  export type Props = {
    open: boolean;
  };

  let { open = $bindable(), onResult }: Props & Modal<FormData> = $props();
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <form
      class="flex flex-col gap-4"
      onsubmit={(e) => {
        e.preventDefault();

        const formData = new FormData(e.target as HTMLFormElement);
        onResult(formData);
        open = false;
      }}
    >
      <Dialog.Header>
        <Dialog.Title>Change Artist Picture</Dialog.Title>
      </Dialog.Header>

      <div class="flex flex-col gap-4">
        <FormItem>
          <Label for="picture">Picture</Label>
          <Input
            id="picture"
            name="picture"
            type="file"
            accept="image/png,image/jpeg"
          />
        </FormItem>
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
