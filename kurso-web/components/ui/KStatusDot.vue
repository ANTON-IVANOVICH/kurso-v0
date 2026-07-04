<script setup lang="ts">
type Tone = 'success' | 'warning' | 'danger' | 'neutral'

withDefaults(defineProps<{ tone?: Tone; pulse?: boolean; size?: number }>(), {
  tone: 'success',
  pulse: false,
  size: 8,
})

const colors: Record<Tone, string> = {
  success: '#2BC58C',
  warning: '#E0A954',
  danger: '#EC5B5B',
  neutral: '#4A5159',
}
</script>

<template>
  <span
    class="inline-block flex-none rounded-full"
    :class="pulse ? 'animate-kpulse' : ''"
    :style="{ width: `${size}px`, height: `${size}px`, background: colors[tone] }"
  />
</template>
