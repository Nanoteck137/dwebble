<script>
  import { Breadcrumb, Button, Input, Label } from "@nanoteck137/nano-ui";
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

<form class="flex flex-col gap-4" method="post">
  <p>Change Password</p>
  <div class="flex flex-col gap-4">
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
  </div>

  <div class="flex justify-end gap-2">
    <Button href="/account" variant="ghost">Back</Button>
    <Button type="submit">Save</Button>
  </div>
</form>
