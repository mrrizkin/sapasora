<script setup lang="ts">
import { computed } from 'vue';

import { links } from '@/js/component/sidebar/links';
import Sidebar from '@/js/component/sidebar/sidebar.vue';
import BaseLayout from '@/js/layout/base.vue';

const groups = computed(() => [
  {
    id: 'links',
    label: 'Go to',
    items: links.flat(),
  },
]);
</script>

<template>
  <BaseLayout>
    <UDashboardGroup unit="rem" storage="local">
      <Sidebar />

      <UDashboardSearch :groups="groups" />
      <slot />
    </UDashboardGroup>
  </BaseLayout>
</template>
