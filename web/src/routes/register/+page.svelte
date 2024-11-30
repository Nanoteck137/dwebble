<script>
  import { Button, Card, Input, Label } from "@nanoteck137/nano-ui";
  import SuperDebug, { superForm } from "sveltekit-superforms";

  const { data } = $props();

  const { form, errors, enhance } = superForm(data.form, { onError: "apply" });
</script>

<form method="post" use:enhance>
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Register</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="username">Username</Label>
        <Input
          id="username"
          name="username"
          type="text"
          bind:value={$form.username}
        />
        {#if $errors.username}
          {#each $errors.username as error}
            {#if error.length > 0}
              <p>{error}</p>
            {/if}
          {/each}
        {/if}
      </div>
      <div class="flex flex-col gap-2">
        <Label for="password">Password</Label>
        <Input
          id="password"
          name="password"
          type="password"
          bind:value={$form.password}
        />
        {#if $errors.password}
          {#each $errors.password as error}
            <p>{error}</p>
          {/each}
        {/if}
      </div>
      <div class="flex flex-col gap-2">
        <Label for="passwordConfirm">Confirm Password</Label>
        <Input
          id="passwordConfirm"
          name="passwordConfirm"
          type="password"
          bind:value={$form.passwordConfirm}
        />
        {#if $errors.passwordConfirm}
          {#each $errors.passwordConfirm as error}
            <p>{error}</p>
          {/each}
        {/if}
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button type="submit">Register</Button>
    </Card.Footer>
  </Card.Root>
</form>

<SuperDebug data={$form} />
