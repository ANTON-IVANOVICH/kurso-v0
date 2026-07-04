<script setup lang="ts">
withDefaults(defineProps<{ active?: boolean; pill?: boolean }>(), {
  active: false,
  pill: false,
})
</script>

<template>
  <button
    :class="[
      'inline-flex items-center gap-1.5 whitespace-nowrap font-medium transition-colors',
      pill ? 'rounded-full px-[13px] py-[7px] text-[13px]' : 'rounded-md px-[14px] py-2 text-sm',
      active
        ? 'bg-brand text-white'
        : 'border border-line-strong bg-surface-raised text-ink hover:border-[#3A4047]',
    ]"
  >
    <slot />
  </button>
</template>
