import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import vueDevTools from "vite-plugin-vue-devtools"
import tailwindcss from "@tailwindcss/vite"

import { fileURLToPath, URL } from "node:url"

export default defineConfig({
  plugins: [vue(), vueDevTools(), tailwindcss()],
  server: {
    hmr: false,
  },
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  build: {
    outDir: "dist",
    sourcemap: true,
    manifest: true,
  },
})
