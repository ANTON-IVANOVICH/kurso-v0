<script setup lang="ts">
// Central icon set. Every glyph is a 24×24 stroked path so weight/size stay
// consistent with the mockups. Inner markup is static and trusted (v-html).
withDefaults(
  defineProps<{ name: keyof typeof ICONS; size?: number | string; stroke?: number | string }>(),
  {
    size: 16,
    stroke: 2,
  },
)

const ICONS = {
  logo: '<path d="M7 4 L7 20 M7 20 L4 17 M7 20 L10 17 M17 20 L17 4 M17 4 L14 7 M17 4 L20 7"/>',
  grid: '<rect x="3" y="3" width="7" height="9" rx="1.5"/><rect x="14" y="3" width="7" height="5" rx="1.5"/><rect x="14" y="12" width="7" height="9" rx="1.5"/><rect x="3" y="16" width="7" height="5" rx="1.5"/>',
  search: '<circle cx="11" cy="11" r="7"/><path d="M21 21l-3.5-3.5"/>',
  lock: '<rect x="5" y="11" width="14" height="9" rx="2"/><path d="M8 11V8a4 4 0 0 1 8 0v3"/>',
  chevron: '<path d="M6 9l6 6 6-6"/>',
  alert:
    '<path d="M12 9v4M12 17h.01"/><path d="M10.3 3.9 1.8 18a2 2 0 0 0 1.7 3h17a2 2 0 0 0 1.7-3L13.7 3.9a2 2 0 0 0-3.4 0Z"/>',
  play: '<path d="M5 3l14 9-14 9V3Z"/>',
  check: '<path d="M5 12l5 5L20 6"/>',
  arrowRight: '<path d="M5 12h14M13 6l6 6-6 6"/>',
  download: '<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3"/>',
  bars: '<path d="M4 19V5M9 19v-7M14 19V9M19 19v-4"/>',
  crossAxes: '<path d="M12 2v20M2 12h20"/>',
  chat: '<path d="M21 15a2 2 0 0 1-2 2H8l-4 4V6a2 2 0 0 1 2-2h13a2 2 0 0 1 2 2Z"/>',
  logout: '<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9"/>',
  refresh: '<path d="M21 12a9 9 0 1 1-3-6.7M21 4v4h-4"/>',
  external:
    '<path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6M15 3h6v6M10 14 21 3"/>',
  plus: '<path d="M12 5v14M5 12h14"/>',
  clock: '<circle cx="12" cy="12" r="9"/><path d="M12 8v4l3 2"/>',
  shield: '<path d="M12 3l7 3v6c0 4-3 7-7 8-4-1-7-4-7-8V6l7-3Z"/>',
  eye: '<path d="M2 12s3.5-7 10-7 10 7 10 7-3.5 7-10 7-10-7-10-7Z"/><circle cx="12" cy="12" r="3"/>',
  layers: '<path d="M12 2 2 7l10 5 10-5-10-5ZM2 17l10 5 10-5M2 12l10 5 10-5"/>',
  filter: '<path d="M3 5h18M6 12h12M10 19h4"/>',
  user: '<circle cx="12" cy="8" r="4"/><path d="M4 21c0-4 3.5-6 8-6s8 2 8 6"/>',
  flag: '<path d="M5 21V5a1 1 0 0 1 1-1h9l1 2h3v8h-8l-1-2H5"/>',
  card: '<rect x="3" y="5" width="18" height="14" rx="2"/><path d="M3 10h18"/>',
} as const
</script>

<template>
  <!-- eslint-disable-next-line vue/no-v-html -->
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    :stroke-width="stroke"
    stroke-linecap="round"
    stroke-linejoin="round"
    v-html="ICONS[name]"
  />
</template>
