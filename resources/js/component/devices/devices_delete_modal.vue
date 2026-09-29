<script setup lang="ts">
import { router } from '@inertiajs/vue3';
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js';
import type { Row } from '@tanstack/table-core';
import { ref } from 'vue';

import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';
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

const toast = useToast();
const open = ref(false);
const submitting = ref(false);

async function onSubmit() {
  const ids = (props.selected ?? []).map((row) => row.original.id);
  if (ids.length === 0) {
    open.value = false;
    return;
  }

  submitting.value = true;

  try {
    await Promise.all(
      ids.map(
        (id) =>
          new Promise<void>((resolve, reject) => {
            router.delete(DeviceController.Destroy.url(id), {
              preserveScroll: true,
              onSuccess: () => resolve(),
              onError: () => reject(new Error(`Failed to delete device ${id}`)),
            });
          }),
      ),
    );

    toast.add({
      title: 'Device deleted',
      description: `${ids.length} device${ids.length > 1 ? 's' : ''} deleted successfully.`,
      icon: 'i-lucide-check',
      color: 'success',
    });

    router.reload();
  } catch {
    toast.add({
      title: 'Delete failed',
      description: 'One or more devices could not be deleted.',
      icon: 'i-lucide-x',
      color: 'error',
    });
  } finally {
    submitting.value = false;
    open.value = false;
  }
}
</script>

<template>
  <UModal v-model:open="open" :title="`Delete ${count} device${count > 1 ? 's' : ''}`" :description="`Are you sure, this action cannot be undone.`">
    <slot />

    <template #body>
      <div class="flex justify-end gap-2">
        <UButton label="Cancel" color="neutral" variant="subtle" :disabled="submitting" @click="open = false" />
        <UButton label="Delete" color="error" variant="solid" :loading="submitting" @click="onSubmit" />
      </div>
    </template>
  </UModal>
</template>
