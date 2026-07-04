<script setup lang="ts">
import { computed } from 'vue'

definePageMeta({ layout: false, middleware: 'auth' })
useSeoMeta({ title: 'Алерт сработал — Kurso' })

const route = useRoute()
const id = computed(() => route.params.id as string)
const { find, pause } = useAlerts()
const { record } = useHistory()
const apiBase = useApiBase()

const alert = computed(() => find(id.value))
const { ranked } = useDirectionRates(
  () => alert.value?.directionSlug ?? '',
  () => alert.value?.reversed ?? false,
)

const unit = computed(() => alert.value?.unit ?? '₽')
const decimals = computed(() => (ranked.value[0] ? rateDecimals(ranked.value[0].rate) : 2))

const best = computed(() => ranked.value[0] ?? null)
const bestRate = computed(() => best.value?.rate ?? null)
const overBy = computed(() =>
  bestRate.value != null && alert.value ? bestRate.value - alert.value.threshold : null,
)
// Alternatives above the threshold (excluding the best), up to 4.
const alternatives = computed(() => {
  const a = alert.value
  if (!a) return []
  return ranked.value
    .slice(1)
    .filter((r) => (a.direction === 'above' ? r.rate >= a.threshold : r.rate <= a.threshold))
    .slice(0, 4)
})

const fmt = (n: number) => fmtNumber(n, n < 1000 ? decimals.value : 0)

function goHref(slug: string) {
  return `${apiBase}/go/${slug}?direction=${alert.value?.directionSlug}`
}
function jump(slug: string, name: string) {
  if (!alert.value || bestRate.value == null) return
  record({
    pair: alert.value.pair,
    exchanger: name,
    slug,
    directionSlug: alert.value.directionSlug,
    amount: `${fmt(bestRate.value)} ${unit.value}`,
  })
}
function pauseAndBack() {
  if (alert.value) pause(alert.value.id)
  navigateTo('/account/alerts')
}
</script>

<template>
  <div class="relative min-h-[100dvh] overflow-x-clip bg-canvas">
    <!-- not found -->
    <div
      v-if="!alert"
      class="flex min-h-[100dvh] flex-col items-center justify-center gap-4 px-6 text-center"
    >
      <p class="text-ink-muted">Алерт не найден или был удалён.</p>
      <KButton @click="navigateTo('/account/alerts')">К моим алертам</KButton>
    </div>

    <div v-else class="mx-auto max-w-[520px]">
      <!-- green hero -->
      <div
        class="bg-[linear-gradient(160deg,#1FB37A,#2BC58C)] px-5 pb-7 pt-10 text-center text-[#06231A]"
      >
        <span
          class="relative mx-auto mb-3.5 flex h-16 w-16 items-center justify-center rounded-full bg-white/25"
        >
          <span class="absolute h-16 w-16 animate-kping rounded-full bg-white/40" />
          <svg
            class="relative"
            width="30"
            height="30"
            viewBox="0 0 24 24"
            fill="none"
            stroke="#06231A"
            stroke-width="2.2"
          >
            <circle cx="12" cy="12" r="8" />
            <circle cx="12" cy="12" r="3" />
          </svg>
        </span>
        <div class="text-[22px] font-extrabold tracking-[-0.02em]">Порог достигнут</div>
        <div class="mt-1 text-[15px] font-semibold opacity-80">{{ alert.pair }}</div>
      </div>

      <div class="px-4 py-5">
        <!-- best result -->
        <div v-if="best" class="mb-4 rounded-[18px] border border-success bg-surface p-[18px]">
          <div class="mb-3.5 flex items-center gap-3">
            <span
              class="flex h-11 w-11 flex-none items-center justify-center rounded-[13px] text-[15px] font-extrabold text-white"
              :style="{
                backgroundColor: exchangerAvatar(best.row.exchangerSlug, best.row.exchangerName)
                  .color,
              }"
              >{{ exchangerAvatar(best.row.exchangerSlug, best.row.exchangerName).initials }}</span
            >
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <span class="truncate text-base font-bold text-ink">{{
                  best.row.exchangerName
                }}</span>
                <KBadge v-if="best.row.partner" tone="brand">Партнёр</KBadge>
              </div>
              <div v-if="best.row.ratingAvg != null" class="tnum mt-0.5 text-xs text-ink-muted">
                <span class="text-warning-deep">★</span> {{ best.row.ratingAvg.toFixed(1) }} ·
                {{ best.row.reviewsCount }}
              </div>
            </div>
          </div>
          <div class="mb-3.5 flex items-end gap-2">
            <span class="tnum text-3xl font-bold text-ink">{{ fmt(best.rate) }} {{ unit }}</span>
            <span
              v-if="overBy != null && overBy >= 0"
              class="tnum mb-0.5 text-xs font-semibold text-success-bright"
            >
              +{{ fmt(overBy) }} {{ unit }} над порогом
            </span>
          </div>
          <div class="mb-3.5 flex flex-wrap gap-2 text-xs text-ink-muted">
            <span
              v-if="best.row.minAmount"
              class="tnum rounded-lg border border-line bg-well px-2.5 py-1.5"
              >мин {{ fmtCompact(Number(best.row.minAmount)) }} {{ unit }}</span
            >
            <span
              v-if="best.row.maxAmount"
              class="tnum rounded-lg border border-line bg-well px-2.5 py-1.5"
              >макс {{ fmtCompact(Number(best.row.maxAmount)) }} {{ unit }}</span
            >
            <span
              v-if="best.row.reserve"
              class="tnum rounded-lg border border-line bg-well px-2.5 py-1.5"
              >резерв {{ fmtCompact(Number(best.row.reserve)) }} {{ unit }}</span
            >
          </div>
          <KButton
            block
            size="lg"
            class="!rounded-[14px]"
            :href="goHref(best.row.exchangerSlug)"
            @click="jump(best.row.exchangerSlug, best.row.exchangerName)"
          >
            Перейти в {{ best.row.exchangerName }}
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.4"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M7 17 17 7M9 7h8v8" />
            </svg>
          </KButton>
          <p class="mt-2.5 text-center text-xs text-ink-faint">
            Курс актуален ~3 минуты · затем пересчитается
          </p>
        </div>

        <div
          v-else
          class="mb-4 rounded-[18px] border border-line bg-surface p-6 text-center text-sm text-ink-muted"
        >
          Курсы по этой паре сейчас недоступны. Попробуйте обновить страницу.
        </div>

        <!-- alternatives -->
        <template v-if="alternatives.length">
          <div class="mb-2.5 text-sm font-bold text-ink">
            Ещё {{ alternatives.length }}
            {{ pluralRu(alternatives.length, 'обменник', 'обменника', 'обменников') }} выше порога
          </div>
          <div class="mb-[18px] flex flex-col gap-2.5">
            <a
              v-for="alt in alternatives"
              :key="alt.row.exchangerId"
              :href="goHref(alt.row.exchangerSlug)"
              target="_blank"
              rel="noopener nofollow"
              class="flex items-center gap-2.5 rounded-2xl border border-line bg-surface px-3.5 py-3 transition-colors hover:border-line-strong"
              @click="jump(alt.row.exchangerSlug, alt.row.exchangerName)"
            >
              <span
                class="flex h-8 w-8 flex-none items-center justify-center rounded-full text-[11px] font-bold text-white"
                :style="{
                  backgroundColor: exchangerAvatar(alt.row.exchangerSlug, alt.row.exchangerName)
                    .color,
                }"
                >{{ exchangerAvatar(alt.row.exchangerSlug, alt.row.exchangerName).initials }}</span
              >
              <div class="min-w-0 flex-1">
                <div class="truncate text-sm font-semibold text-ink">
                  {{ alt.row.exchangerName }}
                </div>
                <div v-if="alt.row.reserve" class="tnum text-[11px] text-ink-faint">
                  резерв {{ fmtCompact(Number(alt.row.reserve)) }} {{ unit }}
                </div>
              </div>
              <span class="tnum flex-none text-[15px] font-bold text-ink"
                >{{ fmt(alt.rate) }} {{ unit }}</span
              >
            </a>
          </div>
        </template>

        <!-- actions -->
        <div class="flex gap-2.5">
          <KButton variant="secondary" block @click="navigateTo('/')">Показать все</KButton>
          <button
            type="button"
            class="flex-1 rounded-[14px] border border-line-strong bg-surface-raised py-3.5 text-sm font-semibold text-ink-muted transition-colors hover:text-ink"
            @click="pauseAndBack"
          >
            Приостановить
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
