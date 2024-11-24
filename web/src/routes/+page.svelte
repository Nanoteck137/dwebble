<script lang="ts">
  import {
    AlertDialog,
    AspectRatio,
    Avatar,
    Breadcrumb,
    Button,
    buttonVariants,
    Card,
    Checkbox,
    Dialog,
    Drawer,
    DropdownMenu,
    Input,
    Label,
    Pagination,
    Popover,
    Progress,
    RadioGroup,
    ScrollArea,
    Select,
    Separator,
    Skeleton,
    Slider,
    Switch,
    Tabs,
    Toggle,
    ToggleGroup,
    Tooltip,
    type ButtonSize,
    type ButtonVariant,
    Command,
    Sheet,
  } from "@nanoteck137/nano-ui";
  import {
    Bold,
    Calculator,
    Calendar,
    CreditCard,
    Italic,
    Settings,
    Smile,
    Underline,
    User,
  } from "lucide-svelte";
  import { onMount } from "svelte";

  const variants: ButtonVariant[] = [
    "default",
    "destructive",
    "ghost",
    "link",
    "outline",
    "secondary",
  ];

  const sizes: ButtonSize[] = ["sm", "default", "lg"];

  let checked = $state(false);

  let progressValue = $state(0);

  onMount(() => {
    const step = 10;
    let negative = false;
    const t = setInterval(() => {
      if (negative) {
        progressValue -= step;
        if (progressValue <= 0) {
          progressValue = 0;
          negative = !negative;
        }
      } else {
        progressValue += step;
        if (progressValue > 100) {
          progressValue = 100;
          negative = !negative;
        }
      }
    }, 500);
    return () => clearInterval(t);
  });

  const tags = Array.from({ length: 50 }).map(
    (_, i, a) => `v1.2.0-beta.${a.length - i}`,
  );

  const fruits = [
    { value: "apple", label: "Apple" },
    { value: "banana", label: "Banana" },
    { value: "blueberry", label: "Blueberry" },
    { value: "grapes", label: "Grapes" },
    { value: "pineapple", label: "Pineapple" },
  ];

  let selectValue = $state("");

  const triggerContent = $derived(
    fruits.find((f) => f.value === selectValue)?.label ?? "Select a fruit",
  );

  const SHEET_SIDES = ["top", "right", "bottom", "left"] as const;
</script>

<p class="p-4 text-xl">Home Page</p>

<Breadcrumb.Root>
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="#">Home</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <DropdownMenu.Root>
        <DropdownMenu.Trigger class="flex items-center gap-1">
          <Breadcrumb.Ellipsis class="size-4" />
          <span class="sr-only">Toggle menu</span>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content align="start">
          <DropdownMenu.Item>Documentation</DropdownMenu.Item>
          <DropdownMenu.Item>Themes</DropdownMenu.Item>
          <DropdownMenu.Item>GitHub</DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Root>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Link href="#">Components</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page>Breadcrumb</Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<div class="flex flex-col gap-2">
  {#each sizes as size}
    <div class="flex gap-2">
      <p class="w-28">size: {size}</p>
      {#each variants as variant}
        <div class="flex w-40 justify-center">
          <Button title="variant: {variant?.toString()}" {size} {variant}>
            Hello World
          </Button>
        </div>
      {/each}
    </div>
  {/each}
</div>

<div class="flex items-center space-x-2">
  <Checkbox id="terms" bind:checked />
  <Label
    for="terms"
    class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
  >
    Checkbox
  </Label>
</div>

<div class="flex items-center space-x-2">
  <Switch id="airplane-mode" />
  <Label for="airplane-mode">Airplane Mode</Label>
</div>

<div class="p-4">
  <Input placeholder="Input" />
</div>

<Separator />

<div class="p-4">
  <Slider value={[50]} min={0} max={100} />
</div>

<div class="h-1"></div>

<div class="flex justify-center">
  <Card.Root class="w-[350px]">
    <Card.Header>
      <Card.Title>Card Title</Card.Title>
      <Card.Description>Card Description</Card.Description>
    </Card.Header>
    <Card.Content>
      <p>Card Content</p>
    </Card.Content>
    <Card.Footer>
      <p>Card Footer</p>
    </Card.Footer>
  </Card.Root>
</div>

<div class="flex items-center space-x-4 p-4">
  <Skeleton class="size-12 rounded-full" />
  <div class="space-y-2">
    <Skeleton class="h-4 w-[250px]" />
    <Skeleton class="h-4 w-[200px]" />
  </div>
</div>

<Progress value={progressValue} max={100} class="w-[60%]" />

<ScrollArea class="h-72 w-48 rounded-md border">
  <div class="p-4">
    <h4 class="mb-4 text-sm font-medium leading-none">Tags</h4>
    {#each tags as tag}
      <div class="text-sm">
        {tag}
      </div>
      <Separator class="my-2" />
    {/each}
  </div>
</ScrollArea>

<Select.Root type="single" name="favoriteFruit" bind:value={selectValue}>
  <Select.Trigger class="w-[180px]">
    {triggerContent}
  </Select.Trigger>
  <Select.Content>
    <Select.Group>
      <Select.GroupHeading>Fruits</Select.GroupHeading>
      {#each fruits as fruit}
        <Select.Item value={fruit.value} label={fruit.label} />
      {/each}
    </Select.Group>
  </Select.Content>
</Select.Root>

<Avatar.Root>
  <Avatar.Image src="https://github.com/nanoteck137.png" alt="@nanoteck137" />
  <Avatar.Fallback>NT</Avatar.Fallback>
</Avatar.Root>

<AspectRatio ratio={16 / 9} class="bg-muted">
  <img
    src="https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=800&dpr=2&q=80"
    alt="Gray by Drew Beamer"
    class="h-full w-full rounded-md object-cover"
  />
</AspectRatio>

<RadioGroup.Root value="comfortable" name="spacing">
  <div class="flex items-center space-x-2">
    <RadioGroup.Item value="default" id="r1" />
    <Label for="r1">Default</Label>
  </div>
  <div class="flex items-center space-x-2">
    <RadioGroup.Item value="comfortable" id="r2" />
    <Label for="r2">Comfortable</Label>
  </div>
  <div class="flex items-center space-x-2">
    <RadioGroup.Item value="compact" id="r3" />
    <Label for="r3">Compact</Label>
  </div>
</RadioGroup.Root>

<ToggleGroup.Root type="multiple">
  <ToggleGroup.Item value="bold" aria-label="Toggle bold">
    <Bold />
  </ToggleGroup.Item>
  <ToggleGroup.Item value="italic" aria-label="Toggle italic">
    <Italic />
  </ToggleGroup.Item>
  <ToggleGroup.Item value="strikethrough" aria-label="Toggle strikethrough">
    <Underline />
  </ToggleGroup.Item>
</ToggleGroup.Root>

<Toggle aria-label="toggle bold">
  <Bold />
</Toggle>

<Command.Root class="max-w-[450px] rounded-lg border shadow-md">
  <Command.Input placeholder="Type a command or search..." />
  <Command.List>
    <Command.Empty>No results found.</Command.Empty>
    <Command.Group heading="Suggestions">
      <Command.Item>
        <Calendar />
        <span>Calendar</span>
      </Command.Item>
      <Command.Item>
        <Smile />
        <span>Search Emoji</span>
      </Command.Item>
      <Command.Item disabled>
        <Calculator />
        <span>Calculator</span>
      </Command.Item>
    </Command.Group>
    <Command.Separator />
    <Command.Group heading="Settings">
      <Command.Item>
        <User />
        <span>Profile</span>
        <Command.Shortcut>⌘P</Command.Shortcut>
      </Command.Item>
      <Command.Item>
        <CreditCard />
        <span>Billing</span>
        <Command.Shortcut>⌘B</Command.Shortcut>
      </Command.Item>
      <Command.Item>
        <Settings />
        <span>Settings</span>
        <Command.Shortcut>⌘S</Command.Shortcut>
      </Command.Item>
    </Command.Group>
  </Command.List>
</Command.Root>

<Tabs.Root value="account" class="w-[400px]">
  <Tabs.List>
    <Tabs.Trigger value="account">Account</Tabs.Trigger>
    <Tabs.Trigger value="password">Password</Tabs.Trigger>
  </Tabs.List>
  <Tabs.Content value="account">
    Make changes to your account here.
  </Tabs.Content>
  <Tabs.Content value="password">Change your password here.</Tabs.Content>
</Tabs.Root>

<div class="grid grid-cols-2 gap-2">
  {#each SHEET_SIDES as side, _ (side)}
    <Sheet.Root>
      <Sheet.Trigger class={buttonVariants({ variant: "outline" })}>
        {side}
      </Sheet.Trigger>
      <Sheet.Content {side}>
        <Sheet.Header>
          <Sheet.Title>Edit profile</Sheet.Title>
          <Sheet.Description>
            Make changes to your profile here. Click save when you're done.
          </Sheet.Description>
        </Sheet.Header>
        <div class="grid gap-4 py-4">
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="name" class="text-right">Name</Label>
            <Input id="name" value="Pedro Duarte" class="col-span-3" />
          </div>
          <div class="grid grid-cols-4 items-center gap-4">
            <Label for="username" class="text-right">Username</Label>
            <Input id="username" value="@peduarte" class="col-span-3" />
          </div>
        </div>
        <Sheet.Footer>
          <Sheet.Close
            class={buttonVariants({ variant: "outline" })}
            type="submit"
          >
            Save changes
          </Sheet.Close>
        </Sheet.Footer>
      </Sheet.Content>
    </Sheet.Root>
  {/each}
</div>

<Drawer.Root>
  <Drawer.Trigger class={buttonVariants({ variant: "default" })}>
    Open Drawer
  </Drawer.Trigger>
  <Drawer.Content>
    <Drawer.Header>
      <Drawer.Title>Are you sure absolutely sure?</Drawer.Title>
      <Drawer.Description>This action cannot be undone.</Drawer.Description>
    </Drawer.Header>
    <Drawer.Footer>
      <Button>Submit</Button>
      <Drawer.Close>Cancel</Drawer.Close>
    </Drawer.Footer>
  </Drawer.Content>
</Drawer.Root>

<DropdownMenu.Root>
  <DropdownMenu.Trigger class={buttonVariants({ variant: "default" })}>
    Open Dropdown
  </DropdownMenu.Trigger>
  <DropdownMenu.Content>
    <DropdownMenu.Group>
      <DropdownMenu.GroupHeading>My Account</DropdownMenu.GroupHeading>
      <DropdownMenu.Separator />
      <DropdownMenu.Item>Profile</DropdownMenu.Item>
      <DropdownMenu.Item>Billing</DropdownMenu.Item>
      <DropdownMenu.Item>Team</DropdownMenu.Item>
      <DropdownMenu.Item>Subscription</DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>

<AlertDialog.Root>
  <AlertDialog.Trigger class={buttonVariants({ variant: "outline" })}>
    Show Alert Dialog
  </AlertDialog.Trigger>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Are you absolutely sure?</AlertDialog.Title>
      <AlertDialog.Description>
        This action cannot be undone. This will permanently delete your account
        and remove your data from our servers.
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>Cancel</AlertDialog.Cancel>
      <AlertDialog.Action>Continue</AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>

<Dialog.Root>
  <Dialog.Trigger class={buttonVariants({ variant: "outline" })}>
    Open Dialog
  </Dialog.Trigger>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Are you sure absolutely sure?</Dialog.Title>
      <Dialog.Description>
        This action cannot be undone. This will permanently delete your account
        and remove your data from our servers.
      </Dialog.Description>
    </Dialog.Header>
  </Dialog.Content>
</Dialog.Root>

<Popover.Root>
  <Popover.Trigger class={buttonVariants({ variant: "outline" })}>
    Open Popover
  </Popover.Trigger>
  <Popover.Content class="w-80">
    <div class="grid gap-4">
      <div class="space-y-2">
        <h4 class="font-medium leading-none">Dimensions</h4>
        <p class="text-sm text-muted-foreground">
          Set the dimensions for the layer.
        </p>
      </div>
      <div class="grid gap-2">
        <div class="grid grid-cols-3 items-center gap-4">
          <Label for="width">Width</Label>
          <Input id="width" value="100%" class="col-span-2 h-8" />
        </div>
        <div class="grid grid-cols-3 items-center gap-4">
          <Label for="maxWidth">Max. width</Label>
          <Input id="maxWidth" value="300px" class="col-span-2 h-8" />
        </div>
        <div class="grid grid-cols-3 items-center gap-4">
          <Label for="height">Height</Label>
          <Input id="height" value="25px" class="col-span-2 h-8" />
        </div>
        <div class="grid grid-cols-3 items-center gap-4">
          <Label for="maxHeight">Max. height</Label>
          <Input id="maxHeight" value="none" class="col-span-2 h-8" />
        </div>
      </div>
    </div>
  </Popover.Content>
</Popover.Root>

<Tooltip.Provider>
  <Tooltip.Root>
    <Tooltip.Trigger class={buttonVariants({ variant: "outline" })}>
      Hover
    </Tooltip.Trigger>
    <Tooltip.Content>
      <p>Add to library</p>
    </Tooltip.Content>
  </Tooltip.Root>
</Tooltip.Provider>

<Pagination.Root count={100} perPage={10}>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.PrevButton />
      </Pagination.Item>
      {#each pages as page (page.key)}
        {#if page.type === "ellipsis"}
          <Pagination.Item>
            <Pagination.Ellipsis />
          </Pagination.Item>
        {:else}
          <Pagination.Item>
            <Pagination.Link {page} isActive={currentPage === page.value}>
              {page.value}
            </Pagination.Link>
          </Pagination.Item>
        {/if}
      {/each}
      <Pagination.Item>
        <Pagination.NextButton />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>
