<div
  class="container mx-auto p-4 md:fixed md:bottom-1/2 md:left-1/2 md:-translate-x-[50%] md:translate-y-1/2"
>
  <div
    class="flex w-full flex-col overflow-hidden rounded shadow-md md:flex-row-reverse"
  >
    <div
      class="flex w-full flex-col items-center justify-center gap-4 bg-gradient-to-t from-[--primary] to-[--accent] p-4 md:bg-gradient-to-r md:p-6"
    >
      <h2 class="text-3xl font-black text-[--bg-color]">Login Page</h2>
      <button
        class="button-secondary-outline rounded-[50px] px-5 py-3 text-sm"
      >
        Create Account
      </button>
    </div>

    <div class="w-full bg-[--card-bg-color] p-4 text-[--card-fg-color] md:p-6">
      <p class="text-2xl font-light">Login</p>
      <div class="h-6"></div>
      <form class="flex flex-col gap-6" method="POST">
        <div class="flex flex-col gap-1">
          <label class="text-sm uppercase tracking-wide" for="username">
            Username
          </label>
          <input
            class="border-1 h-12 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
            id="username"
            name="username"
            type="text"
            placeholder="Username"
          />
        </div>

        <div class="flex flex-col gap-1">
          <label class="text-sm uppercase tracking-wide" for="password">
            Password
          </label>
          <input
            class="border-1 h-12 rounded-[50px] border-[--input-border-color] bg-[--input-bg-color] px-5 text-[--input-fg-color] placeholder:text-[--input-placeholder-color] focus:border-[--input-focus-border-color] focus:ring-0"
            id="password"
            name="password"
            type="password"
            placeholder="Password"
          />
        </div>

        <button
          class="button-primary rounded-[50px] px-5 py-3 text-sm"
          type="submit"
        >
          Login
        </button>
      </form>
    </div>
  </div>
</div>
