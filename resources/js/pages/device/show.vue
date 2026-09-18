<script setup lang="ts">
import { Head } from '@inertiajs/vue3';
import { ref } from 'vue';

import DashboardLayout from '@/js/layout/dashboard.vue';
import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';
import type { Device } from '@/js/types';

interface Props {
  device: Device;
}

const props = defineProps<Props>();

const toast = useToast();

const showQrCode = ref(false);

// Format date helper
const formatDate = (date: Date) => {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

// Get status color
const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    active: 'success',
    inactive: 'neutral',
    pending: 'warning',
    error: 'error',
    disconnected: 'error',
  };
  return colors[status.toLowerCase()] || 'neutral';
};

// Get type icon
const getTypeIcon = (type: string) => {
  const icons: Record<string, string> = {
    whatsapp: 'i-lucide-message-circle',
    telegram: 'i-lucide-send',
    sms: 'i-lucide-smartphone',
    email: 'i-lucide-mail',
  };
  return icons[type.toLowerCase()] || 'i-lucide-radio';
};

// Delete device
const deleteDevice = () => {
  toast.add({
    title: 'Device deleted',
    description: 'The device has been deleted.',
  });
};

// Copy to clipboard
const copyToClipboard = (text: string, label: string) => {
  navigator.clipboard.writeText(text);
  toast.add({
    title: 'Copied',
    description: `${label} copied to clipboard`,
    icon: 'i-lucide-check',
    color: 'success',
  });
};

function shouldShowQRCode() {
  if (props.device.qr_code.length === 0) {
    toast.add({
      title: 'No QR Code',
      description: 'This device does not have a QR Code',
      icon: 'i-lucide-x',
      color: 'error',
    });

    return;
  }

  showQrCode.value = true;
}
</script>

<template>
  <Head>
    <title>{{ device.name }} - Device Details</title>
  </Head>
  <DashboardLayout>
    <UDashboardPanel id="device-show">
      <template #header>
        <UDashboardNavbar :title="device.name" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
        </UDashboardNavbar>
        <UDashboardToolbar>
          <UButton label="Back" icon="i-lucide-arrow-left" variant="outline" color="neutral" :to="DeviceController.Index.url()" />
          <div class="flex items-center gap-3">
            <UButton label="Edit" icon="i-lucide-pencil" variant="outline" color="neutral" :to="DeviceController.Edit.url(device.id)" />
            <UButton label="Delete" icon="i-lucide-trash-2" variant="outline" color="error" @click="deleteDevice" />
          </div>
        </UDashboardToolbar>
      </template>

      <template #body>
        <div class="mx-auto flex w-full flex-col gap-4 sm:gap-6 lg:max-w-4xl">
          <!-- Status Badge -->
          <div class="flex items-center justify-between gap-3">
            <div class="flex items-center gap-3">
              <UBadge :color="getStatusColor(device.status)" variant="subtle" size="lg" :label="device.status" />
              <UIcon :name="getTypeIcon(device.type)" class="h-5 w-5 text-gray-500 dark:text-gray-400" />
              <span class="text-sm text-gray-500 dark:text-gray-400">
                {{ device.type }}
              </span>
            </div>
            <div class="flex items-center gap-3">
              <UButton icon="i-lucide-scan-qr-code" label="Scan QR Code" variant="outline" color="neutral" size="sm" @click="shouldShowQRCode" />
              <UDrawer v-model:open="showQrCode" :ui="{ container: 'max-w-xl mx-auto' }">
                <template #content>
                  <div class="p-4">
                    <h2 class="mb-4 text-lg font-semibold">Device QR Code</h2>
                    <div class="flex justify-center" v-show="props.device.qr_code.length > 0">
                      <img :src="props.device.qr_code" alt="QR Code" class="h-auto w-full max-w-xs" />
                    </div>
                    <p class="mt-4 text-center text-sm text-gray-500" v-show="props.device.qr_code.length > 0">
                      Scan this code to connect your device
                    </p>
                  </div>
                </template>
              </UDrawer>
            </div>
          </div>

          <!-- Device Information -->
          <UPageCard title="Device Information" description="Basic information about this device" variant="subtle">
            <!-- Device ID -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Device ID</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Unique identifier for this device</p>
              </div>
              <div class="flex items-center gap-2">
                <code class="rounded bg-gray-100 px-2 py-1 text-sm dark:bg-gray-800">
                  {{ device.id }}
                </code>
                <UButton icon="i-lucide-copy" variant="ghost" color="gray" size="xs" @click="copyToClipboard(device.id, 'Device ID')" />
              </div>
            </div>

            <USeparator />

            <!-- Device Name -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Device Name</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Display name for this device</p>
              </div>
              <p class="text-sm text-gray-900 dark:text-white">
                {{ device.name }}
              </p>
            </div>

            <USeparator />

            <!-- Device Type -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Device Type</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">The type of messaging platform</p>
              </div>
              <UBadge variant="subtle" color="gray">
                {{ device.type }}
              </UBadge>
            </div>

            <USeparator />

            <!-- JID -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">JID (Jabber ID)</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Unique messaging identifier</p>
              </div>
              <div v-if="device.jid" class="flex items-center gap-2">
                <code class="rounded bg-gray-100 px-2 py-1 text-sm dark:bg-gray-800">
                  {{ device.jid }}
                </code>
                <UButton icon="i-lucide-copy" variant="ghost" color="gray" size="xs" @click="copyToClipboard(device.jid, 'JID')" />
              </div>
              <span v-else class="text-sm text-gray-400 dark:text-gray-500">Not configured</span>
            </div>

            <USeparator />

            <!-- Created At -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Created</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">When this device was created</p>
              </div>
              <p class="text-sm text-gray-900 dark:text-white">
                {{ formatDate(device.created_at) }}
              </p>
            </div>

            <USeparator />

            <!-- Updated At -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Last Updated</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">When this device was last modified</p>
              </div>
              <p class="text-sm text-gray-900 dark:text-white">
                {{ formatDate(device.updated_at) }}
              </p>
            </div>
          </UPageCard>

          <!-- Webhook Configuration -->
          <UPageCard title="Webhook Configuration" description="Event notification settings" variant="subtle">
            <!-- Webhook URL -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Webhook URL</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Endpoint for receiving events</p>
              </div>
              <div v-if="device.webhook" class="flex items-center gap-2">
                <code class="max-w-xs truncate rounded bg-gray-100 px-2 py-1 text-sm dark:bg-gray-800">
                  {{ device.webhook }}
                </code>
                <UButton icon="i-lucide-copy" variant="ghost" color="gray" size="xs" @click="copyToClipboard(device.webhook, 'Webhook URL')" />
                <UButton icon="i-lucide-external-link" variant="ghost" color="gray" size="xs" :to="device.webhook" target="_blank" />
              </div>
              <span v-else class="text-sm text-gray-400 dark:text-gray-500">Not configured</span>
            </div>

            <USeparator />

            <!-- Events -->
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Events</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Subscribed webhook events</p>
              </div>
              <div v-if="device.events" class="flex flex-wrap gap-2">
                <UBadge v-for="event in device.events.split(',')" :key="event" variant="subtle" color="blue">
                  {{ event.trim() }}
                </UBadge>
              </div>
              <span v-else class="text-sm text-gray-400 dark:text-gray-500">No events configured</span>
            </div>
          </UPageCard>

          <!-- Owner Information -->
          <UPageCard title="Owner Information" description="User who created this device" variant="subtle">
            <div class="flex items-start justify-between gap-4 max-sm:flex-col">
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">Owner</p>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Account that owns this device</p>
              </div>
              <div class="flex items-center gap-3">
                <UAvatar :alt="device.user.name" size="sm" />
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">{{ device.user.name }}</p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">@{{ device.user.username }}</p>
                </div>
              </div>
            </div>
          </UPageCard>
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
