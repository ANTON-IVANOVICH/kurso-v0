<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    tone?: 'success' | 'warn' | 'danger' | 'muted' | 'brand'
    pulse?: boolean
    size?: number
  }>(),
  { tone: 'muted', pulse: false, size: 7 },
)

const COLOR: Record<string, string> = {
  success: '#2BC58C',
  warn: '#D99A33',
  danger: '#EC5B5B',
  muted: '#5A616A',
  brand: '#2E7DF2',
}
const color = computed(() => COLOR[props.tone])
</script>

<template>
  <span
    class="inline-block flex-none rounded-full"
    :class="{ 'animate-kpulse': pulse }"
    :style="{ width: `${size}px`, height: `${size}px`, background: color }"
  />
</template>
