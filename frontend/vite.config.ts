import {defineConfig} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'
import * as path from "path";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      $lib: path.resolve("./src/lib"),
      $pages: path.resolve("./src/pages"),
      $wails: path.resolve("./wailsjs"),
      $components: path.resolve("./src/lib/components"),
      // $types: path.resolve("./src/lib/types"),
      $stores: path.resolve("./src/stores"),
    },
  },
  // resolve: {
  //   alias: {
  //     $lib: path.resolve("./src/lib"),
  //   },
  // },
})
