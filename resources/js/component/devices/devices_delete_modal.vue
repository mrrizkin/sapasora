<script setup lang="ts">
import type { Row } from '@tanstack/table-core';
import { ref } from 'vue';

import type { Device } from '@/js/types';

const props = withDefaults(
  defineProps<{
    count?: number;
    selected?: Row<Device>[];
  }>(),
  {
    count: 0,
  },
);

const open = ref(false);

async function onSubmit() {
  console.log('delete', props.selected);
  await new Promise((resolve) => setTimeout(resolve, 1000));
  open.value = false;
}
</script>

<template>
  <UModal v-model:open="open" :title="`Delete ${count} device${count > 1 ? 's' : ''}`" :description="`Are you sure, this action cannot be undone.`">
    <slot />

    <template #body>
      <div class="flex justify-end gap-2">
        <UButton label="Cancel" color="neutral" variant="subtle" @click="open = false" />
        <UButton label="Delete" color="error" variant="solid" loading-auto @click="onSubmit" />
      </div>
    </template>
  </UModal>
</template>
