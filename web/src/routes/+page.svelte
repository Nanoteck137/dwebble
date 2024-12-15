<script lang="ts">
  import { ApiClient } from "$lib/api/client";
  import { getModalState } from "$lib/modal.svelte";
  import { AlertDialog, Button, buttonVariants } from "@nanoteck137/nano-ui";

  const modalState = getModalState();

  const { data } = $props();
</script>

<p class="p-4 text-xl">Home Page</p>

<Button
  onclick={() => {
    modalState.pushModal({
      type: "modal-confirm",
      title: "Are you sure?",
      description: "You are about to delete this",
      confirmDelete: true,
      onConfirm: () => {
        modalState.pushModal({
          type: "modal-confirm",
          title: "Testing",
        });
      },
    });
  }}
>
  Open Modal
</Button>

<Button
  onclick={() => {
    modalState.pushModal({
      type: "modal-query-artist",
      apiClient: new ApiClient(data.apiAddress),
      onArtistSelected: (artist) => {
        console.log(artist);
      },
    });
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
      <AlertDialog.Action>Continue</AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
