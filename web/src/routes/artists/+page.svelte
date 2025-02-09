<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { page } from "$app/stores";
  import { getApiClient, handleApiError } from "$lib";
  import { Artist } from "$lib/api/types";
  import { cn, formatName } from "$lib/utils";
  import {
    Button,
    buttonVariants,
    Checkbox,
    DropdownMenu,
    Pagination,
    Separator,
  } from "@nanoteck137/nano-ui";
  import { EllipsisVertical, Merge, Users, X } from "lucide-svelte";
  import toast from "svelte-5-french-toast";

  let { data } = $props();
  const apiClient = getApiClient();

  let merge = $state<string | null>();
  let selected = $state<string[]>([]);

  function isSelected(id: string) {
    return !!selected.find((i) => i === id);
  }
</script>

{#snippet artistItem(artist: Artist)}
  {@const isMergeTarget = artist.id === merge}
  {@const isSelectedTarget = isSelected(artist.id)}
  <div class="py-2">
    <div
      class={`relative flex items-center gap-2 rounded pr-2 ${isMergeTarget ? "bg-muted text-muted-foreground" : ""}`}
    >
      {#if merge}
        <!-- svelte-ignore a11y_consider_explicit_label -->
        <button
          class="absolute inset-0"
          onclick={() => {
            if (isSelected(artist.id)) {
              selected = selected.filter((i) => i !== artist.id);
            } else {
              selected.push(artist.id);
            }
          }}
        ></button>
      {/if}
      <a href={`/artists/${artist.id}`}>
        <img
          class="inline-flex aspect-square w-14 min-w-14 items-center justify-center rounded border object-cover text-xs"
          src={artist.picture.small}
          alt="cover"
          loading="lazy"
        />
      </a>
      <div class="flex flex-grow flex-col">
        <div class="flex items-center gap-1">
          <a
            class="line-clamp-1 w-fit text-sm font-medium"
            href={`/artists/${artist.id}`}
            title={formatName(artist.name)}
          >
            {formatName(artist.name)}
          </a>
        </div>

        <!-- <ArtistList class="text-muted-foreground" artists={track.allArtists} /> -->

        <p class="line-clamp-1 text-xs text-muted-foreground">
          {#if artist.tags.length > 0}
            {artist.tags.join(", ")}
          {:else}
            No Tags
          {/if}
        </p>
      </div>
      <div class="flex items-center">
        {#if merge}
          {#if !isMergeTarget}
            <Checkbox checked={isSelectedTarget} />
          {/if}
        {:else}
          <DropdownMenu.Root>
            <DropdownMenu.Trigger
              class={cn(
                buttonVariants({ variant: "ghost", size: "icon-lg" }),
                "rounded-full",
              )}
            >
              <EllipsisVertical />
            </DropdownMenu.Trigger>
            <DropdownMenu.Content align="end">
              <DropdownMenu.Group>
                <DropdownMenu.Item
                  onSelect={() => {
                    goto(`/artists/${artist.id}`);
                  }}
                >
                  <Users />
                  Go to Artist
                </DropdownMenu.Item>
                <DropdownMenu.Item
                  onSelect={() => {
                    merge = artist.id;
                    selected = [];
                  }}
                >
                  <Merge />
                  Merge Artist
                </DropdownMenu.Item>
              </DropdownMenu.Group>
            </DropdownMenu.Content>
          </DropdownMenu.Root>
        {/if}
      </div>
    </div>
  </div>
{/snippet}

<div class="flex items-center justify-between">
  <p class="text-bold text-xl">Artists</p>
  <p class="text-sm">{data.page.totalItems} artist(s)</p>
</div>

<div class="flex flex-col">
  {#each data.artists as artist}
    {@render artistItem(artist)}
    <Separator />
  {/each}
</div>

<div class="h-4"></div>

{#if merge}
  <div
    class="bottom-popup sticky border border-border/40 bg-background/95 px-6 py-3 text-foreground backdrop-blur supports-[backdrop-filter]:bg-background/60"
  >
    <p class="text-center">{selected.length} artists selected</p>
    <div class="h-2"></div>

    <div class="flex flex-col justify-center gap-2 md:flex-row">
      <div class="flex gap-2">
        <Button
          variant="outline"
          size="icon"
          onclick={() => {
            merge = null;
            selected = [];
          }}
        >
          <X />
        </Button>

        <Button
          class="flex-grow"
          variant="outline"
          onclick={async () => {
            if (!merge) {
              return;
            }

            const res = await apiClient.mergeArtists(merge, {
              artists: selected,
            });
            if (!res.success) {
              handleApiError(res.error);
              return;
            }

            merge = null;
            selected = [];
            toast.success("Merge artists");
            await invalidateAll();
          }}
        >
          <Merge />
          Merge Artists
        </Button>
      </div>
    </div>
  </div>
{/if}

<Pagination.Root
  page={data.page.page + 1}
  count={data.page.totalItems}
  perPage={data.page.perPage}
  siblingCount={1}
  onPageChange={(p) => {
    const query = $page.url.searchParams;
    query.set("page", (p - 1).toString());

    goto(`?${query.toString()}`, { invalidateAll: true, keepFocus: true });
  }}
>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.PrevButton />
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link
              href="?page={page.value}"
              {page}
              isActive={currentPage === page.value}
            >
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
