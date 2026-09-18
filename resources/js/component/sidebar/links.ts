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
    {
      label: 'Phone book',
      to: '/phonebook',
      icon: 'i-lucide-book-user',
    },
    {
      label: 'Message History',
      to: '/messages',
      icon: 'i-lucide-message-square',
    },
    {
      label: 'Send',
      to: '/send',
      icon: 'i-lucide-send',
    },
    {
      label: 'Templates',
      to: '/templates',
      icon: 'i-lucide-book-dashed',
    },
    {
      label: 'Recurring',
      to: '/recurring',
      icon: 'i-lucide-calendar-sync',
    },
    {
      label: 'AutoReply',
      to: '/autoresponder',
      icon: 'i-lucide-reply-all',
    },
  ],

  // bottom sidebar menu
  [
    {
      label: 'Admins',
      to: '/admins',
      icon: 'i-lucide-user-star',
    },
    {
      label: 'Documentation',
      to: '/documentation',
      icon: 'i-lucide-book-open',
    },
  ],
] satisfies NavigationMenuItem[][];
