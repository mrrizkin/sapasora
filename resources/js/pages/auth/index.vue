<script setup lang="ts">
import { Head, router } from '@inertiajs/vue3';
import type { AuthFormField, FormSubmitEvent } from '@nuxt/ui';
import * as z from 'zod';

import BaseLayout from '@/js/layout/base.vue';
import AuthController from '@/js/lib/actions/app/http/controllers/auth/AuthController';

const fields: AuthFormField[] = [
  {
    name: 'username',
    type: 'text',
    label: 'Username',
    placeholder: 'Enter your username',
    required: true,
  },
  {
    name: 'password',
    label: 'Password',
    type: 'password',
    placeholder: 'Enter your password',
    required: true,
  },
  {
    name: 'remember',
    label: 'Remember me',
    type: 'checkbox',
  },
];

const schema = z.object({
  username: z.string('Invalid email').min(3, 'Must be at least 3 characters'),
  password: z.string('Password is required').min(3, 'Must be at least 3 characters'),
});

type Schema = z.output<typeof schema>;

function onSubmit(payload: FormSubmitEvent<Schema>) {
  router.post(AuthController.Login.url(), payload.data);
}
</script>

<template>
  <Head>
    <title>Login</title>
  </Head>
  <BaseLayout>
    <!-- center this card -->
    <div class="flex h-[100vh] flex-col items-center justify-center gap-4 p-4">
      <UPageCard class="w-full max-w-md">
        <UAuthForm :schema="schema" :fields="fields" title="Welcome back!" icon="i-lucide-lock" @submit="onSubmit">
          <template #validation>
            <UAlert color="error" icon="i-lucide-info" title="Error signing in" v-show="false" />
          </template>
        </UAuthForm>
      </UPageCard>
    </div>
  </BaseLayout>
</template>
