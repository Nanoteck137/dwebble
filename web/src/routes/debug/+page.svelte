<script lang="ts">
  import { Album, Artist, type Track } from "$lib/api/types";
  import AlbumListItem from "$lib/components/AlbumListItem.svelte";
  import ArtistListItem from "$lib/components/ArtistListItem.svelte";
  import TrackListItem from "$lib/components/TrackListItem.svelte";
  import { cn } from "$lib/utils";
  import {
    Button,
    buttonVariants,
    DropdownMenu,
    Separator,
  } from "@nanoteck137/nano-ui";
  import {
    DiscAlbum,
    EllipsisVertical,
    ListPlus,
    Star,
    Users,
  } from "lucide-svelte";

  const testTrack: Track = {
    number: 1,
    id: "123",
    name: "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Praesentium consectetur magnam repellendus animi, voluptate veniam cumque quibusdam error non numquam aliquam inventore quis unde fuga accusantium quaerat, reprehenderit quidem. Consequatur at iure accusantium culpa quia, amet quasi nemo. Sequi, nisi et? Dolore culpa nulla recusandae reiciendis ad delectus enim, corporis qui? Modi vero voluptatibus libero magni officia alias temporibus veniam, dolorum, aperiam quod sed consectetur ut. Maxime architecto accusantium labore officia molestiae ipsa aut dignissimos excepturi nemo. Fugit officia modi pariatur eaque eligendi sunt impedit accusantium excepturi ullam facere, quaerat consequatur quo facilis, animi quae placeat in dignissimos distinctio eius.",
    albumId: "123",
    artistId: "123",
    duration: 320,
    year: 1999,
    originalMediaUrl: "",
    mobileMediaUrl: "",
    coverArt: {
      original: "",
      small: "https://picsum.photos/200/300",
      medium: "",
      large: "",
    },
    albumName: "Test Album",
    // artistName:
    //   "Lorem ipsum dolor sit, amet consectetur adipisicing elit. Praesentium consectetur magnam repellendus animi, voluptate veniam cumque quibusdam error non numquam aliquam inventore quis unde fuga accusantium quaerat, reprehenderit quidem. Consequatur at iure accusantium culpa quia, amet quasi nemo. Sequi, nisi et? Dolore culpa nulla recusandae reiciendis ad delectus enim, corporis qui? Modi vero voluptatibus libero magni officia alias temporibus veniam, dolorum, aperiam quod sed consectetur ut. Maxime architecto accusantium labore officia molestiae ipsa aut dignissimos excepturi nemo. Fugit officia modi pariatur eaque eligendi sunt impedit accusantium excepturi ullam facere, quaerat consequatur quo facilis, animi quae placeat in dignissimos distinctio eius.",
    artistName: "Test Artist",
    created: 0,
    updated: 0,
    tags: ["Good", "Rock", "Metal", "Classic"],
  };

  const testAlbum: Album = {
    id: "123",
    name: "Test Album",
    artistId: "123",
    year: 2013,
    coverArt: {
      original: "",
      small: "https://picsum.photos/200/300",
      medium: "",
      large: "",
    },
    artistName: "Test Artist",
    created: 0,
    updated: 0,
  };

  const testArtist: Artist = {
    id: "123",
    name: "Test Artist",
    picture: {
      original: "",
      small: "https://picsum.photos/200/300",
      medium: "",
      large: "",
    },
  };
</script>

<p class="pb-4 text-lg font-semibold">Track List Item</p>

<TrackListItem track={testTrack} onPlayClicked={() => {}} />
<TrackListItem
  track={{
    ...testTrack,
    tags: [],
    coverArt: {
      original: "",
      small: "",
      medium: "",
      large: "",
    },
  }}
>
  <Button class="rounded-full" variant="ghost" size="icon-lg">
    <Star fill="#fff" />
  </Button>

  <Button class="rounded-full" variant="ghost" size="icon-lg">
    <ListPlus />
  </Button>

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
        <DropdownMenu.Item class="p-0">
          <!-- svelte-ignore a11y_invalid_attribute -->
          <a class="flex h-full w-full items-center gap-2 px-3 py-2" href="#">
            <DiscAlbum size="16" />
            Go to Album
          </a>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="p-0">
          <!-- svelte-ignore a11y_invalid_attribute -->
          <a class="flex h-full w-full items-center gap-2 px-3 py-2" href="#">
            <Users size="16" />
            Go to Artist
          </a>
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
</TrackListItem>

<Separator />

<p class="py-4 text-lg font-semibold">Album List Item</p>

<AlbumListItem album={testAlbum} link />
<AlbumListItem
  album={{
    ...testAlbum,
    coverArt: {
      original: "",
      small: "",
      medium: "",
      large: "",
    },
  }}
>
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
        <DropdownMenu.Item class="p-0">
          <!-- svelte-ignore a11y_invalid_attribute -->
          <a class="flex h-full w-full items-center gap-2 px-3 py-2" href="#">
            <DiscAlbum size="16" />
            Go to Album
          </a>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="p-0">
          <!-- svelte-ignore a11y_invalid_attribute -->
          <a class="flex h-full w-full items-center gap-2 px-3 py-2" href="#">
            <Users size="16" />
            Go to Artist
          </a>
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
</AlbumListItem>

<p class="py-4 text-lg font-semibold">Artist List Item</p>

<ArtistListItem artist={testArtist} link />
<ArtistListItem
  artist={{
    ...testArtist,
    picture: {
      original: "",
      small: "",
      medium: "",
      large: "",
    },
  }}
>
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
        <DropdownMenu.Item class="p-0">
          <!-- svelte-ignore a11y_invalid_attribute -->
          <a class="flex h-full w-full items-center gap-2 px-3 py-2" href="#">
            <DiscAlbum size="16" />
            Go to Album
          </a>
        </DropdownMenu.Item>
        <DropdownMenu.Item class="p-0">
          <!-- svelte-ignore a11y_invalid_attribute -->
          <a class="flex h-full w-full items-center gap-2 px-3 py-2" href="#">
            <Users size="16" />
            Go to Artist
          </a>
        </DropdownMenu.Item>
      </DropdownMenu.Group>
    </DropdownMenu.Content>
  </DropdownMenu.Root>
</ArtistListItem>
