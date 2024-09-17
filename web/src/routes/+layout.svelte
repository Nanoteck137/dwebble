<script lang="ts">
  import {
    DiscAlbum,
    FileMusic,
    Home,
    LayoutDashboard,
    ListMusic,
    LogIn,
    LogOut,
    Menu,
    Server,
    User,
    Users,
  } from "lucide-svelte";
  import "../app.css";
  import AudioPlayer from "$lib/components/AudioPlayer.svelte";
  import Link from "$lib/components/Link.svelte";
  import { browser } from "$app/environment";
  import { fade, fly } from "svelte/transition";

  let { children, data } = $props();

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

<header
  class="fixed left-0 right-0 top-0 z-30 flex h-16 items-center gap-4 bg-[--bg-color] px-4 py-2"
>
  <button
    onclick={() => {
      showSideMenu = true;
    }}
  >
    <Menu size="32" />
  </button>

  <a class="text-3xl font-medium text-[--logo-color]" href="/">Dwebble</a>
</header>

<main class="mt-16">
  {@render children()}
</main>

<footer>
  <AudioPlayer />
</footer>

{#if showSideMenu}
  <button
    class="fixed bottom-0 left-0 right-0 top-0 z-50 bg-[--modal-overlay-bg]"
    onclick={() => {
      showSideMenu = false;
    }}
    transition:fade
  ></button>

  <aside
    class={`fixed bottom-0 top-0 z-50 flex w-80 flex-col bg-[--bg-color] px-4 text-[--fg-color]`}
    transition:fly={{ x: -400 }}
  >
    <div class="flex h-16 items-center gap-4 px-2 py-2">
      <button
        onclick={() => {
          showSideMenu = false;
        }}
      >
        <Menu size="32" />
      </button>
      <a
        class="text-3xl font-medium text-[--logo-color]"
        href="/"
        onclick={() => {
          showSideMenu = false;
        }}
      >
        Dwebble
      </a>
    </div>

    <div class="flex flex-col gap-2 py-2">
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
      {/if}
    </div>
    <div class="flex-grow"></div>
    <div class="flex flex-col gap-2 py-4">
      {#if data.user}
        <Link
          title={data.user.username}
          href="/account"
          icon={User}
          onClick={close}
        />

        {#if data.user.isOwner}
          <Link title="Server" href="/server" icon={Server} onClick={close} />
        {/if}

        <form class="w-full" action="/logout" method="POST">
          <Link title="Logout" icon={LogOut} onClick={close} />
        </form>
      {:else}
        <Link title="Login" href="/login" icon={LogIn} onClick={close} />
      {/if}
    </div>
  </aside>
{/if}

<svelte:head>
  <title>Dwebble</title>
</svelte:head>
