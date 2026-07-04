<script setup lang="ts">
import { computed } from 'vue'
import { alertStatus, type Alert } from '~/composables/useAlerts'

const props = withDefaults(defineProps<{ alert: Alert; actions?: boolean }>(), { actions: false })
const emit = defineEmits<{ pause: []; resume: []; remove: [] }>()

// Current rate is LIVE for the alert's direction (shared query/stream).
const { best } = useDirectionRates(
  () => props.alert.directionSlug,
  () => props.alert.reversed,
)

const status = computed(() => alertStatus(props.alert, best.value))
const isTriggered = computed(() => status.value === 'triggered')

const META = {
  triggered: { label: 'сработал', cls: 'bg-success text-[#06231A]' },
  active: { label: 'активен', cls: 'bg-surface-raised text-ink-muted border border-line-strong' },
  paused: { label: 'на паузе', cls: 'bg-surface-raised text-ink-faint border border-line-strong' },
} as const

const fmt = (n: number) => fmtNumber(n, n < 1000 ? 2 : 0)
const currentText = computed(() =>
  best.value == null ? '—' : `${fmt(best.value)} ${props.alert.unit}`,
)
</script>

<template>
  <div
    class="rounded-2xl border p-4 transition-colors"
    :class="[
      isTriggered ? 'border-success/30 bg-success/[0.08]' : 'border-line bg-surface',
      status === 'paused' ? 'opacity-60' : '',
    ]"
  >
    <div class="mb-2.5 flex items-center justify-between gap-3">
      <div class="flex min-w-0 items-center gap-2">
        <span
          class="flex h-[26px] w-[26px] flex-none items-center justify-center rounded-full text-[11px] font-bold"
          :style="{ backgroundColor: alert.badge.color, color: alert.badge.dark ? '#111' : '#fff' }"
          >{{ alert.badge.symbol }}</span
        >
        <span class="truncate text-[15px] font-semibold text-ink">{{ alert.pair }}</span>
      </div>
      <span
        class="flex-none rounded-md px-2.5 py-1 text-[11px] font-semibold"
        :class="META[status].cls"
        >{{ META[status].label }}</span
      >
    </div>

    <div class="flex items-center justify-between gap-3">
      <span class="text-[13px] text-ink-muted">
        {{ alert.direction === 'above' ? 'выше' : 'ниже' }}
        <span class="tnum text-ink">{{ fmt(alert.threshold) }} {{ alert.unit }}</span>
        · сейчас
        <span class="tnum" :class="isTriggered ? 'text-success-bright' : 'text-ink-muted'">{{
          currentText
        }}</span>
      </span>
      <span class="flex flex-none items-center gap-2">
        <svg
          v-if="alert.channels.includes('telegram')"
          class="text-brand-bright"
          width="13"
          height="13"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path d="M22 3 2 10.5l6 2.2M22 3l-3 17-8-6.3" />
        </svg>
        <svg
          v-if="alert.channels.includes('email')"
          class="text-ink-faint"
          width="13"
          height="13"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <rect x="3" y="5" width="18" height="14" rx="2" />
          <path d="m3 7 9 6 9-6" />
        </svg>
      </span>
    </div>

    <div v-if="actions" class="mt-3.5 flex gap-2 border-t border-line-subtle pt-3.5">
      <button
        v-if="status === 'paused'"
        type="button"
        class="flex-1 rounded-lg border border-line-strong bg-surface-raised py-2 text-[13px] font-semibold text-ink transition-colors hover:border-[#3A4047]"
        @click="emit('resume')"
      >
        Возобновить
      </button>
      <button
        v-else
        type="button"
        class="flex-1 rounded-lg border border-line-strong bg-surface-raised py-2 text-[13px] font-semibold text-ink-muted transition-colors hover:text-ink"
        @click="emit('pause')"
      >
        Приостановить
      </button>
      <button
        type="button"
        class="flex-1 rounded-lg border border-danger/30 bg-danger/[0.08] py-2 text-[13px] font-semibold text-danger transition-colors hover:bg-danger/15"
        @click="emit('remove')"
      >
        Удалить
      </button>
    </div>
  </div>
</template>
