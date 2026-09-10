/** @type {import("prettier").Config} */
const config = {
  useTabs: false,
  singleQuote: false,
  trailingComma: "all",
  printWidth: 79,
  plugins: ["prettier-plugin-svelte", "prettier-plugin-tailwindcss"],
  overrides: [{ files: "*.svelte", options: { parser: "svelte" } }],
  tailwindStylesheet: "./src/app.css",
};

export default config;
