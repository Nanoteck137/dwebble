<script lang="ts">
  import { createEventDispatcher } from "svelte";

  export let progress: number = 0;

  const dispatch = createEventDispatcher();

  function clamp(number: number, min: number, max: number) {
    return Math.min(Math.max(number, min), max);
  }

  let dragging = false;

  let bar: HTMLDivElement;

  $: clampedProgress = clamp(progress, 0.0, 1.0);
  $: right = `${100 - clampedProgress * 100}%`;
  $: head = `${clampedProgress * 100}%`;
</script>

<svelte:window
  on:mouseup={() => {
    dragging = false;
  }}
  on:mousemove={(e) => {
    if (dragging) {
      e.preventDefault();
      e.stopPropagation();
      e.stopImmediatePropagation();
      const rect = bar.getBoundingClientRect();
      const progress = clamp((e.x - rect.x) / rect.width, 0.0, 1.0);
      dispatch("progress", progress);
    }
  }}
/>

<div
  bind:this={bar}
  class="group relative h-2 w-full bg-red-400"
  on:mousedown={(e) => {
    e.preventDefault();
    e.stopPropagation();
    e.stopImmediatePropagation();

    dragging = true;

    const rect = e.currentTarget.getBoundingClientRect();
    const progress = clamp((e.x - rect.x) / rect.width, 0.0, 1.0);
    dispatch("progress", progress);
  }}
>
  <div
    style="right: {right}"
    class="absolute bottom-0 left-0 right-0 top-0 bg-blue-400"
  ></div>
  <div
    style="left: {head}"
    class="absolute bottom-[50%] top-[50%] h-4 w-4 -translate-x-[50%] -translate-y-[50%] rounded-full bg-green-400 opacity-0 group-hover:opacity-100"
  ></div>
</div>
