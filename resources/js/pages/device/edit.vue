<script setup lang="ts">
import { Head, useForm } from '@inertiajs/vue3';
import type { FormSubmitEvent } from '@nuxt/ui';

import DashboardLayout from '@/js/layout/dashboard.vue';
import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';
import DeviceForm from '@/js/pages/device/components/form.vue';
import type { Device } from '@/js/types';

import type { DeviceFormSchema } from './type';

interface Props {
  device?: Device;
}

const props = defineProps<Props>();

// Inertia form for submission
const inertiaForm = useForm({});

const toast = useToast();

// Submit handler
async function onSubmit(event: FormSubmitEvent<DeviceFormSchema>) {
  // Use Inertia to submit a PUT to the actual update endpoint.
  inertiaForm
    .transform(() => event.data)
    .put(DeviceController.Update.url(props.device?.id || '0'), {
      preserveScroll: true,
      onSuccess: () => {
        toast.add({
          title: 'Success',
          description: 'Device updated successfully',
          icon: 'i-lucide-check',
          color: 'success',
        });
      },
      onError: () => {
        toast.add({
          title: 'Error',
          description: 'Failed to update device',
          icon: 'i-lucide-x',
          color: 'error',
        });
      },
    });
}
</script>

<template>
  <Head>
    <title>Edit Device</title>
  </Head>
  <DashboardLayout>
    <UDashboardPanel id="device-edit">
      <template #header>
        <UDashboardNavbar title="Edit Device" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
        </UDashboardNavbar>
        <UDashboardToolbar>
          <UButton label="Back" icon="i-lucide-arrow-left" variant="outline" color="neutral" :to="DeviceController.Index.url()" />
        </UDashboardToolbar>
      </template>

      <template #body>
        <div class="mx-auto flex w-full flex-col gap-4 sm:gap-6 lg:max-w-2xl">
          <DeviceForm :isEdit="true" :onSubmit="onSubmit" :device="props.device" />
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
