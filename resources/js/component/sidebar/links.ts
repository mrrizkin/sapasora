import type { NavigationMenuItem } from '@nuxt/ui';

export const links = [
  // top sidebar menu
  [
    {
      label: 'Dashboard',
      to: '/dashboard',
      icon: 'i-lucide-home',
    },
    {
      label: 'Devices',
      to: '/devices',
      icon: 'i-lucide-smartphone',
    },
  ],

  // bottom sidebar menu
  [
    {
      label: 'Users',
      to: '/users',
      icon: 'i-lucide-user-star',
    },
    {
      label: 'API Keys',
      to: '/api-keys',
      icon: 'i-lucide-key-round',
    },
  ],
] satisfies NavigationMenuItem[][];
