<script setup lang="ts">
import { Head, router } from '@inertiajs/vue3';
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js';
import { onBeforeUnmount, reactive, ref } from 'vue';

import DashboardLayout from '@/js/layout/dashboard.vue';
import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';
import DeviceTokenController from '@/js/lib/actions/app/http/controllers/devicetoken/DeviceTokenController';
import { request } from '@/js/lib/request';
import type { Device, DeviceToken } from '@/js/types';

interface Props {
  device: Device;
}

const props = defineProps<Props>();

const toast = useToast();

// Local, mutable copy of the device so we can update qr_code/status in place
// while polling, without a full Inertia page reload.
const device = reactive<Device>({ ...props.device });

const showQrCode = ref(false);
const connecting = ref(false);
const deleting = ref(false);

let pollTimer: ReturnType<typeof setInterval> | null = null;
let pollDeadline = 0;

function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer);
    pollTimer = null;
  }
  connecting.value = false;
}

async function pollDeviceStatus() {
  try {
    const { data } = await request.get<Device>(DeviceController.Get.url(device.id));
    device.status = data.status;
    device.qr_code = data.qr_code;
    device.jid = data.jid;

    if (data.qr_code && !showQrCode.value) {
      showQrCode.value = true;
    }

    if (data.status === 'connected') {
      stopPolling();
      showQrCode.value = false;
      toast.add({
        title: 'Device connected',
        description: 'The device paired successfully.',
        icon: 'i-lucide-check',
        color: 'success',
      });
      return;
    }
  } catch {
    // Keep polling on transient errors; the deadline below caps total time.
  }

  if (Date.now() > pollDeadline) {
    stopPolling();
    toast.add({
      title: 'Timed out',
      description: 'No response from the device within the expected time. Try again.',
      icon: 'i-lucide-clock-alert',
      color: 'warning',
    });
  }
}

async function connectDevice() {
  connecting.value = true;
  try {
    await request.post(DeviceController.Connect.url(device.id));
    toast.add({
      title: 'Connecting…',
      description: 'Waiting for the QR code to be generated.',
      icon: 'i-lucide-loader',
    });

    pollDeadline = Date.now() + 60_000;
    pollTimer = setInterval(pollDeviceStatus, 1500);
  } catch (err) {
    connecting.value = false;
    toast.add({
      title: 'Connect failed',
      description: err instanceof Error ? err.message : 'Could not start the connection.',
      icon: 'i-lucide-x',
      color: 'error',
    });
  }
}

onBeforeUnmount(() => stopPolling());

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
    connected: 'success',
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
function deleteDevice() {
  deleting.value = true;
  router.delete(DeviceController.Destroy.url(device.id), {
    onSuccess: () => {
      toast.add({
        title: 'Device deleted',
        description: 'The device has been deleted.',
        icon: 'i-lucide-check',
        color: 'success',
      });
      router.visit(DeviceController.Index.url());
    },
    onError: () => {
      deleting.value = false;
      toast.add({
        title: 'Delete failed',
        description: 'Could not delete this device.',
        icon: 'i-lucide-x',
        color: 'error',
      });
    },
  });
}

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
  if (device.qr_code.length === 0) {
    toast.add({
      title: 'No QR Code',
      description: 'This device does not have a QR Code yet. Click "Connect" to generate one.',
      icon: 'i-lucide-x',
      color: 'error',
    });

    return;
  }

  showQrCode.value = true;
}

// --- Device tokens (used by other apps to authenticate against /api/v1/gateway/*) ---
const tokens = ref<DeviceToken[]>([]);
const loadingTokens = ref(false);
const creatingToken = ref(false);
const newToken = ref<string | null>(null);

async function loadTokens() {
  loadingTokens.value = true;
  try {
    const { data } = await request.get<{ data: DeviceToken[] }>(
      DeviceTokenController.List.url({ query: { device_id: device.id, limit: 100 } }),
    );
    tokens.value = data.data ?? [];
  } catch {
    // non-fatal, the section will just show empty state
  } finally {
    loadingTokens.value = false;
  }
}

async function createToken() {
  creatingToken.value = true;
  try {
    const { data } = await request.post<DeviceToken>(DeviceTokenController.Store.url(), {
      device_id: device.id,
    });
    newToken.value = data.token;
    await loadTokens();
  } catch (err) {
    toast.add({
      title: 'Failed to create token',
      description: err instanceof Error ? err.message : 'Please try again.',
      icon: 'i-lucide-x',
      color: 'error',
    });
  } finally {
    creatingToken.value = false;
  }
}

async function revokeToken(token: DeviceToken) {
  try {
    await request.delete(DeviceTokenController.Destroy.url(token.id));
    toast.add({
      title: 'Token revoked',
      icon: 'i-lucide-check',
      color: 'success',
    });
    await loadTokens();
  } catch {
    toast.add({
      title: 'Failed to revoke token',
      icon: 'i-lucide-x',
      color: 'error',
    });
  }
}

loadTokens();

// --- Test send message (only relevant once the device is connected) ---
const testPhone = ref('');
const testBody = ref('Hello from Sapasora! This is a test message.');
const sendingTest = ref(false);
const testResult = ref<{ success: boolean; message: string } | null>(null);

async function sendTestMessage() {
  if (!testPhone.value.trim() || !testBody.value.trim()) {
    toast.add({ title: 'Phone number and message are required', icon: 'i-lucide-x', color: 'error' });
    return;
  }

  sendingTest.value = true;
  testResult.value = null;
  try {
    await request.post(DeviceController.SendTestMessage.url(device.id), {
      phone: testPhone.value.trim(),
      body: testBody.value.trim(),
    });
    testResult.value = { success: true, message: 'Message sent successfully.' };
    toast.add({ title: 'Test message sent', icon: 'i-lucide-check', color: 'success' });
  } catch (err) {
    // request.ts's response interceptor throws a Fetch API Response carrying
    // the backend's { title, message } JSON body, not a plain Error.
    let description = 'Please try again.';
    if (err instanceof Response) {
      try {
        const body = await err.clone().json();
        description = body?.message ?? description;
      } catch {
        // body wasn't JSON, fall back to the generic message
      }
    } else if (err instanceof Error) {
      description = err.message;
    }
    testResult.value = { success: false, message: description };
    toast.add({
      title: 'Failed to send test message',
      description,
      icon: 'i-lucide-x',
      color: 'error',
    });
  } finally {
    sendingTest.value = false;
  }
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
            <UButton label="Delete" icon="i-lucide-trash-2" variant="outline" color="error" :loading="deleting" @click="deleteDevice" />
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
              <UButton
                icon="i-lucide-plug-zap"
                :label="connecting ? 'Connecting…' : 'Connect'"
                color="primary"
                size="sm"
                :loading="connecting"
                @click="connectDevice" />
              <UButton icon="i-lucide-scan-qr-code" label="Scan QR Code" variant="outline" color="neutral" size="sm" @click="shouldShowQRCode" />
              <UDrawer v-model:open="showQrCode" :ui="{ container: 'max-w-xl mx-auto' }">
                <template #content>
                  <div class="p-4">
                    <h2 class="mb-4 text-lg font-semibold">Device QR Code</h2>
                    <div class="flex justify-center" v-show="device.qr_code.length > 0">
                      <img :src="device.qr_code" alt="QR Code" class="h-auto w-full max-w-xs" />
                    </div>
                    <p class="mt-4 text-center text-sm text-gray-500" v-show="device.qr_code.length > 0">
                      Scan this code with WhatsApp (Linked devices → Link a device) to pair.
                    </p>
                    <p class="mt-4 text-center text-sm text-gray-500" v-show="device.qr_code.length === 0 && connecting">
                      Waiting for the QR code to be generated…
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

          <!-- API Access / Device Tokens -->
          <UPageCard
            title="API Access"
            description="Device tokens are used by other apps to call the gateway API (send messages, check status, etc.) on behalf of this device."
            variant="subtle">
            <template #header>
              <UButton
                label="Generate new token"
                icon="i-lucide-key-round"
                size="sm"
                :loading="creatingToken"
                @click="createToken" />
            </template>

            <UAlert
              v-if="newToken"
              color="warning"
              variant="subtle"
              icon="i-lucide-triangle-alert"
              title="Copy this token now — it won't be shown again"
              class="mb-4">
              <template #description>
                <div class="mt-2 flex items-center gap-2">
                  <code class="flex-1 truncate rounded bg-gray-100 px-2 py-1 text-xs dark:bg-gray-800">{{ newToken }}</code>
                  <UButton icon="i-lucide-copy" size="xs" variant="ghost" @click="copyToClipboard(newToken!, 'Device token')" />
                </div>
              </template>
            </UAlert>

            <div v-if="loadingTokens" class="py-4 text-center text-sm text-gray-500">Loading tokens…</div>
            <div v-else-if="tokens.length === 0" class="py-4 text-center text-sm text-gray-500">
              No device tokens yet. Generate one to let another app call the gateway for this device.
            </div>
            <div v-else class="divide-y divide-gray-200 dark:divide-gray-800">
              <div v-for="token in tokens" :key="token.id" class="flex items-center justify-between gap-4 py-3">
                <div class="flex items-center gap-2">
                  <code class="rounded bg-gray-100 px-2 py-1 text-xs dark:bg-gray-800">{{ token.token.slice(0, 14) }}…</code>
                  <UBadge :color="token.status === 'active' ? 'success' : 'neutral'" variant="subtle" size="xs">{{ token.status }}</UBadge>
                </div>
                <UButton icon="i-lucide-trash-2" color="error" variant="ghost" size="xs" @click="revokeToken(token)" />
              </div>
            </div>
          </UPageCard>

          <!-- Test Send Message -->
          <UPageCard
            v-if="device.status === 'connected'"
            title="Test Send Message"
            description="Send a quick test message from this device to verify the pairing actually works."
            variant="subtle">
            <div class="flex flex-col gap-3">
              <UFormField label="Recipient phone number" description="Include country code, e.g. 62812xxxxxxx">
                <UInput v-model="testPhone" placeholder="62812xxxxxxx" class="w-full" />
              </UFormField>
              <UFormField label="Message">
                <UTextarea v-model="testBody" :rows="3" class="w-full" />
              </UFormField>
              <div>
                <UButton label="Send test message" icon="i-lucide-send" :loading="sendingTest" @click="sendTestMessage" />
              </div>

              <UAlert
                v-if="testResult"
                :color="testResult.success ? 'success' : 'error'"
                variant="subtle"
                :icon="testResult.success ? 'i-lucide-check' : 'i-lucide-triangle-alert'"
                :title="testResult.message" />
            </div>
          </UPageCard>
          <UPageCard
            v-else
            title="Test Send Message"
            description="Connect this device first to send a test message."
            variant="subtle" />

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
