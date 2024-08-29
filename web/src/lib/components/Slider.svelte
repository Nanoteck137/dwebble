<script lang="ts">
  export let value: number;
  export let onValue: (value: number) => void;

  let dragging = false;
  let dragValue = value;

  let sliderDiv: HTMLDivElement;

  // $: console.log(value);
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="relative h-2 w-full touch-none bg-red-400"
  bind:this={sliderDiv}
  onclick={(e) => {
    const target = e.target! as HTMLDivElement;
    const rect = target.getBoundingClientRect();
    const x = e.clientX - rect.x;
    const percent = x / rect.width;
    onValue(percent);
  }}
>
  <div
    class="absolute top-1/2 h-4 w-4 -translate-x-1/2 -translate-y-1/2 rounded-full bg-white"
    style={`left: ${(dragging ? dragValue : value) * 100}%`}
    onclick={(e) => {
      e.preventDefault();
      e.stopPropagation();
      e.stopImmediatePropagation();
    }}
    onmousedown={(e) => {
      e.preventDefault();

      dragging = true;
      dragValue = value;

      const move = (e: MouseEvent) => {
        e.preventDefault();

        const rect = sliderDiv.getBoundingClientRect();
        const x = e.clientX - rect.x;
        let percent = x / rect.width;
        console.log(percent);

        if (percent > 1.0) percent = 1.0;
        if (percent < 0.0) percent = 0.0;

        dragValue = percent;
      };

      const up = (e: MouseEvent) => {
        e.stopPropagation();

        document.removeEventListener("mousemove", move);
        document.removeEventListener("mouseup", up);

        dragging = false;
        onValue(dragValue);
        value = dragValue;
      };

      document.addEventListener("mousemove", move);
      document.addEventListener("mouseup", up);
    }}
    ontouchstart={(e) => {
      e.preventDefault();
      e.stopPropagation();
      e.stopImmediatePropagation();

      dragging = true;
      dragValue = value;

      const move = (e: TouchEvent) => {
        e.preventDefault();
        e.stopPropagation();
        e.stopImmediatePropagation();

        const rect = sliderDiv.getBoundingClientRect();
        const x = e.touches[0].clientX - rect.x;
        let percent = x / rect.width;
        console.log(percent);

        if (percent > 1.0) percent = 1.0;
        if (percent < 0.0) percent = 0.0;

        dragValue = percent;
      };

      const up = (e: TouchEvent) => {
        e.preventDefault();
        e.stopPropagation();
        e.stopImmediatePropagation();

        document.removeEventListener("touchmove", move);
        document.removeEventListener("touchend", up);

        dragging = false;
        onValue(dragValue);
        value = dragValue;
      };

      document.addEventListener("touchmove", move);
      document.addEventListener("touchend", up);
    }}
  ></div>
</div>
