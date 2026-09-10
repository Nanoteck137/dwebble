<script>
  import {
    AlertCircleIcon,
    CheckCircle2Icon,
    CheckIcon,
    ChevronRight,
    FileTextIcon,
    FolderCodeIcon,
    MoonIcon,
    PopcornIcon,
    SearchIcon,
    SunIcon,
    XIcon,
  } from "@lucide/svelte";
  import { Button } from "$lib/components/ui/button";

  import {
    Accordion,
    Alert,
    AlertDialog,
    AspectRatio,
    Attachment,
    Avatar,
    Badge,
    Breadcrumb,
    Bubble,
    ButtonGroup,
    buttonVariants,
    Card,
    Carousel,
    Checkbox,
    Collapsible,
    Command,
    ContextMenu,
    Dialog,
    Drawer,
    DropdownMenu,
    Empty,
    Field,
    HoverCard,
    Input,
    InputGroup,
    InputOTP,
    Item,
    Kbd,
    Label,
    Marker,
    Menubar,
    Message,
    NativeSelect,
    NavigationMenu,
    Pagination,
    Popover,
    Progress,
    RadioGroup,
    Resizable,
    ScrollArea,
    Select,
    Separator,
    Sheet,
    Skeleton,
    Slider,
    Spinner,
    Switch,
    Table,
    Tabs,
    Textarea,
    Toggle,
    ToggleGroup,
    Tooltip,
  } from "$lib/components/ui";
  import { resetMode, setMode } from "mode-watcher";

  let value = $state(33);
</script>

<DropdownMenu.Root>
  <DropdownMenu.Trigger
    class={buttonVariants({ variant: "outline", size: "icon" })}
  >
    <SunIcon
      class="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 !transition-all dark:scale-0 dark:-rotate-90"
    />
    <MoonIcon
      class="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 !transition-all dark:scale-100 dark:rotate-0"
    />
    <span class="sr-only">Toggle theme</span>
  </DropdownMenu.Trigger>
  <DropdownMenu.Content align="end">
    <DropdownMenu.Item onclick={() => setMode("light")}
      >Light</DropdownMenu.Item
    >
    <DropdownMenu.Item onclick={() => setMode("dark")}>Dark</DropdownMenu.Item>
    <DropdownMenu.Item onclick={() => resetMode()}>System</DropdownMenu.Item>
  </DropdownMenu.Content>
</DropdownMenu.Root>

<div class="flex gap-2">
  <p class="w-25">Size: xs</p>
  <Button size="xs" variant="default">Default</Button>
  <Button size="xs" variant="outline">Outline</Button>
  <Button size="xs" variant="secondary">Secondary</Button>
  <Button size="xs" variant="ghost">Ghost</Button>
  <Button size="xs" variant="destructive">Destructive</Button>
  <Button size="xs" variant="link">Link</Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: sm</p>
  <Button size="sm" variant="default">Default</Button>
  <Button size="sm" variant="outline">Outline</Button>
  <Button size="sm" variant="secondary">Secondary</Button>
  <Button size="sm" variant="ghost">Ghost</Button>
  <Button size="sm" variant="destructive">Destructive</Button>
  <Button size="sm" variant="link">Link</Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: default</p>
  <Button size="default" variant="default">Default</Button>
  <Button size="default" variant="outline">Outline</Button>
  <Button size="default" variant="secondary">Secondary</Button>
  <Button size="default" variant="ghost">Ghost</Button>
  <Button size="default" variant="destructive">Destructive</Button>
  <Button size="default" variant="link">Link</Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: lg</p>
  <Button size="lg" variant="default">Default</Button>
  <Button size="lg" variant="outline">Outline</Button>
  <Button size="lg" variant="secondary">Secondary</Button>
  <Button size="lg" variant="ghost">Ghost</Button>
  <Button size="lg" variant="destructive">Destructive</Button>
  <Button size="lg" variant="link">Link</Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: icon</p>
  <Button size="icon" variant="default"><ChevronRight /></Button>
  <Button size="icon" variant="outline"><ChevronRight /></Button>
  <Button size="icon" variant="secondary"><ChevronRight /></Button>
  <Button size="icon" variant="ghost"><ChevronRight /></Button>
  <Button size="icon" variant="destructive"><ChevronRight /></Button>
  <Button size="icon" variant="link"><ChevronRight /></Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: icon-xs</p>
  <Button size="icon-xs" variant="default"><ChevronRight /></Button>
  <Button size="icon-xs" variant="outline"><ChevronRight /></Button>
  <Button size="icon-xs" variant="secondary"><ChevronRight /></Button>
  <Button size="icon-xs" variant="ghost"><ChevronRight /></Button>
  <Button size="icon-xs" variant="destructive"><ChevronRight /></Button>
  <Button size="icon-xs" variant="link"><ChevronRight /></Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: icon-sm</p>
  <Button size="icon-sm" variant="default"><ChevronRight /></Button>
  <Button size="icon-sm" variant="outline"><ChevronRight /></Button>
  <Button size="icon-sm" variant="secondary"><ChevronRight /></Button>
  <Button size="icon-sm" variant="ghost"><ChevronRight /></Button>
  <Button size="icon-sm" variant="destructive"><ChevronRight /></Button>
  <Button size="icon-sm" variant="link"><ChevronRight /></Button>
</div>

<div class="h-2"></div>

<div class="flex gap-2">
  <p class="w-25">Size: icon-lg</p>
  <Button size="icon-lg" variant="default"><ChevronRight /></Button>
  <Button size="icon-lg" variant="outline"><ChevronRight /></Button>
  <Button size="icon-lg" variant="secondary"><ChevronRight /></Button>
  <Button size="icon-lg" variant="ghost"><ChevronRight /></Button>
  <Button size="icon-lg" variant="destructive"><ChevronRight /></Button>
  <Button size="icon-lg" variant="link"><ChevronRight /></Button>
</div>

<Accordion.Root type="single" class="w-full sm:max-w-[70%]" value="item-1">
  <Accordion.Item value="item-1">
    <Accordion.Trigger>Product Information</Accordion.Trigger>
    <Accordion.Content class="flex flex-col gap-4 text-balance">
      <p>
        Our flagship product combines cutting-edge technology with sleek
        design. Built with premium materials, it offers unparalleled
        performance and reliability.
      </p>
      <p>
        Key features include advanced processing capabilities, and an intuitive
        user interface designed for both beginners and experts.
      </p>
    </Accordion.Content>
  </Accordion.Item>
  <Accordion.Item value="item-2">
    <Accordion.Trigger>Shipping Details</Accordion.Trigger>
    <Accordion.Content class="flex flex-col gap-4 text-balance">
      <p>
        We offer worldwide shipping through trusted courier partners. Standard
        delivery takes 3-5 business days, while express shipping ensures
        delivery within 1-2 business days.
      </p>
      <p>
        All orders are carefully packaged and fully insured. Track your
        shipment in real-time through our dedicated tracking portal.
      </p>
    </Accordion.Content>
  </Accordion.Item>
  <Accordion.Item value="item-3">
    <Accordion.Trigger>Return Policy</Accordion.Trigger>
    <Accordion.Content class="flex flex-col gap-4 text-balance">
      <p>
        We stand behind our products with a comprehensive 30-day return policy.
        If you&apos;re not completely satisfied, simply return the item in its
        original condition.
      </p>
      <p>
        Our hassle-free return process includes free return shipping and full
        refunds processed within 48 hours of receiving the returned item.
      </p>
    </Accordion.Content>
  </Accordion.Item>
</Accordion.Root>

<div class="grid w-full max-w-xl items-start gap-4">
  <Alert.Root>
    <CheckCircle2Icon />
    <Alert.Title>Success! Your changes have been saved</Alert.Title>
    <Alert.Description
      >This is an alert with icon, title and description.</Alert.Description
    >
  </Alert.Root>
  <Alert.Root>
    <PopcornIcon />
    <Alert.Title
      >This Alert has a title and an icon. No description.</Alert.Title
    >
  </Alert.Root>
  <Alert.Root variant="destructive">
    <AlertCircleIcon />
    <Alert.Title>Unable to process your payment.</Alert.Title>
    <Alert.Description>
      <p>Please verify your billing information and try again.</p>
      <ul class="list-inside list-disc text-sm">
        <li>Check your card details</li>
        <li>Ensure sufficient funds</li>
        <li>Verify billing address</li>
      </ul>
    </Alert.Description>
  </Alert.Root>
</div>

<AlertDialog.Root>
  <AlertDialog.Trigger class={buttonVariants({ variant: "outline" })}>
    Show Dialog
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

<AspectRatio ratio={16 / 9} class="rounded-lg bg-muted">
  <img
    src="https://images.unsplash.com/photo-1588345921523-c2dcdb7f1dcd?w=800&dpr=2&q=80"
    alt="Gray by Drew Beamer"
    class="h-full w-full rounded-lg object-cover dark:brightness-[0.2] dark:grayscale"
  />
</AspectRatio>

<Attachment.Root>
  <Attachment.Media>
    <FileTextIcon />
  </Attachment.Media>
  <Attachment.Content>
    <Attachment.Title>sales-dashboard.pdf</Attachment.Title>
    <Attachment.Description>PDF · 2.4 MB</Attachment.Description>
  </Attachment.Content>
  <Attachment.Actions>
    <Attachment.Action aria-label="Remove sales-dashboard.pdf">
      <XIcon />
    </Attachment.Action>
  </Attachment.Actions>
</Attachment.Root>

<Avatar.Root>
  <Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" />
  <Avatar.Fallback>CN</Avatar.Fallback>
</Avatar.Root>

<Badge variant="outline">Badge</Badge>

<Breadcrumb.Root>
  <Breadcrumb.List>
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/">Home</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Link href="/components">Components</Breadcrumb.Link>
    </Breadcrumb.Item>
    <Breadcrumb.Separator />
    <Breadcrumb.Item>
      <Breadcrumb.Page>Breadcrumb</Breadcrumb.Page>
    </Breadcrumb.Item>
  </Breadcrumb.List>
</Breadcrumb.Root>

<Bubble.Root>
  <Bubble.Content>
    I checked the registry output and removed the stale route.
  </Bubble.Content>
  <Bubble.Reactions>
    <span>👍</span>
  </Bubble.Reactions>
</Bubble.Root>

<ButtonGroup.Root>
  <Button>Button 1</Button>
  <Button>Button 2</Button>
</ButtonGroup.Root>

<!-- TODO(patrik): More examples -->
<Card.Root>
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

<!-- TODO(patrik): More examples -->
<Carousel.Root>
  <Carousel.Content>
    <Carousel.Item>...</Carousel.Item>
    <Carousel.Item>...</Carousel.Item>
    <Carousel.Item>...</Carousel.Item>
  </Carousel.Content>
  <Carousel.Previous />
  <Carousel.Next />
</Carousel.Root>

<Checkbox />

<Collapsible.Root>
  <Collapsible.Trigger>Can I use this in my project?</Collapsible.Trigger>
  <Collapsible.Content>
    Yes. Free to use for personal and commercial projects. No attribution
    required.
  </Collapsible.Content>
</Collapsible.Root>

<!-- TODO(patrik): Combobox Example -->

<Command.Root>
  <Command.Input placeholder="Type a command or search..." />
  <Command.List>
    <Command.Empty>No results found.</Command.Empty>
    <Command.Group heading="Suggestions">
      <Command.Item>Calendar</Command.Item>
      <Command.Item>Search Emoji</Command.Item>
      <Command.Item>Calculator</Command.Item>
    </Command.Group>
    <Command.Separator />
    <Command.Group heading="Settings">
      <Command.Item>Profile</Command.Item>
      <Command.Item>Billing</Command.Item>
      <Command.Item>Settings</Command.Item>
    </Command.Group>
  </Command.List>
</Command.Root>

<ContextMenu.Root>
  <ContextMenu.Trigger>Right click</ContextMenu.Trigger>
  <ContextMenu.Content>
    <ContextMenu.Item>Profile</ContextMenu.Item>
    <ContextMenu.Item>Billing</ContextMenu.Item>
    <ContextMenu.Item>Team</ContextMenu.Item>
    <ContextMenu.Item>Subscription</ContextMenu.Item>
  </ContextMenu.Content>
</ContextMenu.Root>

<Dialog.Root>
  <Dialog.Trigger>Open</Dialog.Trigger>
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

<Drawer.Root>
  <Drawer.Trigger>Open</Drawer.Trigger>
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
  <DropdownMenu.Trigger>Open</DropdownMenu.Trigger>
  <DropdownMenu.Content>
    <DropdownMenu.Group>
      <DropdownMenu.Label>My Account</DropdownMenu.Label>
      <DropdownMenu.Separator />
      <DropdownMenu.Item>Profile</DropdownMenu.Item>
      <DropdownMenu.Item>Billing</DropdownMenu.Item>
      <DropdownMenu.Item>Team</DropdownMenu.Item>
      <DropdownMenu.Item>Subscription</DropdownMenu.Item>
    </DropdownMenu.Group>
  </DropdownMenu.Content>
</DropdownMenu.Root>

<Empty.Root>
  <Empty.Header>
    <Empty.Media variant="icon">
      <FolderCodeIcon />
    </Empty.Media>
    <Empty.Title>No data</Empty.Title>
    <Empty.Description>No data found</Empty.Description>
  </Empty.Header>
  <Empty.Content>
    <Button>Add data</Button>
  </Empty.Content>
</Empty.Root>

<Field.Set>
  <Field.Legend>Profile</Field.Legend>
  <Field.Description>This appears on invoices and emails.</Field.Description>
  <Field.Group>
    <Field.Field>
      <Field.Label for="name">Full name</Field.Label>
      <Input id="name" autocomplete="off" placeholder="Evil Rabbit" />
      <Field.Description
        >This appears on invoices and emails.</Field.Description
      >
    </Field.Field>
    <Field.Field>
      <Field.Label for="username">Username</Field.Label>
      <Input id="username" autocomplete="off" aria-invalid />
      <Field.Error>Choose another username.</Field.Error>
    </Field.Field>
    <Field.Field orientation="horizontal">
      <Switch id="newsletter" />
      <Field.Label for="newsletter">Subscribe to the newsletter</Field.Label>
    </Field.Field>
  </Field.Group>
</Field.Set>

<HoverCard.Root>
  <HoverCard.Trigger>Hover</HoverCard.Trigger>
  <HoverCard.Content>
    SvelteKit - Web development, streamlined
  </HoverCard.Content>
</HoverCard.Root>

<Input type="email" placeholder="Email" class="max-w-xs" />

<InputGroup.Root>
  <InputGroup.Input placeholder="Search..." />
  <InputGroup.Addon>
    <SearchIcon />
  </InputGroup.Addon>
  <InputGroup.Addon align="inline-end">
    <InputGroup.Button>Search</InputGroup.Button>
  </InputGroup.Addon>
</InputGroup.Root>

<InputOTP.Root maxlength={6}>
  {#snippet children({ cells })}
    <InputOTP.Group>
      {#each cells.slice(0, 3) as cell}
        <InputOTP.Slot {cell} />
      {/each}
    </InputOTP.Group>
    <InputOTP.Separator />
    <InputOTP.Group>
      {#each cells.slice(3, 6) as cell}
        <InputOTP.Slot {cell} />
      {/each}
    </InputOTP.Group>
  {/snippet}
</InputOTP.Root>

<Item.Root>
  <Item.Header>Item Header</Item.Header>
  <Item.Media />
  <Item.Content>
    <Item.Title>Item</Item.Title>
    <Item.Description>Item</Item.Description>
  </Item.Content>
  <Item.Actions />
  <Item.Footer>Item Footer</Item.Footer>
</Item.Root>

<Kbd.Root>B</Kbd.Root>

<Label for="email">Your email address</Label>

<Marker.Root>
  <Marker.Icon>
    <CheckIcon />
  </Marker.Icon>
  <Marker.Content>Explored 4 files</Marker.Content>
</Marker.Root>

<Menubar.Root>
  <Menubar.Menu>
    <Menubar.Trigger>File</Menubar.Trigger>
    <Menubar.Content>
      <Menubar.Item>
        New Tab
        <Menubar.Shortcut>⌘T</Menubar.Shortcut>
      </Menubar.Item>
      <Menubar.Item>New Window</Menubar.Item>
      <Menubar.Separator />
      <Menubar.Item>Share</Menubar.Item>
      <Menubar.Separator />
      <Menubar.Item>Print</Menubar.Item>
    </Menubar.Content>
  </Menubar.Menu>
</Menubar.Root>

<Message.Root>
  <Message.Avatar>
    <Avatar.Root>
      <Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" />
      <Avatar.Fallback>CN</Avatar.Fallback>
    </Avatar.Root>
  </Message.Avatar>
  <Message.Content>
    <Bubble.Root>
      <Bubble.Content>How can I help you today?</Bubble.Content>
    </Bubble.Root>
  </Message.Content>
</Message.Root>

<NativeSelect.Root>
  <NativeSelect.Option value="">Select a fruit</NativeSelect.Option>
  <NativeSelect.Option value="apple">Apple</NativeSelect.Option>
  <NativeSelect.Option value="banana">Banana</NativeSelect.Option>
  <NativeSelect.Option value="blueberry">Blueberry</NativeSelect.Option>
  <NativeSelect.Option value="grapes" disabled>Grapes</NativeSelect.Option>
  <NativeSelect.Option value="pineapple">Pineapple</NativeSelect.Option>
</NativeSelect.Root>

<NavigationMenu.Root>
  <NavigationMenu.List>
    <NavigationMenu.Item>
      <NavigationMenu.Trigger>Item One</NavigationMenu.Trigger>
      <NavigationMenu.Content>
        <NavigationMenu.Link>Link</NavigationMenu.Link>
      </NavigationMenu.Content>
    </NavigationMenu.Item>
  </NavigationMenu.List>
</NavigationMenu.Root>

<Pagination.Root count={100} perPage={10}>
  {#snippet children({ pages, currentPage })}
    <Pagination.Content>
      <Pagination.Item>
        <Pagination.Previous />
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
        <Pagination.Next />
      </Pagination.Item>
    </Pagination.Content>
  {/snippet}
</Pagination.Root>

<Popover.Root>
  <Popover.Trigger>Open</Popover.Trigger>
  <Popover.Content>Place content for the popover here.</Popover.Content>
</Popover.Root>

<Progress value={33} />

<RadioGroup.Root value="option-one">
  <div class="flex items-center space-x-2">
    <RadioGroup.Item value="option-one" id="option-one" />
    <Label for="option-one">Option One</Label>
  </div>
  <div class="flex items-center space-x-2">
    <RadioGroup.Item value="option-two" id="option-two" />
    <Label for="option-two">Option Two</Label>
  </div>
</RadioGroup.Root>

<!-- TODO(patrik): Range Calender -->

<Resizable.PaneGroup direction="horizontal">
  <Resizable.Pane>One</Resizable.Pane>
  <Resizable.Handle />
  <Resizable.Pane>Two</Resizable.Pane>
</Resizable.PaneGroup>

<ScrollArea class="h-50 w-88 rounded-md border p-4">
  Jokester began sneaking into the castle in the middle of the night and
  leaving jokes all over the place: under the king's pillow, in his soup, even
  in the royal toilet. The king was furious, but he couldn't seem to stop
  Jokester. And then, one day, the people of the kingdom discovered that the
  jokes left by Jokester were so funny that they couldn't help but laugh. And
  once they started laughing, they couldn't stop.
</ScrollArea>

<Select.Root type="single">
  <Select.Trigger class="w-45"></Select.Trigger>
  <Select.Content>
    <Select.Item value="light">Light</Select.Item>
    <Select.Item value="dark">Dark</Select.Item>
    <Select.Item value="system">System</Select.Item>
  </Select.Content>
</Select.Root>

<Separator />

<Sheet.Root>
  <Sheet.Trigger>Open</Sheet.Trigger>
  <Sheet.Content>
    <Sheet.Header>
      <Sheet.Title>Are you sure absolutely sure?</Sheet.Title>
      <Sheet.Description>
        This action cannot be undone. This will permanently delete your account
        and remove your data from our servers.
      </Sheet.Description>
    </Sheet.Header>
  </Sheet.Content>
</Sheet.Root>

<!-- TODO(patrik): Missing "Sidebar" -->

<Skeleton class="h-5 w-25 rounded-full" />

<Slider type="single" bind:value max={100} step={1} />

<!-- TODO(patrik): Missing "Sonner" -->

<!-- TODO(patrik): Not working -->
<Spinner class="h-10 w-10" />

<Switch />

<Table.Root>
  <Table.Caption>A list of your recent invoices.</Table.Caption>
  <Table.Header>
    <Table.Row>
      <Table.Head class="w-25">Invoice</Table.Head>
      <Table.Head>Status</Table.Head>
      <Table.Head>Method</Table.Head>
      <Table.Head class="text-end">Amount</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    <Table.Row>
      <Table.Cell class="font-medium">INV001</Table.Cell>
      <Table.Cell>Paid</Table.Cell>
      <Table.Cell>Credit Card</Table.Cell>
      <Table.Cell class="text-end">$250.00</Table.Cell>
    </Table.Row>
  </Table.Body>
</Table.Root>

<Tabs.Root value="account" class="w-100">
  <Tabs.List>
    <Tabs.Trigger value="account">Account</Tabs.Trigger>
    <Tabs.Trigger value="password">Password</Tabs.Trigger>
  </Tabs.List>
  <Tabs.Content value="account">
    Make changes to your account here.
  </Tabs.Content>
  <Tabs.Content value="password">Change your password here.</Tabs.Content>
</Tabs.Root>

<Textarea />

<Toggle>Toggle</Toggle>

<ToggleGroup.Root type="single">
  <ToggleGroup.Item value="a">A</ToggleGroup.Item>
  <ToggleGroup.Item value="b">B</ToggleGroup.Item>
  <ToggleGroup.Item value="c">C</ToggleGroup.Item>
</ToggleGroup.Root>

<Tooltip.Root>
  <Tooltip.Trigger>Hover</Tooltip.Trigger>
  <Tooltip.Content>
    <p>Add to library</p>
  </Tooltip.Content>
</Tooltip.Root>

<div>
  <h1 class="scroll-m-20 text-4xl font-extrabold tracking-tight text-balance">
    Taxing Laughter: The Joke Tax Chronicles
  </h1>
  <p
    class="text-xl leading-7 text-muted-foreground [&:not(:first-child)]:mt-6"
  >
    Once upon a time, in a far-off land, there was a very lazy king who spent
    all day lounging on his throne. One day, his advisors came to him with a
    problem: the kingdom was running out of money.
  </p>
  <h2
    class="mt-10 scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight transition-colors first:mt-0"
  >
    The King&apos;s Plan
  </h2>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    The king thought long and hard, and finally came up with
    <a href="##" class="font-medium text-primary underline underline-offset-4">
      a brilliant plan
    </a>
    : he would tax the jokes in the kingdom.
  </p>
  <blockquote class="mt-6 border-s-2 ps-6 italic">
    &quot;After all,&quot; he said, &quot;everyone enjoys a good joke, so
    it&apos;s only fair that they should pay for the privilege.&quot;
  </blockquote>
  <h3 class="mt-8 scroll-m-20 text-2xl font-semibold tracking-tight">
    The Joke Tax
  </h3>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    The king&apos;s subjects were not amused. They grumbled and complained, but
    the king was firm:
  </p>
  <ul class="my-6 ms-6 list-disc [&>li]:mt-2">
    <li>1st level of puns: 5 gold coins</li>
    <li>2nd level of jokes: 10 gold coins</li>
    <li>3rd level of one-liners : 20 gold coins</li>
  </ul>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    As a result, people stopped telling jokes, and the kingdom fell into a
    gloom. But there was one person who refused to let the king&apos;s
    foolishness get him down: a court jester named Jokester.
  </p>
  <h3 class="mt-8 scroll-m-20 text-2xl font-semibold tracking-tight">
    Jokester&apos;s Revolt
  </h3>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    Jokester began sneaking into the castle in the middle of the night and
    leaving jokes all over the place: under the king&apos;s pillow, in his
    soup, even in the royal toilet. The king was furious, but he couldn&apos;t
    seem to stop Jokester.
  </p>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    And then, one day, the people of the kingdom discovered that the jokes left
    by Jokester were so funny that they couldn&apos;t help but laugh. And once
    they started laughing, they couldn&apos;t stop.
  </p>
  <h3 class="mt-8 scroll-m-20 text-2xl font-semibold tracking-tight">
    The People&apos;s Rebellion
  </h3>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    The people of the kingdom, feeling uplifted by the laughter, started to
    tell jokes and puns again, and soon the entire kingdom was in on the joke.
  </p>
  <div class="my-6 w-full overflow-y-auto">
    <table class="w-full">
      <thead>
        <tr class="m-0 border-t p-0 even:bg-muted">
          <th
            class="border px-4 py-2 text-start font-bold [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            King&apos;s Treasury
          </th>
          <th
            class="border px-4 py-2 text-start font-bold [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            People&apos;s happiness
          </th>
        </tr>
      </thead>
      <tbody>
        <tr class="m-0 border-t p-0 even:bg-muted">
          <td
            class="border px-4 py-2 text-start [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            Empty
          </td>
          <td
            class="border px-4 py-2 text-start [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            Overflowing
          </td>
        </tr>
        <tr class="m-0 border-t p-0 even:bg-muted">
          <td
            class="border px-4 py-2 text-start [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            Modest
          </td>
          <td
            class="border px-4 py-2 text-start [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            Satisfied
          </td>
        </tr>
        <tr class="m-0 border-t p-0 even:bg-muted">
          <td
            class="border px-4 py-2 text-start [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            Full
          </td>
          <td
            class="border px-4 py-2 text-start [&[align=center]]:text-center [&[align=right]]:text-end"
          >
            Ecstatic
          </td>
        </tr>
      </tbody>
    </table>
  </div>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    The king, seeing how much happier his subjects were, realized the error of
    his ways and repealed the joke tax. Jokester was declared a hero, and the
    kingdom lived happily ever after.
  </p>
  <p class="leading-7 [&:not(:first-child)]:mt-6">
    The moral of the story is: never underestimate the power of a good laugh
    and always be careful of bad ideas.
  </p>
</div>
