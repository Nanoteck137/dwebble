<script lang="ts">
  import { Breadcrumb, Button, Input } from "@nanoteck137/nano-ui";
  import { Eye, EyeOff, TypeOutlineIcon } from "lucide-svelte";

  const { data } = $props();

  const tokens = $state(data.tokens.map((t) => ({ ...t, revealed: false })));
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/account">Account</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Tokens</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="h-4"></div>

<Button href="/account/tokens/new">New Token</Button>

{#each tokens as token, i}
  <div>
    <p>{token.name}</p>
    <div class="flex gap-2">
      <Input value={token.id} type={token.revealed ? "text" : "password"} />
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          tokens[i].revealed = !tokens[i].revealed;
        }}
      >
        {#if token.revealed}
          <EyeOff />
        {:else}
          <Eye />
        {/if}
      </Button>
    </div>
  </div>
{/each}
