<script setup lang="ts">
import { Head } from '@inertiajs/vue3';

import DashboardLayout from '@/js/layout/dashboard.vue';
import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';

interface Stats {
  total_devices: number;
  connected_devices: number;
  disconnected_devices: number;
  inactive_devices: number;
  total_api_keys: number;
}

interface Props {
  stats: Stats;
}

const props = defineProps<Props>();

const cards = [
  {
    label: 'Your Devices',
    value: () => props.stats.total_devices,
    icon: 'i-lucide-smartphone',
    color: 'primary',
  },
  {
    label: 'Connected',
    value: () => props.stats.connected_devices,
    icon: 'i-lucide-plug-zap',
    color: 'success',
  },
  {
    label: 'Disconnected / Inactive',
    value: () => props.stats.disconnected_devices + props.stats.inactive_devices,
    icon: 'i-lucide-plug-2',
    color: 'neutral',
  },
  {
    label: 'Your API Keys',
    value: () => props.stats.total_api_keys,
    icon: 'i-lucide-key-round',
    color: 'warning',
  },
] as const;
</script>

<template>
  <Head>
    <title>Dashboard</title>
  </Head>

  <DashboardLayout>
    <UDashboardPanel id="dashboard">
      <template #header>
        <UDashboardNavbar title="Dashboard" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
          <template #right>
            <UButton label="New device" icon="i-lucide-plus" :to="DeviceController.Create.url()" />
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <UPageCard v-for="card in cards" :key="card.label" variant="subtle" class="gap-2">
            <div class="flex items-center gap-3">
              <div class="rounded-lg bg-gray-100 p-2 dark:bg-gray-800">
                <UIcon :name="card.icon" class="h-5 w-5" />
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ card.label }}</p>
                <p class="text-2xl font-semibold text-gray-900 dark:text-white">{{ card.value() }}</p>
              </div>
            </div>
          </UPageCard>
        </div>

        <UPageCard v-if="stats.total_devices === 0" title="Get started" description="You don't have any devices yet." variant="subtle" class="mt-4">
          <UButton label="Create your first device" icon="i-lucide-plus" :to="DeviceController.Create.url()" />
        </UPageCard>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
