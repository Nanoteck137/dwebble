<script>
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Card, Input, Label } from "@nanoteck137/nano-ui";
  import { superForm } from "sveltekit-superforms";

  const { data } = $props();

  const { form, errors, enhance } = superForm(data.form, { onError: "apply" });
</script>

<form method="post" use:enhance>
  <Card.Root class="mx-auto max-w-[465px]">
    <Card.Header>
      <Card.Title>Login</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <FormItem>
        <Label for="username">Username</Label>
        <Input
          id="username"
          name="username"
          type="text"
          bind:value={$form.username}
        />
        <Errors errors={$errors.username} />
      </FormItem>

      <FormItem>
        <Label for="password">Password</Label>
        <Input
          id="password"
          name="password"
          type="password"
          bind:value={$form.password}
        />
        <Errors errors={$errors.password} />
      </FormItem>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/register" variant="ghost">Register</Button>
      <Button type="submit">Login</Button>
    </Card.Footer>
  </Card.Root>
</form>

<!-- <SuperDebug data={$form} /> -->
