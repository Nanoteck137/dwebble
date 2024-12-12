<script>
  import Errors from "$lib/components/Errors.svelte";
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import SuperDebug, { superForm } from "sveltekit-superforms";

  const { data } = $props();

  const { form, errors, enhance } = superForm(data.form, { onError: "apply" });
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/account">Account</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/account/tokens">Tokens</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>New</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="h-4"></div>

<form method="post" use:enhance>
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>New Api Token</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/account/tokens" variant="outline">Back</Button>
      <Button type="submit">Create</Button>
    </Card.Footer>
  </Card.Root>
</form>

<SuperDebug data={$form} />
