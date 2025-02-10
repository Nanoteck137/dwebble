<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError, openConfirm } from "$lib";
  import { cn } from "$lib/utils";
  import {
    Breadcrumb,
    buttonVariants,
    DropdownMenu,
  } from "@nanoteck137/nano-ui";
  import { Edit, EllipsisVertical, Trash, Wallpaper } from "lucide-svelte";
  import toast from "svelte-5-french-toast";
  import Image from "$lib/components/Image.svelte";
  import ChangeArtistPictureModal from "./ChangeArtistPictureModal.svelte";
  import EditArtistDetailsModal from "./EditArtistDetailsModal.svelte";

  const { data } = $props();
  const apiClient = getApiClient();

  let openEditArtistDetailsModal = $state(false);
  let openArtistPictureModal = $state(false);
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
    <Image
      class="h-full w-full"
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
                openEditArtistDetailsModal = true;
              }}
            >
              <Edit />
              Edit Artist
            </DropdownMenu.Item>

            <DropdownMenu.Item
              onSelect={() => {
                openArtistPictureModal = true;
              }}
            >
              <Wallpaper />
              Change Picture
            </DropdownMenu.Item>

            <DropdownMenu.Separator />

            <DropdownMenu.Item
              onSelect={async () => {
                const confirmed = await openConfirm({
                  title: "Are you sure?",
                  description: "You are about to delete this artist",
                });
                if (confirmed) {
                  const res = await apiClient.deleteArtist(data.artist.id);
                  if (!res.success) {
                    handleApiError(res.error);
                    return;
                  }

                  toast.success("Deleted artist");
                  goto("/artists", { invalidateAll: true });
                }
              }}
            >
              <Trash />
              Delete Artist
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

    {#if data.artist.name.other}
      <p class="text-xs">Other Name: {data.artist.name.other}</p>
    {/if}

    {#if data.artist.tags.length > 0}
      <p class="text-xs">Tags: {data.artist.tags.join(", ")}</p>
    {/if}
  </div>
</div>

<EditArtistDetailsModal
  bind:open={openEditArtistDetailsModal}
  artist={data.artist}
  onResult={async (resultData) => {
    const res = await apiClient.editArtist(data.artist.id, {
      name: resultData.name,
      otherName: resultData.otherName,
      tags: resultData.tags.split(","),
    });
    if (!res.success) {
      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.success("Successfully updated artist");
    invalidateAll();
  }}
/>

<ChangeArtistPictureModal
  bind:open={openArtistPictureModal}
  onResult={async (formData) => {
    const toastId = toast.loading("Uploading image...");

    const res = await apiClient.changeArtistPicture(data.artist.id, formData);
    if (!res.success) {
      toast.dismiss(toastId);

      handleApiError(res.error);
      invalidateAll();
      return;
    }

    toast.dismiss(toastId);
    toast.success("Successfully uploaded image, force reload page to see it");
  }}
/>
