<script setup lang="ts">
import { useForm } from '@inertiajs/vue3';
import type { FormSubmitEvent } from '@nuxt/ui';
import { reactive } from 'vue';

import type { Device } from '@/js/types';

import { type DeviceFormSchema, deviceFormSchema } from '../type';

interface Props {
  isEdit: boolean;
  device?: Device;
  onSubmit: (event: FormSubmitEvent<DeviceFormSchema>) => void;
}

const props = defineProps<Props>();

// Type options
const typeOptions = [
  { label: 'WhatsApp', value: 'whatsapp' },
  { label: 'Telegram', value: 'telegram' },
  { label: 'SMS', value: 'sms' },
  { label: 'Email', value: 'email' },
];

// Form state
const device = reactive<Partial<DeviceFormSchema>>({
  name: props.device?.name,
  type: props.device?.type,
  webhook: props.device?.webhook,
});

// Inertia form for submission
const inertiaForm = useForm({});

// Submit handler
async function onSubmit(event: FormSubmitEvent<DeviceFormSchema>) {
  props.onSubmit(event);
}
</script>

<template>
  <UForm id="device-form" :schema="deviceFormSchema" :state="device" @submit="onSubmit">
    <UPageCard
      :title="props.isEdit ? 'Edit Device' : 'Create New Device'"
      description="Add a new device to your account and configure its settings."
      variant="naked"
      orientation="horizontal"
      class="mb-4">
      <div class="flex gap-3 lg:ms-auto">
        <UButton label="Cancel" variant="outline" color="neutral" to="/devices" :disabled="inertiaForm.processing" />
        <UButton
          form="device-form"
          :label="props.isEdit ? 'Edit Device' : 'Create New Device'"
          icon="i-lucide-plus"
          type="submit"
          :loading="inertiaForm.processing" />
      </div>
    </UPageCard>

    <UPageCard variant="subtle">
      <!-- Device Name -->
      <UFormField
        name="name"
        label="Device Name"
        description="A friendly name to identify this device"
        required
        class="flex items-start justify-between gap-4 max-sm:flex-col">
        <UInput
          v-model="device.name"
          placeholder="e.g., Production WhatsApp Bot"
          autocomplete="off"
          :disabled="inertiaForm.processing"
          :ui="{ base: 'w-[300px]' }" />
      </UFormField>

      <USeparator />

      <!-- Device Type -->
      <UFormField
        name="type"
        label="Device Type"
        description="Select the type of device you want to create"
        required
        class="flex items-start justify-between gap-4 max-sm:flex-col">
        <USelectMenu
          v-model="device.type"
          :items="typeOptions"
          value-key="value"
          placeholder="Select device type"
          :disabled="inertiaForm.processing"
          :ui="{ base: 'w-[300px]' }" />
      </UFormField>

      <USeparator />

      <!-- Webhook URL -->
      <UFormField
        name="webhook"
        label="Webhook URL"
        description="URL to receive event notifications (optional)"
        class="flex items-start justify-between gap-4 max-sm:flex-col">
        <UInput
          v-model="device.webhook"
          type="url"
          placeholder="https://your-domain.com/webhook"
          autocomplete="off"
          :disabled="inertiaForm.processing"
          :ui="{ base: 'w-[300px]' }" />
      </UFormField>
    </UPageCard>
  </UForm>
</template>
