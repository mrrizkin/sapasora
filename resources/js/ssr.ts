import { createInertiaApp } from '@inertiajs/vue3';
import createServer from '@inertiajs/vue3/server';
import ui from '@nuxt/ui/vue-plugin';
import { VueQueryPlugin } from '@tanstack/vue-query';
import { createHead, renderSSRHead } from '@unhead/vue/server';
import { renderToString } from '@vue/server-renderer';
import { resolvePageComponent } from 'laravel-vite-plugin/inertia-helpers';
import { type DefineComponent, createSSRApp, h } from 'vue';

const appName = import.meta.env.VITE_APP_NAME || 'Sapasora App';

createServer(async (page) => {
  const head = createHead();

  return createInertiaApp({
    page,
    title: (title) => (title ? `${title} - ${appName}` : appName),
    render: renderToString,
    resolve(name) {
      return resolvePageComponent(`./pages/${name}.vue`, import.meta.glob<DefineComponent>('./pages/**/*.vue'));
    },
    setup({ App, props, plugin }) {
      return createSSRApp({ render: () => h(App, props) })
        .use(ui)
        .use(VueQueryPlugin)
        .use(head)
        .use(plugin);
    },
  }).then(async (app) => {
    const payload = await renderSSRHead(head);
    app.head.push(payload.headTags);
    return app;
  });
});
