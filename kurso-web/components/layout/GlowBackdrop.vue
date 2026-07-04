<script setup lang="ts">
// Ambient radial glow anchored top-centre, matching the design language.
type Tone = 'brand' | 'danger' | 'warning'
withDefaults(defineProps<{ tone?: Tone }>(), { tone: 'brand' })

const toneClass: Record<Tone, string> = {
  brand: 'bg-[radial-gradient(ellipse_at_center,rgba(46,125,242,0.16),rgba(46,125,242,0)_62%)]',
  danger: 'bg-[radial-gradient(ellipse_at_center,rgba(236,91,74,0.12),rgba(236,91,74,0)_64%)]',
  warning: 'bg-[radial-gradient(ellipse_at_center,rgba(217,154,51,0.12),rgba(217,154,51,0)_64%)]',
}
</script>

<template>
  <div
    aria-hidden="true"
    class="pointer-events-none absolute left-1/2 top-[-320px] z-0 h-[720px] w-[1080px] max-w-none -translate-x-1/2"
    :class="toneClass[tone]"
  />
</template>
