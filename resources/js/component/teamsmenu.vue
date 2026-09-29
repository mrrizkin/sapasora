<script setup lang="ts">
defineProps<{
  collapsed?: boolean;
}>();
</script>

<template>
  <div class="flex items-center gap-2 px-2 py-2">
    <UIcon name="i-lucide-message-square" class="size-5 shrink-0 text-primary" />
    <span v-if="!collapsed" class="truncate font-semibold text-highlighted">Sapasora</span>
  </div>
</template>
