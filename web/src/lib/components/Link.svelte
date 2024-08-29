<script lang="ts">
  export let title: string;
  export let href: string | undefined = undefined;
  export let icon;
  export let onClick: (() => void) | undefined = undefined;
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<svelte:element
  this={href ? "a" : "button"}
  class="hover:text-[--on-primary flex w-full items-center gap-2 rounded px-2 py-2 hover:bg-[--primary] hover:text-[--on-primary]"
  onclick={() => onClick?.()}
  {href}
>
  <svelte:component this={icon} size="32"></svelte:component>
  <p class="text-xl">{title}</p>
</svelte:element>
