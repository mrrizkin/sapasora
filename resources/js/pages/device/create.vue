<script setup lang="ts">
import { Head, useForm } from '@inertiajs/vue3';
import type { FormSubmitEvent } from '@nuxt/ui';

import DashboardLayout from '@/js/layout/dashboard.vue';
import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';
import DeviceForm from '@/js/pages/device/components/form.vue';

import type { DeviceFormSchema } from './type';

// Inertia form for submission
const inertiaForm = useForm({});

const toast = useToast();

// Submit handler
async function onSubmit(event: FormSubmitEvent<DeviceFormSchema>) {
  // Use Inertia to submit
  inertiaForm
    .transform(() => event.data)
    .post(DeviceController.Store.url(), {
      preserveScroll: true,
      onSuccess: () => {
        toast.add({
          title: 'Success',
          description: 'Device created successfully',
          icon: 'i-lucide-check',
          color: 'success',
        });
      },
      onError: () => {
        toast.add({
          title: 'Error',
          description: 'Failed to create device',
          icon: 'i-lucide-x',
          color: 'error',
        });
      },
    });
}
</script>

<template>
  <Head>
    <title>Create new Device</title>
  </Head>
  <DashboardLayout>
    <UDashboardPanel id="device-new">
      <template #header>
        <UDashboardNavbar title="Create Device" :ui="{ right: 'gap-3' }">
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
          <DeviceForm :isEdit="false" :onSubmit="onSubmit" />
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
