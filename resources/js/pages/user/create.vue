<script setup lang="ts">
import { Head, useForm } from '@inertiajs/vue3';
import type { FormSubmitEvent } from '@nuxt/ui';
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js';
import { reactive } from 'vue';

import DashboardLayout from '@/js/layout/dashboard.vue';
import AccountController from '@/js/lib/actions/app/http/controllers/account/AccountController';
import type { Role } from '@/js/types';

import { type UserFormSchema, userFormSchema } from './type';

interface Props {
  roles: Role[];
}

const props = defineProps<Props>();

const toast = useToast();

const roleOptions = props.roles.map((role) => ({ label: role.name, value: role.id }));

const state = reactive<Partial<UserFormSchema>>({
  name: '',
  username: '',
  password: '',
  role_id: roleOptions[0]?.value,
});

const inertiaForm = useForm({});

async function onSubmit(event: FormSubmitEvent<UserFormSchema>) {
  inertiaForm
    .transform(() => event.data)
    .post(AccountController.Store.url(), {
      preserveScroll: true,
      onSuccess: () => {
        toast.add({
          title: 'Success',
          description: 'User created successfully',
          icon: 'i-lucide-check',
          color: 'success',
        });
      },
      onError: () => {
        toast.add({
          title: 'Error',
          description: 'Failed to create user',
          icon: 'i-lucide-x',
          color: 'error',
        });
      },
    });
}
</script>

<template>
  <Head>
    <title>Create new User</title>
  </Head>
  <DashboardLayout>
    <UDashboardPanel id="user-new">
      <template #header>
        <UDashboardNavbar title="Create User" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
        </UDashboardNavbar>
        <UDashboardToolbar>
          <UButton label="Back" icon="i-lucide-arrow-left" variant="outline" color="neutral" :to="AccountController.Index.url()" />
        </UDashboardToolbar>
      </template>

      <template #body>
        <div class="mx-auto flex w-full flex-col gap-4 sm:gap-6 lg:max-w-2xl">
          <UForm id="user-form" :schema="userFormSchema" :state="state" @submit="onSubmit">
            <UPageCard
              title="Create New User"
              description="Add a new user account. They can log in and create their own API keys and device tokens."
              variant="naked"
              orientation="horizontal"
              class="mb-4">
              <div class="flex gap-3 lg:ms-auto">
                <UButton label="Cancel" variant="outline" color="neutral" to="/users" :disabled="inertiaForm.processing" />
                <UButton form="user-form" label="Create User" icon="i-lucide-plus" type="submit" :loading="inertiaForm.processing" />
              </div>
            </UPageCard>

            <UPageCard variant="subtle">
              <UFormField
                name="name"
                label="Full Name"
                description="Display name for this user"
                required
                class="flex items-start justify-between gap-4 max-sm:flex-col">
                <UInput v-model="state.name" placeholder="e.g., Jane Doe" autocomplete="off" :ui="{ base: 'w-[300px]' }" />
              </UFormField>

              <USeparator />

              <UFormField
                name="username"
                label="Username"
                description="Used to sign in"
                required
                class="flex items-start justify-between gap-4 max-sm:flex-col">
                <UInput v-model="state.username" placeholder="e.g., jane" autocomplete="off" :ui="{ base: 'w-[300px]' }" />
              </UFormField>

              <USeparator />

              <UFormField
                name="password"
                label="Password"
                description="Minimum 3 characters"
                required
                class="flex items-start justify-between gap-4 max-sm:flex-col">
                <UInput v-model="state.password" type="password" autocomplete="new-password" :ui="{ base: 'w-[300px]' }" />
              </UFormField>

              <USeparator />

              <UFormField
                name="role_id"
                label="Role"
                description="Determines what this user can access"
                required
                class="flex items-start justify-between gap-4 max-sm:flex-col">
                <USelectMenu v-model="state.role_id" :items="roleOptions" value-key="value" placeholder="Select role" :ui="{ base: 'w-[300px]' }" />
              </UFormField>
            </UPageCard>
          </UForm>
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
