<script setup lang="ts">
import { Head } from '@inertiajs/vue3';
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js';
import { ref } from 'vue';

import DashboardLayout from '@/js/layout/dashboard.vue';
import APIKeyController from '@/js/lib/actions/app/http/controllers/apikey/APIKeyController';
import { request } from '@/js/lib/request';
import type { APIKey, Paginated } from '@/js/types';

interface Props {
  apikeys: Paginated<APIKey>;
}

const props = defineProps<Props>();

const toast = useToast();

const keys = ref<APIKey[]>(props.apikeys.data ?? []);
const newKeyName = ref('');
const creating = ref(false);
const newKey = ref<string | null>(null);

function copyToClipboard(value: string, label: string) {
  navigator.clipboard.writeText(value);
  toast.add({
    title: `${label} copied`,
    icon: 'i-lucide-check',
    color: 'success',
  });
}

async function reload() {
  try {
    const { data } = await request.get<{ data: APIKey[] }>(APIKeyController.List.url({ query: { limit: 100 } }));
    keys.value = data.data ?? [];
  } catch {
    // non-fatal
  }
}

async function createKey() {
  if (!newKeyName.value.trim()) {
    toast.add({ title: 'Name is required', icon: 'i-lucide-x', color: 'error' });
    return;
  }

  creating.value = true;
  try {
    const { data } = await request.post<APIKey>(APIKeyController.Store.url(), {
      name: newKeyName.value.trim(),
      permissions: {
        can_list_device: true,
        can_get_device: true,
        can_update_device: true,
        can_store_device: true,
        can_delete_device: true,
        can_list_device_token: true,
        can_get_device_token: true,
        can_store_device_token: true,
        can_delete_device_token: true,
      },
    });
    newKey.value = data.key;
    newKeyName.value = '';
    await reload();
  } catch (err) {
    toast.add({
      title: 'Failed to create API key',
      description: err instanceof Error ? err.message : 'Please try again.',
      icon: 'i-lucide-x',
      color: 'error',
    });
  } finally {
    creating.value = false;
  }
}

async function revokeKey(key: APIKey) {
  try {
    await request.delete(APIKeyController.Destroy.url(key.id));
    toast.add({ title: 'API key revoked', icon: 'i-lucide-check', color: 'success' });
    await reload();
  } catch {
    toast.add({ title: 'Failed to revoke key', icon: 'i-lucide-x', color: 'error' });
  }
}
</script>

<template>
  <Head>
    <title>API Keys</title>
  </Head>

  <DashboardLayout>
    <UDashboardPanel id="api-keys">
      <template #header>
        <UDashboardNavbar title="API Keys" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <div class="mx-auto flex w-full flex-col gap-4 sm:gap-6 lg:max-w-3xl">
          <UPageCard
            title="API Keys"
            description="API keys (sk-sak-...) let other apps manage your account (create devices, users, etc.) programmatically. They are NOT used to send WhatsApp messages — for that, generate a device token from a specific device's page instead."
            variant="subtle">
            <template #header>
              <div class="flex items-center gap-2">
                <UInput v-model="newKeyName" placeholder="Key name, e.g. 'CRM integration'" class="w-64" />
                <UButton label="Generate new key" icon="i-lucide-key-round" size="sm" :loading="creating" @click="createKey" />
              </div>
            </template>

            <UAlert
              v-if="newKey"
              color="warning"
              variant="subtle"
              icon="i-lucide-triangle-alert"
              title="Copy this key now — it won't be shown again"
              class="mb-4">
              <template #description>
                <div class="mt-2 flex items-center gap-2">
                  <code class="flex-1 truncate rounded bg-gray-100 px-2 py-1 text-xs dark:bg-gray-800">{{ newKey }}</code>
                  <UButton icon="i-lucide-copy" size="xs" variant="ghost" @click="copyToClipboard(newKey!, 'API key')" />
                </div>
              </template>
            </UAlert>

            <div v-if="keys.length === 0" class="py-4 text-center text-sm text-gray-500">No API keys yet. Generate one above.</div>
            <div v-else class="divide-y divide-gray-200 dark:divide-gray-800">
              <div v-for="key in keys" :key="key.id" class="flex items-center justify-between gap-4 py-3">
                <div class="flex flex-1 items-center gap-2">
                  <span class="text-sm font-medium text-gray-900 dark:text-white">{{ key.name }}</span>
                  <code class="rounded bg-gray-100 px-2 py-1 text-xs dark:bg-gray-800">{{ key.key.slice(0, 14) }}…</code>
                  <UBadge :color="key.status === 'active' ? 'success' : 'neutral'" variant="subtle" size="xs">{{ key.status }}</UBadge>
                </div>
                <UButton icon="i-lucide-trash-2" color="error" variant="ghost" size="xs" @click="revokeKey(key)" />
              </div>
            </div>
          </UPageCard>
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
