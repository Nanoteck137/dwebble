<script>
  import { Button, Card, Input, Label } from "@nanoteck137/nano-ui";
</script>

<form method="post">
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Login</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="username">Username</Label>
        <Input id="username" name="username" type="text" />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="password">Password</Label>
        <Input id="password" name="password" type="password" />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/account" variant="outline">Back</Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
