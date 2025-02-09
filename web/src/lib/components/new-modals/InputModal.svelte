<script lang="ts">
  import type { Modal } from "$lib/components/new-modals";
  import { Button, Dialog, Input } from "@nanoteck137/nano-ui";

  export type Props = {
    title?: string;
    placeholder?: string;
  };

  const {
    title,
    placeholder,

    class: className,
    children,
    onResult,
  }: Props & Modal<string> = $props();

  let open = $state(false);
  let input = $state("");
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
        onResult(input);
      }}
    >
      <Dialog.Header>
        <Dialog.Title>{title ?? "Enter input"}</Dialog.Title>
      </Dialog.Header>
      <Input {placeholder} bind:value={input} />
      <Dialog.Footer>
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>
        <Button type="submit">Ok</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
