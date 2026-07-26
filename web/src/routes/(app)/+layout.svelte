<script lang="ts">
  import {
    DiscAlbum,
    FileMusic,
    Heart,
    ListMusic,
    ListVideo,
    LogOut,
    Search,
    Server,
    User,
    Users,
  } from "lucide-svelte";
  import AudioPlayer from "$lib/components/audio/AudioPlayer.svelte";
  import MobilePlayer from "$lib/components/audio/MobilePlayer.svelte";
  import { getMusicManager } from "$lib/music-manager.svelte";
  import { onMount } from "svelte";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import toast, { Toaster } from "svelte-5-french-toast";
  import { getApiAddress, handleApiError, setApiClient } from "$lib";
  import { setMusicManager } from "$lib/music-manager.svelte";
  import { goto, invalidateAll } from "$app/navigation";
  import { setQuickPlaylist } from "$lib/quick-playlist.svelte";
  import { setFavorites } from "$lib/favorites.svelte";
  import {
    initPlaylistModalManager,
    showPlaylistModal,
  } from "$lib/playlist-modal.svelte";
  import PlaylistSelectorModal from "$lib/components/new-modals/PlaylistSelectorModal.svelte";
  import { isRoleAdmin } from "$lib/utils.js";
  import { page } from "$app/state";

  let { children, data } = $props();

  let apiClient = setApiClient(
    getApiAddress(page.url),
    localStorage.getItem("token") ?? undefined,
  );

  const musicManager = setMusicManager(apiClient);

  onMount(() => {
    if (data.user) {
      musicManager.initQueue();
    }
  });

  setFavorites(apiClient);

  initPlaylistModalManager(apiClient);

  let quickPlaylist = setQuickPlaylist(apiClient);

  $effect(() => {
    quickPlaylist.setPlaylistId(data.user?.quickPlaylist ?? null);
  });
</script>

<svelte:head>
  <title>Tunebook</title>
</svelte:head>

<Toaster position="bottom-right" />

<PlaylistSelectorModal />

{#if data.user}
  <header
    class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
  >
    <div
      class="container flex h-14 max-w-screen-2xl items-center gap-4 px-4 sm:px-8"
    >
      <a
        class="bg-gradient-to-tr from-logo-1 via-logo-2 to-logo-3 bg-clip-text text-2xl font-medium text-transparent"
        href="/">Tunebook</a
      >

      <div class="hidden items-center gap-1 md:flex">
        <a
          href="/albums"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
            '/albums',
          )
            ? 'bg-accent text-accent-foreground'
            : 'text-muted-foreground'}"
        >
          Albums
        </a>
        <a
          href="/artists"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
            '/artists',
          )
            ? 'bg-accent text-accent-foreground'
            : 'text-muted-foreground'}"
        >
          Artists
        </a>
        <a
          href="/tracks"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
            '/tracks',
          )
            ? 'bg-accent text-accent-foreground'
            : 'text-muted-foreground'}"
        >
          Tracks
        </a>
        <a
          href="/playlists"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
            '/playlists',
          )
            ? 'bg-accent text-accent-foreground'
            : 'text-muted-foreground'}"
        >
          Playlists
        </a>
        <a
          href="/favorites"
          class="rounded-md px-3 py-1.5 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground {page.url.pathname.startsWith(
            '/favorites',
          )
            ? 'bg-accent text-accent-foreground'
            : 'text-muted-foreground'}"
        >
          Favorites
        </a>
      </div>

      <div class="flex-grow"></div>

      <div class="flex items-center gap-2">
        <button
          class={buttonVariants({ variant: "ghost", size: "icon" })}
          onclick={async () => {
            const id = await showPlaylistModal({
              selectedId: data.user?.quickPlaylist ?? undefined,
              title: "Set Quick Playlist",
              description: "Choose which playlist to use for quick adds",
            });
            if (!id) return;

            const res = await apiClient.setQuickPlaylist({
              playlistId: id,
            });
            if (!res.success) {
              handleApiError(res.error);
              invalidateAll();
              return;
            }

            toast.success("Successfully set quick playlist");
            invalidateAll();
          }}
        >
          <ListVideo />
        </button>

        <Button href="/search" size="icon" variant="ghost">
          <Search />
        </Button>

        <DropdownMenu.Root>
          <DropdownMenu.Trigger>
            <img
              class="w-8 min-w-8 rounded-full"
              src={data.user.picture.small}
              alt=""
            />
          </DropdownMenu.Trigger>
          <DropdownMenu.Content class="w-56" align="end">
            <DropdownMenu.Group>
              <DropdownMenu.GroupHeading>
                {data.user.displayName}
              </DropdownMenu.GroupHeading>

              <DropdownMenu.Separator />

              <DropdownMenu.Item
                onSelect={() => {
                  if (!data.user) return;

                  goto(`/users/${data.user.id}`);
                }}
              >
                <User />
                Account
              </DropdownMenu.Item>

              {#if isRoleAdmin(data.user.role)}
                <DropdownMenu.Item
                  onSelect={() => {
                    goto(`/server`);
                  }}
                >
                  <Server />
                  Server
                </DropdownMenu.Item>
              {/if}

              <DropdownMenu.Separator />

              <DropdownMenu.Item
                onSelect={() => {
                  localStorage.removeItem("token");
                  musicManager.reset();
                  goto("/", { invalidateAll: true });
                }}
              >
                <LogOut />
                Logout
              </DropdownMenu.Item>
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>
      </div>
    </div>
  </header>
{/if}

<main class="container px-4 py-4 sm:px-8">
  {@render children()}
</main>

{#if getMusicManager().showPlayer}
  <MobilePlayer />
{/if}

<footer class="fixed bottom-0 z-40 w-full">
  <AudioPlayer />

  {#if data.user}
    <nav
      class="flex items-center justify-around border-t bg-background py-1 md:hidden"
    >
      <a
        href="/albums"
        class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
          '/albums',
        )
          ? 'text-primary'
          : 'text-muted-foreground'}"
      >
        <DiscAlbum size={18} />
        Albums
      </a>
      <a
        href="/artists"
        class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
          '/artists',
        )
          ? 'text-primary'
          : 'text-muted-foreground'}"
      >
        <Users size={18} />
        Artists
      </a>
      <a
        href="/tracks"
        class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
          '/tracks',
        )
          ? 'text-primary'
          : 'text-muted-foreground'}"
      >
        <FileMusic size={18} />
        Tracks
      </a>
      <a
        href="/playlists"
        class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
          '/playlists',
        )
          ? 'text-primary'
          : 'text-muted-foreground'}"
      >
        <ListMusic size={18} />
        Playlists
      </a>
      <a
        href="/favorites"
        class="flex flex-col items-center gap-0.5 px-3 py-1 text-xs font-medium transition-colors {page.url.pathname.startsWith(
          '/favorites',
        )
          ? 'text-primary'
          : 'text-muted-foreground'}"
      >
        <Heart size={18} />
        Favorites
      </a>
    </nav>
  {/if}
</footer>
