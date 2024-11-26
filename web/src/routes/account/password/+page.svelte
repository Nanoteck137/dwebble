<script>
  import {
    Breadcrumb,
    Button,
    Card,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
</script>

<div class="py-2">
  <Breadcrumb.Root>
    <Breadcrumb.List>
      <Breadcrumb.Item>
        <Breadcrumb.Link href="/account">Account</Breadcrumb.Link>
      </Breadcrumb.Item>
      <Breadcrumb.Separator />
      <Breadcrumb.Item>
        <Breadcrumb.Page>Change Password</Breadcrumb.Page>
      </Breadcrumb.Item>
    </Breadcrumb.List>
  </Breadcrumb.Root>
</div>

<div class="h-4"></div>

<form method="post">
  <Card.Root class="mx-auto max-w-[450px]">
    <Card.Header>
      <Card.Title>Change Password</Card.Title>
    </Card.Header>
    <Card.Content class="flex flex-col gap-4">
      <div class="flex flex-col gap-2">
        <Label for="currentPassword">Current Password</Label>
        <Input id="currentPassword" name="currentPassword" type="password" />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="newPassword">New Password</Label>
        <Input id="newPassword" name="newPassword" type="password" />
      </div>
      <div class="flex flex-col gap-2">
        <Label for="newPasswordConfirm">Confirm New Password</Label>
        <Input
          id="newPasswordConfirm"
          name="newPasswordConfirm"
          type="password"
        />
      </div>
    </Card.Content>
    <Card.Footer class="flex justify-end gap-4">
      <Button href="/account" variant="outline">Back</Button>
      <Button type="submit">Save</Button>
    </Card.Footer>
  </Card.Root>
</form>
