<script setup lang="ts">
const { t } = useI18n()

definePageMeta({ layout: 'account', middleware: 'auth' })
useSeoMeta({ title: () => t('historyPage.seoTitle') })

const { history } = useHistory()
</script>

<template>
  <div class="max-w-3xl">
    <h1 class="text-2xl font-extrabold tracking-[-0.02em] text-ink">
      {{ t('historyPage.title') }}
    </h1>
    <p class="mb-5 mt-1 text-sm text-ink-faint">{{ t('historyPage.subtitle') }}</p>

    <div v-if="history.length" class="overflow-hidden rounded-2xl border border-line bg-surface">
      <AccountHistoryRow v-for="h in history" :key="h.id" :entry="h" />
    </div>

    <div
      v-else
      class="flex flex-col items-center gap-3 rounded-2xl border border-dashed border-line-strong bg-surface/50 px-6 py-12 text-center"
    >
      <span
        class="flex h-12 w-12 items-center justify-center rounded-full bg-success/[0.12] text-success-bright"
      >
        <svg
          width="22"
          height="22"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <circle cx="12" cy="12" r="9" />
          <path d="M12 8v4l3 2" />
        </svg>
      </span>
      <p class="max-w-xs text-sm text-ink-muted">
        {{ t('historyPage.empty') }}
      </p>
      <NuxtLink to="/" class="text-sm font-semibold text-brand-bright">{{
        t('historyPage.toRates')
      }}</NuxtLink>
    </div>
  </div>
</template>
