<script setup lang="ts">
// Mobile-only top bar for nested content pages: a back chevron + page title.
// Replaces the site header/bottom-nav so legal pages read as nested screens.
defineProps<{ title: string }>()

const { t } = useI18n()
const router = useRouter()

function goBack() {
  if (import.meta.client && window.history.length > 1) router.back()
  else router.push('/')
}
</script>

<template>
  <div class="relative z-20 flex items-center gap-3 border-b border-line-subtle px-4 py-3">
    <button
      type="button"
      :aria-label="t('legalBackBar.back')"
      class="flex h-9 w-9 flex-none items-center justify-center rounded-full border border-line bg-surface text-ink transition-colors hover:border-line-strong"
      @click="goBack"
    >
      <svg
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M15 6l-6 6 6 6" />
      </svg>
    </button>
    <div class="flex-1 truncate text-base font-bold text-ink">{{ title }}</div>
  </div>
</template>
