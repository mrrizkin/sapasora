<script setup lang="ts">
import { Head, router } from '@inertiajs/vue3';
import type { TableColumn } from '@nuxt/ui';
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js';
import type { PaginationState } from '@tanstack/table-core';
import { computed, h, ref, resolveComponent, watch } from 'vue';

import DashboardLayout from '@/js/layout/dashboard.vue';
import AccountController from '@/js/lib/actions/app/http/controllers/account/AccountController';
import { debounce } from '@/js/lib/utils';
import type { Account, Paginated } from '@/js/types';

import type { ColumnFilter } from './type';

interface Props {
  users: Paginated<Account>;
  filters: ColumnFilter;
  pagination: PaginationState;
}

const props = defineProps<Props>();

const toast = useToast();

const state = ref<{
  columnFilters: ColumnFilter;
  pagination: PaginationState;
}>({
  columnFilters: props.filters,
  pagination: props.pagination,
});

watch(
  () => [state.value.columnFilters, state.value.pagination] as const,
  () => {
    router.visit(
      AccountController.Index.url({
        query: {
          page: state.value.pagination.pageIndex + 1,
          limit: state.value.pagination.pageSize,
          ...state.value.columnFilters,
        },
      }),
      {
        preserveState: true,
        preserveScroll: true,
        replace: true,
      },
    );
  },
  { deep: true },
);

const updateFilters = debounce((filters: ColumnFilter) => {
  if (state.value.columnFilters.search !== filters.search) {
    state.value.columnFilters.search = filters.search;
    state.value.pagination.pageIndex = 0;
  }
}, 500);

const searchQuery = computed({
  get: () => state.value.columnFilters.search || '',
  set: (value: string) => {
    updateFilters({ ...state.value.columnFilters, search: value });
  },
});

function deleteUser(user: Account) {
  router.delete(AccountController.Destroy.url(user.id), {
    preserveScroll: true,
    onSuccess: () => {
      toast.add({
        title: 'User deleted',
        description: `${user.name} was removed.`,
        icon: 'i-lucide-check',
        color: 'success',
      });
    },
    onError: () => {
      toast.add({
        title: 'Delete failed',
        description: `Could not delete ${user.name}.`,
        icon: 'i-lucide-x',
        color: 'error',
      });
    },
  });
}

const UButton = resolveComponent('UButton');
const UDropdownMenu = resolveComponent('UDropdownMenu');
const UBadge = resolveComponent('UBadge');

const columns: TableColumn<Account>[] = [
  {
    accessorKey: 'id',
    header: 'ID',
    cell: ({ row }) => h('span', { class: 'font-medium text-muted font-mono' }, [row.original.id]),
  },
  {
    accessorKey: 'name',
    header: 'Name',
    cell: ({ row }) => h('p', { class: 'font-medium text-highlighted' }, [row.original.name]),
  },
  {
    accessorKey: 'username',
    header: 'Username',
  },
  {
    id: 'role',
    header: 'Role',
    cell: ({ row }) => (row.original.role ? h(UBadge, { variant: 'subtle', color: 'primary' }, () => row.original.role!.name) : h('span', '—')),
  },
  {
    id: 'actions',
    cell: ({ row }) =>
      h('div', { class: 'text-right' }, [
        h(
          UDropdownMenu,
          {
            content: { align: 'end' },
            items: [
              { type: 'label', label: 'Actions' },
              { type: 'separator' },
              {
                label: 'Delete user',
                icon: 'i-lucide-trash',
                color: 'error',
                onSelect: () => deleteUser(row.original),
              },
            ],
          },
          [
            h(UButton, {
              icon: 'i-lucide-ellipsis-vertical',
              color: 'neutral',
              variant: 'ghost',
              class: 'ml-auto',
            }),
          ],
        ),
      ]),
  },
];
</script>

<template>
  <Head>
    <title>Users</title>
  </Head>

  <DashboardLayout>
    <UDashboardPanel id="users">
      <template #header>
        <UDashboardNavbar title="Users" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
          <template #right>
            <UButton label="New user" icon="i-lucide-plus" :to="AccountController.Create.url()" />
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <div class="flex flex-wrap items-center justify-between gap-1.5">
          <UInput v-model="searchQuery" class="max-w-sm" icon="i-lucide-search" placeholder="Search users..." />
        </div>

        <UTable :data="props.users.data ?? []" :columns="columns" class="shrink-0" />

        <div class="mt-auto flex items-center justify-between gap-3 border-t border-default pt-4">
          <div class="text-sm text-muted">{{ props.users?.total ?? 0 }} user(s) total.</div>
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
