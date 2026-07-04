<script setup lang="ts">
// Kurso brand mark: gradient square with the up/down arrows glyph + wordmark.
withDefaults(defineProps<{ size?: number; wordmark?: boolean }>(), {
  size: 40,
  wordmark: true,
})
</script>

<template>
  <span class="inline-flex items-center gap-2.5">
    <span
      class="flex items-center justify-center rounded-xl bg-[linear-gradient(150deg,#4A90F5,#2E7DF2)] shadow-glow"
      :style="{ width: `${size}px`, height: `${size}px` }"
    >
      <!-- Size via :style, not :width/:height — SVG width/height are read-only
           DOM props, so binding them throws during hydration patching. -->
      <svg
        viewBox="0 0 24 24"
        fill="none"
        stroke="#fff"
        stroke-width="2.1"
        stroke-linecap="round"
        stroke-linejoin="round"
        :style="{ width: `${Math.round(size * 0.48)}px`, height: `${Math.round(size * 0.48)}px` }"
      >
        <path d="M7 4 L7 20 M7 20 L4 17 M7 20 L10 17 M17 20 L17 4 M17 4 L14 7 M17 4 L20 7" />
      </svg>
    </span>
    <span
      v-if="wordmark"
      class="font-extrabold tracking-[-0.02em]"
      :style="{ fontSize: `${Math.round(size * 0.52)}px` }"
      >Kurso</span
    >
  </span>
</template>
