import { createInertiaApp } from '@inertiajs/vue3';
import ui from '@nuxt/ui/vue-plugin';
import { VueQueryPlugin } from '@tanstack/vue-query';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import { type DefineComponent, createSSRApp, h } from 'vue';

import '@/css/app.css';

const appName = import.meta.env.VITE_APP_NAME || 'Sapasora App';

createInertiaApp({
  title: (title) => (title ? `${title} - ${appName}` : appName),
  resolve(name) {
    return resolvePageComponent(`./pages/${name}.vue`, import.meta.glob<DefineComponent>('./pages/**/*.vue'));
  },
  setup({ el, App, props, plugin }) {
    createSSRApp({ render: () => h(App, props) })
      .use(plugin)
      .use(ui)
      .use(VueQueryPlugin)
      .mount(el);
  },
});
