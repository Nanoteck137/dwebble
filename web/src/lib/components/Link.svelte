<script lang="ts">
  interface Props {
    title: string;
    href?: string | undefined;
    icon: any;
    onClick?: (() => void) | undefined;
  }

  let {
    title,
    href = undefined,
    icon,
    onClick = undefined
  }: Props = $props();

  const SvelteComponent = $derived(icon);
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<svelte:element
  this={href ? "a" : "button"}
  class="flex w-full items-center gap-2 rounded px-2 py-2 hover:bg-[--primary] hover:text-[--on-primary]"
  onclick={() => onClick?.()}
  {href}
>
  <SvelteComponent size="32" strokeWidth="1.5"></SvelteComponent>
  <p class="text-xl">{title}</p>
</svelte:element>
