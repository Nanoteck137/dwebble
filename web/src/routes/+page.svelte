<script lang="ts">
  import { ApiClient } from "$lib/api/client.js";
  import ConfirmModal from "$lib/components/modals/ConfirmModal.svelte";
  import QueryArtistModal from "$lib/components/modals/QueryArtistModal.svelte";
  import { AlertDialog, Button, buttonVariants } from "@nanoteck137/nano-ui";
  import { modals } from "svelte-modals";

  const { data } = $props();

  async function test() {
    const confirmed = await modals.open(ConfirmModal, {
      title: "Are you sure?",
      description: "Hello World",
      confirmDelete: true,
    });
    console.log(confirmed);
  }
</script>

<p class="p-4 text-xl">Home Page</p>

<Button
  onclick={async () => {
    test();
  }}>Open Modal</Button
>

<Button
  onclick={async () => {
    const res = await modals.open(QueryArtistModal, {
      apiClient: new ApiClient(data.apiAddress),
    });
    console.log("artist", res);
  }}
>
  Open Artist Query
</Button>

<AlertDialog.Root>
  <AlertDialog.Trigger class={buttonVariants({ variant: "outline" })}>
    Show Dialog
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
      <AlertDialog.Description>
        This action cannot be undone. This will permanently delete your account
        and remove your data from our servers.
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
      <Button variant="destructive">Delete</Button>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
