<script setup lang="ts">
import { Head, router } from '@inertiajs/vue3';
import type { TableColumn } from '@nuxt/ui';
import { useToast } from '@nuxt/ui/runtime/composables/useToast.js';
import type { PaginationState, Row, RowSelectionState, VisibilityState } from '@tanstack/table-core';
import { titleCase } from 'scule';
import { type GlobalComponents, computed, h, ref, resolveComponent, useTemplateRef, watch } from 'vue';
import type { ComponentExposed } from 'vue-component-type-helpers';

import DevicesDeleteModal from '@/js/component/devices/devices_delete_modal.vue';
import DashboardLayout from '@/js/layout/dashboard.vue';
import DeviceController from '@/js/lib/actions/app/http/controllers/device/DeviceController';
import { debounce } from '@/js/lib/utils';
import type { Device, Paginated } from '@/js/types';

import type { ColumnFilter } from './type';

interface Props {
  devices: Paginated<Device>;
  filters: ColumnFilter;
  pagination: PaginationState;
}

const props = defineProps<Props>();

// Composables
const toast = useToast();
const table = useTemplateRef<ComponentExposed<GlobalComponents['UTable']>>('table');

// State Management
const state = ref<{
  columnFilters: ColumnFilter;
  columnVisibility: VisibilityState;
  rowSelection: RowSelectionState;
  pagination: PaginationState;
}>({
  columnFilters: props.filters,
  columnVisibility: {},
  rowSelection: {},
  pagination: props.pagination,
});

// watch state update
watch(
  () => [state.value.columnFilters, state.value.pagination] as const,
  () => {
    router.visit(
      DeviceController.Index({
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
  { deep: true }, // Skip immediate run on mount
);

// Filters
const updateFilters = debounce((filters: ColumnFilter) => {
  if (state.value.columnFilters.search !== filters.search) {
    state.value.columnFilters.search = filters.search;
    state.value.pagination.pageIndex = 0; // reset pagination on search
  }
}, 500);

// Search
const searchQuery = computed({
  get: () => state.value.columnFilters.search || '',
  set: (value: string) => {
    updateFilters({ ...state.value.columnFilters, search: value });
  },
});

// Selection
const selected = computed(() => (table.value?.tableApi?.getSelectedRowModel()?.rows ?? []) as Row<Device>[]);

const selectedCount = computed(() => selected.value.length);

// Column Visibility
const columnVisibilityItems = computed(
  () =>
    table.value?.tableApi
      ?.getAllColumns()
      ?.filter((col: any) => col.getCanHide())
      ?.map((col: any) => ({
        label: titleCase(col.id),
        type: 'checkbox' as const,
        checked: col.getIsVisible(),
        onUpdateChecked: (checked: boolean) => {
          table.value?.tableApi.getColumn(col.id)?.toggleVisibility(!!checked);
        },
        onSelect: (e?: Event) => e?.preventDefault(),
      })) ?? [],
);

// Actions
function copyDeviceId(deviceId: string) {
  navigator.clipboard.writeText(deviceId.toString());
  toast.add({
    title: 'Copied to clipboard',
    description: 'Device ID copied to clipboard',
  });
}

function deleteDevice() {
  toast.add({
    title: 'Device deleted',
    description: 'The device has been deleted.',
  });
}

function getRowActions(row: Row<Device>) {
  return [
    { type: 'label', label: 'Actions' },
    {
      label: 'Copy device ID',
      icon: 'i-lucide-copy',
      onSelect: () => copyDeviceId(row.original.id),
    },
    { type: 'separator' },
    {
      label: 'View device details',
      icon: 'i-lucide-list',
      onSelect: () => router.visit(DeviceController.Show(row.original.id)),
    },
    { type: 'separator' },
    {
      label: 'Delete device',
      icon: 'i-lucide-trash',
      color: 'error',
      onSelect: deleteDevice,
    },
  ];
}

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

const getTypeIcon = (type: string) => {
  const icons: Record<string, string> = {
    whatsapp: 'i-lucide-message-circle',
    telegram: 'i-lucide-send',
    sms: 'i-lucide-smartphone',
    email: 'i-lucide-mail',
  };
  return icons[type.toLowerCase()] || 'i-lucide-radio';
};

// Column Renderers
const UButton = resolveComponent('UButton');
const UDropdownMenu = resolveComponent('UDropdownMenu');
const UCheckbox = resolveComponent('UCheckbox');
const UBadge = resolveComponent('UBadge');
const UIcon = resolveComponent('UIcon');

// Table Configuration
const columns: TableColumn<Device>[] = [
  {
    id: 'select',
    header: ({ table }) =>
      h(UCheckbox, {
        modelValue: table.getIsSomePageRowsSelected() ? 'indeterminate' : table.getIsAllPageRowsSelected(),
        'onUpdate:modelValue': (value: boolean | 'indeterminate') => table.toggleAllPageRowsSelected(!!value),
        ariaLabel: 'Select all',
      }),
    cell: ({ row }) =>
      h(UCheckbox, {
        modelValue: row.getIsSelected(),
        'onUpdate:modelValue': (value: boolean | 'indeterminate') => row.toggleSelected(!!value),
        ariaLabel: 'Select row',
      }),
  },
  {
    accessorKey: 'id',
    header: 'ID',
    cell: ({ row }) => h('span', { class: 'font-medium text-muted font-mono' }, [row.original.id]),
  },
  {
    accessorKey: 'name',
    header: ({ column }) => {
      const isSorted = column.getIsSorted();
      const icon = isSorted ? (isSorted === 'asc' ? 'i-lucide-arrow-up-narrow-wide' : 'i-lucide-arrow-down-wide-narrow') : 'i-lucide-arrow-up-down';

      return h(UButton, {
        color: 'neutral',
        variant: 'ghost',
        label: 'Name',
        icon,
        class: '-mx-2.5',
        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
      });
    },
    cell: ({ row }) => h('div', { class: 'flex items-center gap-3' }, [h('p', { class: 'font-medium text-highlighted' }, [row.original.name])]),
  },
  {
    accessorKey: 'type',
    header: 'Type',
    cell: ({ row }: any) =>
      h('div', { class: 'flex items-center gap-3' }, [
        h('div', { class: 'flex items-center gap-3' }, [
          h(UIcon, { name: getTypeIcon(row.original.type), class: 'h-5 w-5 text-gray-500 dark:text-gray-400' }),
          h('span', { class: 'text-sm text-gray-500 dark:text-gray-400' }, row.original.type),
        ]),
      ]),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    filterFn: 'equals',
    cell: ({ row }) =>
      h('div', { class: 'flex items-center gap-3' }, [
        h(UBadge, { color: getStatusColor(row.original.status), variant: 'subtle', size: 'lg', label: row.original.status }),
      ]),
  },
  {
    id: 'actions',
    cell: ({ row }) =>
      h('div', { class: 'text-right' }, [
        h(UDropdownMenu, { content: { align: 'end' }, items: getRowActions(row) }, [
          h(UButton, {
            icon: 'i-lucide-ellipsis-vertical',
            color: 'neutral',
            variant: 'ghost',
            class: 'ml-auto',
          }),
        ]),
      ]),
  },
];
</script>

<template>
  <Head>
    <title>Devices</title>
  </Head>

  <DashboardLayout>
    <UDashboardPanel id="devices">
      <template #header>
        <UDashboardNavbar title="Devices" :ui="{ right: 'gap-3' }">
          <template #leading>
            <UDashboardSidebarCollapse />
          </template>
          <template #right>
            <UButton label="New device" icon="i-lucide-plus" :to="DeviceController.Create.url()" />
          </template>
        </UDashboardNavbar>
      </template>

      <template #body>
        <!-- Toolbar -->
        <div class="flex flex-wrap items-center justify-between gap-1.5">
          <UInput v-model="searchQuery" class="max-w-sm" icon="i-lucide-search" placeholder="Search devices..." />

          <div class="flex flex-wrap items-center gap-1.5">
            <DevicesDeleteModal :count="selectedCount" :selected="selected">
              <UButton v-if="selectedCount" label="Delete" color="error" variant="subtle" icon="i-lucide-trash">
                <template #trailing>
                  <UKbd>{{ selectedCount }}</UKbd>
                </template>
              </UButton>
            </DevicesDeleteModal>

            <UDropdownMenu :items="columnVisibilityItems" :content="{ align: 'end' }">
              <UButton label="Display" color="neutral" variant="outline" trailing-icon="i-lucide-settings-2" />
            </UDropdownMenu>
          </div>
        </div>

        <!-- Table -->
        <UTable
          ref="table"
          v-model:column-visibility="state.columnVisibility"
          v-model:row-selection="state.rowSelection"
          v-model:pagination="state.pagination"
          class="shrink-0"
          :data="props.devices.data ?? []"
          :columns="columns"
          :pagination-options="{
            manualPagination: true,
            rowCount: props.devices.total ?? 0,
          }"
          :ui="{
            base: 'table-fixed border-separate border-spacing-0',
            thead: '[&>tr]:bg-elevated/50 [&>tr]:after:content-none',
            tbody: '[&>tr]:last:[&>td]:border-b-0',
            th: 'py-2 first:rounded-l-lg last:rounded-r-lg border-y border-default first:border-l last:border-r',
            td: 'border-b border-default',
            separator: 'h-0',
          }" />

        <!-- Footer -->
        <div class="mt-auto flex items-center justify-between gap-3 border-t border-default pt-4">
          <div class="text-sm text-muted">{{ selectedCount }} of {{ props.devices?.total ?? 0 }} row(s) selected.</div>

          <UPagination
            :default-page="(table?.tableApi.getState().pagination.pageIndex || 0) + 1"
            :items-per-page="table?.tableApi.getState().pagination.pageSize"
            :total="props.devices?.total ?? 0"
            @update:page="(p: number) => table?.tableApi.setPageIndex(p - 1)" />
        </div>
      </template>
    </UDashboardPanel>
  </DashboardLayout>
</template>
