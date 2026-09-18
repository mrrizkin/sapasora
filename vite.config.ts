import ui from '@nuxt/ui/vite';
import tailwindcss from '@tailwindcss/vite';
import vue from '@vitejs/plugin-vue';
import laravel from 'laravel-vite-plugin';
import { URL, fileURLToPath } from 'node:url';
import { defineConfig } from 'vite';

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    laravel({
      input: ['resources/css/app.css', 'resources/js/app.ts', 'resources/js/colormode.ts'],
      ssr: 'resources/js/ssr.ts',
      ssrOutputDirectory: 'build/ssr',
    }),
    vue(),
    tailwindcss(),
    ui({
      inertia: true,
      ui: {
        colors: {
          primary: 'green',
          neutral: 'neutral',
        },
      },
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./resources', import.meta.url)),
    },
  },
});
