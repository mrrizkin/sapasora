<script setup lang="ts">
import type { DropdownMenuItem } from '@nuxt/ui';
import { computed, ref } from 'vue';

defineProps<{
  collapsed?: boolean;
}>();

const teams = ref([
  {
    label: 'Vue',
    avatar: { src: 'https://avatars.githubusercontent.com/u/6128107?v=4', alt: 'Vue', crossOrigin: 'anonymous' },
  },
  {
    label: 'Vite',
    avatar: { src: 'https://avatars.githubusercontent.com/u/65625612?v=4', alt: 'Vite', crossOrigin: 'anonymous' },
  },
  {
    label: 'Vitest',
    avatar: { src: 'https://avatars.githubusercontent.com/u/95747107?v=4', alt: 'Vitest', crossOrigin: 'anonymous' },
  },
]);
const selectedTeam = ref(teams.value[0]);

const items = computed<DropdownMenuItem[][]>(() => {
  return [
    teams.value.map((team) => ({
      ...team,
      onSelect() {
        selectedTeam.value = team;
      },
    })),
    [
      {
        label: 'Create team',
        icon: 'i-lucide-circle-plus',
      },
      {
        label: 'Manage teams',
        icon: 'i-lucide-cog',
      },
    ],
  ];
});
</script>

<template>
  <UDropdownMenu
    :items="items"
    :content="{ align: 'center', collisionPadding: 12 }"
    :ui="{ content: collapsed ? 'w-40' : 'w-(--reka-dropdown-menu-trigger-width)' }">
    <UButton
      v-bind="{
        ...selectedTeam,
        label: collapsed ? undefined : selectedTeam?.label,
        trailingIcon: collapsed ? undefined : 'i-lucide-chevrons-up-down',
      }"
      color="neutral"
      variant="ghost"
      block
      :square="collapsed"
      class="data-[state=open]:bg-elevated"
      :class="[!collapsed && 'py-2']"
      :ui="{
        trailingIcon: 'text-dimmed',
      }" />
  </UDropdownMenu>
</template>
