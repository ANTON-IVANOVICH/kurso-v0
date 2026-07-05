<script setup lang="ts">
import { ref } from 'vue'

const { t } = useI18n()

definePageMeta({ layout: 'legal', legalTitle: 'Файлы cookies' })

useSeoMeta({
  title: () => t('cookies.seoTitle'),
  description: () => t('cookies.seoDescription'),
})

const analytics = ref(true)
const marketing = ref(false)

function acceptAll() {
  analytics.value = true
  marketing.value = true
}
function onlyNecessary() {
  analytics.value = false
  marketing.value = false
}
</script>

<template>
  <div class="mx-auto max-w-[820px] px-4 pt-8 md:px-6 md:pt-12">
    <div class="font-label text-[11px] uppercase tracking-[0.14em] text-ink-faint">
      {{ t('cookies.eyebrow') }}
    </div>
    <h1
      class="mt-3 text-[26px] font-extrabold tracking-[-0.025em] text-ink md:text-[34px] md:font-black"
    >
      {{ t('cookies.title') }}
    </h1>
    <p class="mt-4 max-w-[660px] leading-relaxed text-ink-muted">
      {{ t('cookies.intro') }}
    </p>

    <!-- categories -->
    <div class="mt-7 overflow-hidden rounded-2xl border border-line bg-surface">
      <!-- necessary (always on) -->
      <div class="flex items-center justify-between gap-4 border-b border-line p-5">
        <div class="max-w-[520px]">
          <div class="font-bold text-ink">{{ t('cookies.necessaryTitle') }}</div>
          <div class="mt-1 text-[13px] leading-relaxed text-ink-muted">
            {{ t('cookies.necessaryDesc') }}
          </div>
        </div>
        <div
          class="relative h-[26px] w-[46px] flex-none rounded-full bg-brand opacity-65"
          aria-hidden="true"
        >
          <span class="absolute right-[3px] top-[3px] h-5 w-5 rounded-full bg-white" />
        </div>
      </div>

      <!-- analytics -->
      <div class="flex items-center justify-between gap-4 border-b border-line p-5">
        <div class="max-w-[520px]">
          <div class="font-bold text-ink">{{ t('cookies.analyticsTitle') }}</div>
          <div class="mt-1 text-[13px] leading-relaxed text-ink-muted">
            {{ t('cookies.analyticsDesc') }}
          </div>
        </div>
        <button
          type="button"
          role="switch"
          :aria-checked="analytics"
          class="relative h-[26px] w-[46px] flex-none rounded-full transition-colors"
          :class="analytics ? 'bg-brand' : 'border border-line-strong bg-surface-chip'"
          @click="analytics = !analytics"
        >
          <span
            class="absolute top-[3px] h-5 w-5 rounded-full transition-all"
            :class="analytics ? 'right-[3px] bg-white' : 'left-[3px] bg-ink-faint'"
          />
        </button>
      </div>

      <!-- marketing -->
      <div class="flex items-center justify-between gap-4 p-5">
        <div class="max-w-[520px]">
          <div class="font-bold text-ink">{{ t('cookies.marketingTitle') }}</div>
          <div class="mt-1 text-[13px] leading-relaxed text-ink-muted">
            {{ t('cookies.marketingDesc') }}
          </div>
        </div>
        <button
          type="button"
          role="switch"
          :aria-checked="marketing"
          class="relative h-[26px] w-[46px] flex-none rounded-full transition-colors"
          :class="marketing ? 'bg-brand' : 'border border-line-strong bg-surface-chip'"
          @click="marketing = !marketing"
        >
          <span
            class="absolute top-[3px] h-5 w-5 rounded-full transition-all"
            :class="marketing ? 'right-[3px] bg-white' : 'left-[3px] bg-ink-faint'"
          />
        </button>
      </div>
    </div>

    <div class="mt-5 flex flex-col gap-3 sm:flex-row sm:flex-wrap">
      <KButton>{{ t('cookies.save') }}</KButton>
      <KButton variant="secondary" @click="acceptAll">{{ t('cookies.acceptAll') }}</KButton>
      <KButton variant="ghost" class="border border-line-strong" @click="onlyNecessary">
        {{ t('cookies.onlyNecessary') }}
      </KButton>
    </div>

    <div class="mt-9 max-w-[660px] space-y-6">
      <div>
        <h2 class="text-[17px] font-bold text-ink">{{ t('cookies.whatTitle') }}</h2>
        <p class="mt-2 leading-relaxed text-ink-muted">
          {{ t('cookies.whatDesc') }}
        </p>
      </div>
      <div>
        <h2 class="text-[17px] font-bold text-ink">{{ t('cookies.manageTitle') }}</h2>
        <p class="mt-2 leading-relaxed text-ink-muted">
          {{ t('cookies.manageDesc') }}
        </p>
      </div>
    </div>
  </div>
</template>
