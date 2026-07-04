<script setup lang="ts">
import { computed, ref } from 'vue'

withDefaults(defineProps<{ compact?: boolean }>(), { compact: false })

const locale = useLocale()
const { t } = useI18n()
const open = ref(false)
const current = computed(() => locales.find((o) => o.code === locale.value) ?? locales[0])

function select(code: Locale) {
  locale.value = code
  open.value = false
}
</script>

<template>
  <div class="relative">
    <button
      v-if="compact"
      type="button"
      :aria-label="t('header.langAria')"
      class="flex h-[38px] w-[38px] items-center justify-center rounded-full border border-line bg-surface text-ink-muted"
      @click="open = !open"
    >
      <svg
        width="18"
        height="18"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="9" />
        <path d="M3 12h18" />
        <path d="M12 3a15 15 0 0 1 0 18a15 15 0 0 1 0-18" />
      </svg>
    </button>
    <button
      v-else
      type="button"
      class="inline-flex h-[42px] items-center gap-2 rounded-full border border-line-strong bg-surface px-3.5 text-ink transition-colors hover:border-[#3A4047]"
      @click="open = !open"
    >
      <svg
        class="text-ink-muted"
        width="17"
        height="17"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <circle cx="12" cy="12" r="9" />
        <path d="M3 12h18" />
        <path d="M12 3a15 15 0 0 1 0 18a15 15 0 0 1 0-18" />
      </svg>
      <span class="w-[22px] text-sm font-semibold">{{ current.label }}</span>
      <svg
        class="text-ink-faint transition-transform"
        :class="open ? 'rotate-180' : ''"
        width="13"
        height="13"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2.4"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <path d="M6 9l6 6 6-6" />
      </svg>
    </button>

    <template v-if="open">
      <div class="fixed inset-0 z-30" @click="open = false" />
      <div
        class="absolute right-0 z-40 mt-2 w-44 rounded-xl border border-line-strong bg-surface-hi p-1.5 shadow-pop"
      >
        <button
          v-for="o in locales"
          :key="o.code"
          type="button"
          class="flex w-full items-center gap-3 rounded-lg px-2.5 py-2 text-left transition-colors hover:bg-white/[0.04]"
          :class="o.code === locale ? 'bg-brand/[0.12]' : ''"
          @click="select(o.code)"
        >
          <span class="w-7 font-label text-xs font-semibold text-ink-muted">{{ o.label }}</span>
          <span class="flex-1 text-sm">{{ o.name }}</span>
          <svg
            v-if="o.code === locale"
            class="text-brand-bright"
            width="15"
            height="15"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2.6"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <path d="M5 12l5 5L20 6" />
          </svg>
        </button>
      </div>
    </template>
  </div>
</template>
