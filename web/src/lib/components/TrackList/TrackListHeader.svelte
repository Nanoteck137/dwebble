<script lang="ts">
  import { EllipsisVertical, ListPlus, Play } from "lucide-svelte";


  interface Props {
    cover: string;
    name: string;
    onPlay: () => void;
  }

  let { cover, name, onPlay }: Props = $props();
</script>

<div class="w-80">
  <img class="aspect-square w-80 rounded object-cover" src={cover} alt="" />
  <div class="h-2"></div>
  <p class="line-clamp-2 text-center font-bold">{name}</p>
  <div class="h-2"></div>
  <div class="flex items-center justify-center gap-5">
    <button
      class="button-secondary flex h-10 w-10 items-center justify-center rounded-full p-2 hover:brightness-75 active:scale-95"
      onclick={() => {}}
      title="Quick Add"
    >
      <ListPlus size="30" />
    </button>

    <button
      class="button-primary rounded-full p-2 hover:brightness-75 active:scale-95"
      onclick={onPlay}
      title="Play"
    >
      <Play size="30" />
    </button>

    <button
      class="button-secondary flex h-10 w-10 items-center justify-center rounded-full p-2 hover:brightness-75 active:scale-95"
      onclick={() => {}}
      title="More"
    >
      <EllipsisVertical size="30" />
    </button>
  </div>
  <div class="h-2"></div>
</div>
