<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { AlertChannel } from '~/composables/useAlerts'

// Alert builder — "smart alert" variant: the rule reads as one sentence, backed
// by a live best-rate chart and a forecast computed from REAL observed history
// (how often the threshold was actually reached). Creating an alert persists it;
// the account/triggered screens evaluate it live against the rates API.
definePageMeta({ layout: false, middleware: 'auth' })
useSeoMeta({ title: 'Новый алерт — Kurso' })

const { give, get, giveCode, directionSlug, reversed, unitRate } = useCalculator()
const { openPicker } = useCurrencyPicker()
const { add, alerts } = useAlerts()

const unit = computed(() => currencyUnit(get.value.code))
const decimals = computed(() => (unitRate.value != null ? rateDecimals(unitRate.value) : 2))

// --- threshold ---
const direction = ref<'above' | 'below'>('above')
const threshold = ref<number>(0)
const touched = ref(false)
const activePct = ref<number | null>(1)

const round = (n: number, d: number) => {
  const f = 10 ** d
  return Math.round(n * f) / f
}
function applyPct(pct: number) {
  activePct.value = pct
  touched.value = true
  if (unitRate.value != null) {
    const raw = unitRate.value * (1 + (direction.value === 'above' ? pct : -pct) / 100)
    threshold.value = round(raw, decimals.value)
  }
}
watch(
  [unitRate, direction],
  () => {
    if (!touched.value && unitRate.value != null) applyPct(activePct.value ?? 1)
  },
  { immediate: true },
)

const thresholdStr = computed({
  get: () => (threshold.value ? fmtNumber(threshold.value, decimals.value) : ''),
  set: (v: string) => {
    touched.value = true
    activePct.value = null
    const n = parseFloat(
      v
        .replace(/[^\d.,]/g, '')
        .replace(/\s/g, '')
        .replace(',', '.'),
    )
    threshold.value = Number.isFinite(n) ? n : 0
  },
})
const PCTS = [1, 1.25, 2, 5]

// --- live chart history (real, from /rates/{dir}/history) ---
const { series: rawSeries } = useRateHistory(directionSlug, 7)
// Match the displayed unit: a fiat→crypto selection inverts the direction rate.
const series = computed(() =>
  reversed.value
    ? rawSeries.value.map((r) => (r > 0 ? 1 / r : 0)).filter((n) => n > 0)
    : rawSeries.value,
)

// --- forecast from real observed history ---
const movePct = computed(() => {
  if (unitRate.value == null || unitRate.value === 0) return null
  const diff =
    direction.value === 'above'
      ? threshold.value - unitRate.value
      : unitRate.value - threshold.value
  return (diff / unitRate.value) * 100
})
const reachedCount = computed(
  () =>
    series.value.filter((v) =>
      direction.value === 'above' ? v >= threshold.value : v <= threshold.value,
    ).length,
)
function plural(n: number, one: string, few: string, many: string) {
  const m10 = n % 10
  const m100 = n % 100
  if (m10 === 1 && m100 !== 11) return one
  if (m10 >= 2 && m10 <= 4 && (m100 < 10 || m100 >= 20)) return few
  return many
}

// --- filters + channels ---
const scope = ref<'all' | 'favorites' | 'manual'>('all')
const minReserve = ref(true)
const ratingFilter = ref(true)
const validity = ref<'7' | '30' | 'forever'>('30')
const channels = ref<Record<AlertChannel, boolean>>({ telegram: true, email: false, push: false })
const showFilters = ref(false)

const scopeLabel = computed(
  () => ({ all: 'Все обменники', favorites: 'Избранные', manual: 'Вручную' })[scope.value],
)
const validityLabel = computed(
  () => ({ '7': '7 дней', '30': '30 дней', forever: 'бессрочно' })[validity.value],
)

const canCreate = computed(() => !!directionSlug.value && threshold.value > 0)

function create() {
  if (!canCreate.value) return
  const selected = (Object.keys(channels.value) as AlertChannel[]).filter((c) => channels.value[c])
  add({
    directionSlug: directionSlug.value,
    reversed: reversed.value,
    pair: `${giveCode.value} → ${get.value.name}`,
    badge: { symbol: give.value.symbol, color: give.value.color, dark: give.value.dark },
    unit: unit.value,
    direction: direction.value,
    threshold: threshold.value,
    channels: selected.length ? selected : ['telegram'],
    scope: scope.value,
    minReserve: minReserve.value,
    ratingFilter: ratingFilter.value,
    validity: validity.value,
  })
  navigateTo('/account/alerts')
}

const scopeOpts = [
  { value: 'all', label: 'Все' },
  { value: 'favorites', label: 'Избранные' },
  { value: 'manual', label: 'Вручную' },
]
const validityOpts = [
  { value: '7', label: '7 дней' },
  { value: '30', label: '30 дней' },
  { value: 'forever', label: 'Бессрочно' },
]
</script>

<template>
  <div class="relative min-h-[100dvh] overflow-x-clip bg-canvas">
    <GlowBackdrop />
    <div class="relative z-10 mx-auto flex min-h-[100dvh] max-w-[880px] flex-col">
      <!-- nav bar -->
      <div class="flex items-center justify-between border-b border-line-faint px-4 py-3.5">
        <button
          type="button"
          class="text-[15px] text-ink-muted"
          @click="navigateTo('/account/alerts')"
        >
          Отмена
        </button>
        <span class="text-base font-bold text-ink">Новый алерт</span>
        <button
          type="button"
          class="text-[15px] font-semibold text-brand-bright disabled:opacity-40"
          :disabled="!canCreate"
          @click="create"
        >
          Создать
        </button>
      </div>

      <div class="grid gap-4 px-4 py-5 md:grid-cols-2 md:gap-6 md:px-0">
        <!-- LEFT: rule + channels -->
        <div class="flex flex-col gap-4">
          <!-- rule sentence -->
          <div class="rounded-2xl border border-line bg-surface p-4">
            <div class="mb-2.5 text-xs uppercase tracking-[0.05em] text-ink-faint">Правило</div>
            <div
              class="flex flex-wrap items-center gap-x-2 gap-y-2.5 text-[17px] leading-relaxed text-ink"
            >
              <span>Уведомить, когда</span>
              <button
                type="button"
                class="inline-flex items-center gap-1.5 rounded-full border border-line-strong bg-surface-chip px-2.5 py-1.5 text-[14px] font-semibold transition-colors hover:border-brand"
                @click="openPicker('give')"
              >
                <span
                  class="flex h-5 w-5 items-center justify-center rounded-full text-[11px] font-bold"
                  :style="{ backgroundColor: give.color, color: give.dark ? '#111' : '#fff' }"
                  >{{ give.symbol }}</span
                >{{ give.code }}
              </button>
              <svg
                class="text-ink-faint"
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M5 12h14M13 6l6 6-6 6" />
              </svg>
              <button
                type="button"
                class="inline-flex items-center gap-1.5 rounded-full border border-line-strong bg-surface-chip px-2.5 py-1.5 text-[14px] font-semibold transition-colors hover:border-brand"
                @click="openPicker('get')"
              >
                <span
                  class="flex h-5 w-5 items-center justify-center rounded-full text-[11px] font-bold"
                  :style="{ backgroundColor: get.color, color: get.dark ? '#111' : '#fff' }"
                  >{{ get.symbol }}</span
                >{{ get.name }}
              </button>
              <span>станет</span>
              <button
                type="button"
                class="inline-flex items-center gap-1 rounded-full border border-line-strong bg-surface-chip px-2.5 py-1.5 text-[14px] font-semibold transition-colors hover:border-brand"
                @click="direction = direction === 'above' ? 'below' : 'above'"
              >
                <svg
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  :stroke="direction === 'above' ? '#2BC58C' : '#EC7B7A'"
                  stroke-width="2.2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path :d="direction === 'above' ? 'M5 14l7-7 7 7' : 'M5 10l7 7 7-7'" />
                </svg>
                {{ direction === 'above' ? 'выше' : 'ниже' }}
              </button>
              <span
                class="inline-flex items-baseline gap-1 rounded-full border border-brand bg-brand/[0.08] px-2.5 py-1"
              >
                <input
                  v-model="thresholdStr"
                  inputmode="decimal"
                  class="tnum min-w-0 flex-none bg-transparent text-[16px] font-bold text-ink focus:outline-none"
                  :style="{ width: `${Math.max(3, thresholdStr.length + 1)}ch` }"
                />
                <span class="text-[14px] font-semibold text-ink-faint">{{ unit }}</span>
              </span>
            </div>

            <p v-if="!directionSlug" class="mt-3 text-xs text-warning">
              Для этой пары пока нет курсов — выберите крипту и фиат.
            </p>

            <!-- percent presets -->
            <div class="mt-3.5 flex gap-2">
              <button
                v-for="p in PCTS"
                :key="p"
                type="button"
                class="flex-1 rounded-full py-1.5 text-[13px] font-semibold transition-colors"
                :class="
                  activePct === p
                    ? 'border border-line-strong bg-surface-chip text-ink'
                    : 'border border-line bg-surface text-ink-muted hover:text-ink'
                "
                :disabled="unitRate == null"
                @click="applyPct(p)"
              >
                {{ direction === 'above' ? '+' : '−' }}{{ p }}%
              </button>
            </div>
          </div>

          <!-- channels -->
          <div class="rounded-2xl border border-line bg-surface p-4">
            <div class="mb-2.5 text-xs uppercase tracking-[0.05em] text-ink-faint">
              Куда уведомить
            </div>
            <div class="flex flex-col gap-2.5">
              <label class="flex items-center gap-3">
                <span
                  class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-brand text-white"
                >
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M22 3 2 10.5l6 2.2M22 3l-3 17-8-6.3M22 3 8 12.7m0 0v5.3l3-3.6" />
                  </svg>
                </span>
                <span class="flex flex-1 items-center gap-2 text-[15px] font-semibold text-ink"
                  >Telegram<span
                    class="rounded bg-brand/15 px-1.5 py-0.5 text-[10px] font-semibold text-brand-bright"
                    >рекомендуем</span
                  ></span
                >
                <KToggle v-model="channels.telegram" />
              </label>
              <label class="flex items-center gap-3">
                <span
                  class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-[#3A4452] text-white"
                >
                  <svg
                    width="16"
                    height="16"
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
                <span class="flex-1 text-[15px] font-semibold text-ink">Email</span>
                <KToggle v-model="channels.email" />
              </label>
              <label class="flex items-center gap-3">
                <span
                  class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] bg-[#3A4452] text-white"
                >
                  <svg
                    width="16"
                    height="16"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
                    <path d="M10.5 20a1.8 1.8 0 0 0 3 0" />
                  </svg>
                </span>
                <span class="flex-1 text-[15px] font-semibold text-ink">Push</span>
                <KToggle v-model="channels.push" />
              </label>
            </div>
          </div>

          <!-- filters summary (collapsible) -->
          <div class="rounded-2xl border border-line bg-surface">
            <button
              type="button"
              class="flex w-full items-center gap-2 px-4 py-3.5 text-left"
              @click="showFilters = !showFilters"
            >
              <span class="flex-1 truncate text-[13px] text-ink-muted">
                {{ scopeLabel }}<template v-if="minReserve"> · резерв 5M+</template
                ><template v-if="ratingFilter"> · ★ ≥ 4.5</template> · {{ validityLabel }}
              </span>
              <span class="text-[13px] font-semibold text-brand-bright">{{
                showFilters ? 'Скрыть' : 'Настроить'
              }}</span>
            </button>
            <div v-if="showFilters" class="flex flex-col gap-3 border-t border-line px-4 py-4">
              <div>
                <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Обменники</div>
                <KSegmented v-model="scope" :options="scopeOpts" />
              </div>
              <label class="flex items-center justify-between">
                <span class="text-sm text-ink">Только с резервом от 5M ₽</span>
                <KToggle v-model="minReserve" />
              </label>
              <label class="flex items-center justify-between">
                <span class="text-sm text-ink">Учитывать рейтинг ≥ 4.5</span>
                <KToggle v-model="ratingFilter" />
              </label>
              <div>
                <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">
                  Срок действия
                </div>
                <KSegmented v-model="validity" :options="validityOpts" />
              </div>
            </div>
          </div>
        </div>

        <!-- RIGHT: chart + forecast -->
        <div class="flex flex-col gap-4">
          <div class="rounded-2xl border border-line bg-surface p-4">
            <div class="mb-3 flex items-center justify-between">
              <span class="text-[13px] font-semibold text-ink">Курс · 7 дней</span>
              <span v-if="unitRate != null" class="tnum text-[13px] text-ink-faint"
                >сейчас
                <span class="text-ink-muted"
                  >{{ fmtNumber(unitRate, decimals) }} {{ unit }}</span
                ></span
              >
            </div>
            <AlertChart :series="series" :threshold="threshold" :current="unitRate" />
            <div class="mt-2 flex items-center gap-2">
              <span class="inline-block h-0 w-4 border-t-2 border-dashed border-brand-bright" />
              <span class="tnum text-xs font-semibold text-brand-bright"
                >порог {{ thresholdStr }} {{ unit }}</span
              >
            </div>
          </div>

          <!-- forecast (real) -->
          <div class="rounded-2xl border border-brand/25 bg-brand/[0.07] p-4">
            <div class="mb-1.5 flex items-center gap-2">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="#6BA6FF"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <circle cx="12" cy="12" r="9" />
                <path d="M12 7v5l3 3" />
              </svg>
              <span class="text-[13px] font-semibold text-brand-bright">Прогноз срабатывания</span>
            </div>
            <p class="text-[13px] leading-relaxed text-ink-muted">
              <template v-if="unitRate == null">Курс загружается…</template>
              <template v-else-if="movePct != null && movePct <= 0">
                Порог уже достигнут — алерт сработает при следующей проверке.
              </template>
              <template v-else-if="movePct != null">
                Нужен {{ direction === 'above' ? 'рост' : 'спад' }} на
                <span class="font-semibold text-ink">{{ fmtNumber(movePct, 1) }}%</span>.
                <template v-if="reachedCount > 0">
                  За неделю курс достигал
                  <span class="tnum font-semibold text-ink">{{ thresholdStr }} {{ unit }}</span>
                  {{ reachedCount }} {{ plural(reachedCount, 'раз', 'раза', 'раз') }}.
                </template>
                <template v-else>За последнюю неделю курс не доходил до порога.</template>
              </template>
            </p>
          </div>

          <KButton block size="lg" class="!rounded-2xl" :disabled="!canCreate" @click="create">
            Создать алерт
          </KButton>
          <p class="text-center text-xs leading-relaxed text-ink-faint">
            До <span class="tnum">{{ ALERT_LIMIT }}</span> активных алертов на бесплатном тарифе ·
            сейчас
            <span class="tnum">{{ alerts.length }}</span>
          </p>
        </div>
      </div>
    </div>

    <CurrencyPickerDrawer />
  </div>
</template>
