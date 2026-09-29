<script setup lang="ts">
import { router, usePage } from '@inertiajs/vue3';
import type { DropdownMenuItem } from '@nuxt/ui';
import { computed } from 'vue';

import AuthController from '@/js/lib/actions/app/http/controllers/auth/AuthController';

defineProps<{
  collapsed?: boolean;
}>();

interface AuthPageProps {
  [key: string]: unknown;
  auth?: {
    user?: {
      id: string;
      name: string;
      username: string;
    };
  };
}

const page = usePage<AuthPageProps>();

const user = computed(() => ({
  name: page.props.auth?.user?.name ?? 'Unknown user',
  username: page.props.auth?.user?.username ?? '',
}));

function logout() {
  router.delete(AuthController.Logout.url(), {
    onSuccess: () => router.visit('/auth'),
  });
}

const items = computed<DropdownMenuItem[][]>(() => [
  [
    {
      type: 'label',
      label: user.value.name,
      icon: 'i-lucide-user',
    },
  ],
  [
    {
      label: 'Log out',
      icon: 'i-lucide-log-out',
      onSelect: logout,
    },
  ],
]);
</script>

<template>
  <UDropdownMenu
    :items="items"
    :content="{ align: 'center', collisionPadding: 12 }"
    :ui="{ content: collapsed ? 'w-48' : 'w-(--reka-dropdown-menu-trigger-width)' }">
    <UButton
      :label="collapsed ? undefined : user.name"
      icon="i-lucide-user-circle"
      :trailing-icon="collapsed ? undefined : 'i-lucide-chevrons-up-down'"
      color="neutral"
      variant="ghost"
      block
      :square="collapsed"
      class="data-[state=open]:bg-elevated"
      :ui="{
        trailingIcon: 'text-dimmed',
      }" />
  </UDropdownMenu>
</template>
