<script setup lang="ts">
defineProps<{ title: string; subtitle?: string }>()
</script>

<template>
  <div class="mb-5 flex flex-wrap items-end justify-between gap-3">
    <div>
      <h1 class="text-xl font-extrabold tracking-tight md:text-[22px]">{{ title }}</h1>
      <p v-if="subtitle" class="mt-1 text-sm text-ink-faint">{{ subtitle }}</p>
    </div>
    <div v-if="$slots.actions" class="flex items-center gap-2">
      <slot name="actions" />
    </div>
  </div>
</template>
