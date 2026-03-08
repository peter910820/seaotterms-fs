import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import tsconfigPaths from "vite-tsconfig-paths";
import vuetify, { transformAssetUrls } from "vite-plugin-vuetify";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue({
      template: { transformAssetUrls },
    }),
    tsconfigPaths(),
    vuetify({ autoImport: true }),
  ],
  resolve: {
    alias: {
      "@": "/src",
    },
  },
  build: {
    rollupOptions: {
      output: { // 檔名完全使用hash
        entryFileNames: 'assets/[hash].js',
        chunkFileNames: 'assets/[hash].js',
        assetFileNames: 'assets/[hash].[ext]',
      }
    },
    minify: "terser", // 用 terser 取代預設的 esbuild
    terserOptions: {
      compress: {
        drop_console: false, // 是否刪除 console.log (跟你原本一樣)
        drop_debugger: true, // 刪除 debugger
      },
      output: {
        comments: false, // 移除註解
      },
      mangle: true, // 混淆變數名稱
    },
  },
});
