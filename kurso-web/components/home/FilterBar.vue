<script setup lang="ts">
import type { HomeFilter, HomeSort } from '~/composables/useHomeFilters'

// Working filters for the home list — wired to shared state the page reads.
// Only data-backed filters are shown (partner flag + reserve threshold).
const { filter, sort } = useHomeFilters()
const { t } = useI18n()

const chips: { key: HomeFilter; i18n: string }[] = [
  { key: 'all', i18n: 'home.filters.all' },
  { key: 'partners', i18n: 'home.filters.partners' },
  { key: 'reserve5m', i18n: 'home.filters.reserve5m' },
]
const sorts: { key: HomeSort; i18n: string }[] = [
  { key: 'best', i18n: 'home.filters.byBest' },
  { key: 'rating', i18n: 'home.filters.byRating' },
]
const sortLabel = computed(() => t(sorts.find((s) => s.key === sort.value)?.i18n ?? ''))
function cycleSort() {
  const i = sorts.findIndex((s) => s.key === sort.value)
  sort.value = sorts[(i + 1) % sorts.length].key
}

// KChip is pill-shaped on mobile; these !important md overrides switch it back to
// the squarer desktop chip so the whole bar lives in one responsive DOM.
const desk = 'md:!rounded-md md:!px-[14px] md:!py-2 md:!text-sm'
</script>

<template>
  <!-- Full-bleed on mobile: chips scroll to the screen edge, not into the page
       gutter. -mx-4/px-4 cancels the page padding, then re-insets the first chip;
       reset at md where the bar lives inside the results column. -->
  <div
    class="scrollx -mx-4 mb-4 flex items-center gap-2.5 overflow-x-auto border-b border-line-subtle px-4 pb-3.5 md:mx-0 md:px-0"
  >
    <KChip
      v-for="c in chips"
      :key="c.key"
      :active="filter === c.key"
      pill
      :class="desk"
      @click="filter = c.key"
    >
      {{ t(c.i18n) }}
    </KChip>
    <KChip pill :class="[desk, 'ml-auto !bg-surface !text-ink-muted']" @click="cycleSort">
      <svg
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M3 6h18M6 12h12M10 18h4" />
      </svg>
      {{ sortLabel }}
    </KChip>
  </div>
</template>
