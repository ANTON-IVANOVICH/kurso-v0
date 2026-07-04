<script setup lang="ts">
import { computed } from 'vue'

type Tone = 'brand' | 'success' | 'warning' | 'danger' | 'violet' | 'neutral'

const props = withDefaults(defineProps<{ tone?: Tone }>(), { tone: 'neutral' })

const tones: Record<Tone, string> = {
  brand: 'text-brand-bright bg-brand/15',
  success: 'text-success-bright bg-success/10 border border-success/30',
  warning: 'text-warning bg-warning/[0.12] border border-warning/30',
  danger: 'text-danger bg-danger/[0.12] border border-danger/30',
  violet: 'text-violet-light bg-violet/15',
  neutral: 'text-ink-muted bg-surface-raised border border-line-strong',
}

const classes = computed(() => [
  'inline-flex items-center gap-1 rounded-sm px-2.5 py-1 text-[11px] font-semibold',
  tones[props.tone],
])
</script>

<template>
  <span :class="classes"><slot /></span>
</template>
