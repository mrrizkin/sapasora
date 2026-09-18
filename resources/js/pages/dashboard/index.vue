<script setup lang="ts">
import { Head } from '@inertiajs/vue3';

import DashboardLayout from '@/js/layout/dashboard.vue';
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
        </UDashboardNavbar>
      </template>

      <template #body></template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
