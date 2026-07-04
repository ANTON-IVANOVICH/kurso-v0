<script setup lang="ts">
// On/off switch matching the design's pill toggle (blue when on, knob slides).
const model = defineModel<boolean>({ default: false })
withDefaults(defineProps<{ size?: 'sm' | 'md' }>(), { size: 'md' })
</script>

<template>
  <button
    type="button"
    role="switch"
    :aria-checked="model"
    class="relative flex-none rounded-full transition-colors"
    :class="[
      model ? 'bg-brand' : 'bg-line-strong',
      size === 'md' ? 'h-[26px] w-[46px]' : 'h-6 w-[42px]',
    ]"
    @click="model = !model"
  >
    <span
      class="absolute top-[3px] rounded-full transition-all"
      :class="[
        size === 'md' ? 'h-5 w-5' : 'h-[18px] w-[18px]',
        model
          ? (size === 'md' ? 'left-[23px]' : 'left-[21px]') + ' bg-white'
          : 'left-[3px] bg-ink-muted',
      ]"
    />
  </button>
</template>
