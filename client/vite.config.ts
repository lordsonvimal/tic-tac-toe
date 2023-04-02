import { defineConfig } from "vite";
import solidPlugin from "vite-plugin-solid";
import devtools from "solid-devtools/vite";

export default defineConfig({
  plugins: [devtools({
    autoname: true,
    locator: {
      componentLocation: true,
      jsxLocation: true
    }}),
    solidPlugin()],
  build: {
    target: "esnext"
  },
});
