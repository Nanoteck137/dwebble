<script lang="ts">
  import {
    DiscAlbum,
    FileMusic,
    Home,
    ListMusic,
    LogIn,
    LogOut,
    Menu,
    Search,
    Server,
    Tags,
    User,
    Users,
  } from "lucide-svelte";
  import "../app.css";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";
  import Link from "$lib/components/Link.svelte";
  import { browser } from "$app/environment";
  import { fade, fly } from "svelte/transition";
  import { Button } from "@nanoteck137/nano-ui";
  import { setModalState } from "$lib/modal.svelte";
  import ConfirmModal from "$lib/components/modals/ConfirmModal.svelte";
  import QueryArtistModal from "$lib/components/modals/QueryArtistModal.svelte";

  let { children, data } = $props();

  const modalState = setModalState();

  let showSideMenu = $state(false);

  function close() {
    showSideMenu = false;
  }

  $effect(() => {
    if (showSideMenu) {
      if (browser) document.body.style.overflow = "hidden";
    } else {
      if (browser) document.body.style.overflow = "";
    }
  });
</script>

<svelte:head>
  <title>Dwebble</title>
</svelte:head>

<!-- svelte-ignore a11y_consider_explicit_label -->
{#each modalState.modals as modal}
  {#if modal.data.type === "modal-confirm"}
    <ConfirmModal {modalState} data={modal.data} />
  {:else if modal.data.type === "modal-query-artist"}
    <QueryArtistModal {modalState} data={modal.data} />
  {/if}
{/each}

<header
  class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
>
  <div class="container flex h-14 max-w-screen-2xl items-center gap-4">
    <button
      onclick={() => {
        showSideMenu = true;
      }}
    >
      <Menu size="20" />
    </button>

    <a class="text-2xl font-medium text-[--logo-color]" href="/">Dwebble</a>
    <p>{data.user?.role}</p>

    <div class="flex-grow"></div>

    <Button href="/search" size="icon" variant="ghost">
      <Search />
    </Button>
  </div>
</header>

<main class="container py-4">
  {@render children()}
</main>

<footer class="fixed bottom-0 w-screen">
  <AudioPlayer />
</footer>

{#if showSideMenu}
  <!-- svelte-ignore a11y_consider_explicit_label -->
  <button
    class="fixed inset-0 z-50 bg-black/80"
    onclick={() => {
      showSideMenu = false;
    }}
    transition:fade={{ duration: 200 }}
  ></button>

  <aside
    class={`fixed bottom-0 top-0 z-50 flex w-72 flex-col bg-sidebar text-sidebar-foreground`}
    transition:fly={{ x: -400 }}
  >
    <div class="flex h-14 items-center gap-4 border-b px-8">
      <button
        onclick={() => {
          showSideMenu = false;
        }}
      >
        <Menu size="20" />
      </button>
      <a
        class="text-2xl font-medium"
        href="/"
        onclick={() => {
          showSideMenu = false;
        }}
      >
        Dwebble
      </a>
    </div>

    <div class="flex flex-col gap-2 px-4 py-4">
      <Link title="Home" href="/" icon={Home} onClick={close} />
      <Link title="Artists" href="/artists" icon={Users} onClick={close} />
      <Link title="Albums" href="/albums" icon={DiscAlbum} onClick={close} />
      <Link title="Tracks" href="/tracks" icon={FileMusic} onClick={close} />

      {#if data.user}
        <Link
          title="Playlists"
          href="/playlists"
          icon={ListMusic}
          onClick={close}
        />

        <Link title="Taglists" href="/taglists" icon={Tags} onClick={close} />
      {/if}
    </div>
    <div class="flex-grow"></div>
    <div class="flex flex-col gap-2 px-4 py-2">
      {#if data.user}
        <Link
          title={data.user.username}
          href="/account"
          icon={User}
          onClick={close}
        />

        {#if data.user.role === "super_user"}
          <Link title="Server" href="/server" icon={Server} onClick={close} />
        {/if}

        <form class="w-full" action="/logout" method="POST">
          <Link title="Logout" icon={LogOut} onClick={close} />
        </form>
      {:else}
        <Link title="Login" href="/login" icon={LogIn} onClick={close} />
      {/if}
    </div>
    <div class="h-4"></div>
  </aside>
{/if}
