<script setup lang="ts">
import { computed } from 'vue'

// Sparkline behind the alert builder: real best-rate history with the threshold
// drawn as a dashed line. Everything is mapped into a 320×130 viewBox so the
// SVG scales fluidly. Degrades to a flat line when history is sparse.
const props = withDefaults(
  defineProps<{
    series: number[] // chronological best-rate points (already unit-adjusted)
    threshold: number
    current: number | null
    height?: number
  }>(),
  { height: 130 },
)

const W = 320
const H = 130
const PAD = 12

// Points to plot — fall back to a flat line at the current rate if history is thin.
const pts = computed<number[]>(() => {
  if (props.series.length >= 2) return props.series
  const c = props.current ?? props.series[0] ?? props.threshold
  return [c, c]
})

const bounds = computed(() => {
  const vals = [...pts.value, props.threshold]
  if (props.current != null) vals.push(props.current)
  let min = Math.min(...vals)
  let max = Math.max(...vals)
  if (max - min < 1e-9) {
    min -= Math.abs(min) * 0.01 + 1
    max += Math.abs(max) * 0.01 + 1
  }
  const span = max - min
  return { min: min - span * 0.12, max: max + span * 0.12 }
})

const x = (i: number, n: number) => PAD + (i / Math.max(1, n - 1)) * (W - 2 * PAD)
const y = (v: number) => {
  const { min, max } = bounds.value
  return H - PAD - ((v - min) / (max - min)) * (H - 2 * PAD)
}

const linePath = computed(() => {
  const n = pts.value.length
  return pts.value
    .map((v, i) => `${i === 0 ? 'M' : 'L'}${x(i, n).toFixed(1)},${y(v).toFixed(1)}`)
    .join(' ')
})
const areaPath = computed(() => {
  const n = pts.value.length
  const line = pts.value.map((v, i) => `L${x(i, n).toFixed(1)},${y(v).toFixed(1)}`).join(' ')
  return `M${x(0, n).toFixed(1)},${(H - PAD).toFixed(1)} ${line.slice(1)} L${x(n - 1, n).toFixed(1)},${(H - PAD).toFixed(1)} Z`
})
const thresholdY = computed(() => y(props.threshold))
const lastX = computed(() => x(pts.value.length - 1, pts.value.length))
const lastY = computed(() => y(pts.value[pts.value.length - 1]))
const gid = `alertfill-${Math.round((props.threshold || 1) * 1000) % 100000}`
</script>

<template>
  <svg
    :viewBox="`0 0 ${W} ${H}`"
    class="w-full"
    :style="{ height: `${height}px` }"
    preserveAspectRatio="none"
  >
    <defs>
      <linearGradient :id="gid" x1="0" y1="0" x2="0" y2="1">
        <stop offset="0%" stop-color="rgba(46,125,242,0.28)" />
        <stop offset="100%" stop-color="rgba(46,125,242,0)" />
      </linearGradient>
    </defs>
    <line
      :x1="0"
      :y1="thresholdY"
      :x2="W"
      :y2="thresholdY"
      stroke="#6BA6FF"
      stroke-width="1.4"
      stroke-dasharray="5 5"
      vector-effect="non-scaling-stroke"
    />
    <path :d="areaPath" :fill="`url(#${gid})`" />
    <path
      :d="linePath"
      fill="none"
      stroke="#2E7DF2"
      stroke-width="2.5"
      stroke-linecap="round"
      stroke-linejoin="round"
      vector-effect="non-scaling-stroke"
    />
    <circle :cx="lastX" :cy="lastY" r="4.5" fill="#2E7DF2" stroke="#0E1013" stroke-width="3" />
  </svg>
</template>
