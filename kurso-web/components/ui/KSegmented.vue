<script setup lang="ts">
// Segmented control (2–3 equal pills in a well). `#icon` slot renders a leading
// glyph per option (e.g. up/down chevrons on the condition toggle).
export interface SegOption {
  value: string
  label: string
}

defineProps<{ options: SegOption[]; size?: 'sm' | 'md' }>()
const model = defineModel<string>({ required: true })
</script>

<template>
  <div class="flex gap-1.5 rounded-xl border border-line bg-well p-1.5">
    <button
      v-for="o in options"
      :key="o.value"
      type="button"
      class="flex flex-1 items-center justify-center gap-1.5 rounded-lg font-semibold transition-colors"
      :class="[
        size === 'sm' ? 'py-2 text-[13px]' : 'py-2.5 text-sm',
        model === o.value
          ? 'bg-surface-chip text-ink shadow-[0_2px_8px_rgba(0,0,0,0.3)]'
          : 'text-ink-dim hover:text-ink',
      ]"
      @click="model = o.value"
    >
      <slot name="icon" :option="o" :active="model === o.value" />
      {{ o.label }}
    </button>
  </div>
</template>
