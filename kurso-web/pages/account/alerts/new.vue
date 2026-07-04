<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import type { AlertChannel } from '~/composables/useAlerts'

definePageMeta({ layout: false, middleware: 'auth' })
useSeoMeta({ title: 'Новый алерт — Kurso' })

const { give, get, giveCode, directionSlug, reversed, unitRate } = useCalculator()
const { openPicker } = useCurrencyPicker()
const { add } = useAlerts()

const unit = computed(() => currencyUnit(get.value.code))
const decimals = computed(() => (unitRate.value != null ? rateDecimals(unitRate.value) : 2))

// --- threshold ---
const direction = ref<'above' | 'below'>('above')
const threshold = ref<number>(0)
const touched = ref(false)
const activePct = ref<number | null>(1)

function applyPct(pct: number) {
  activePct.value = pct
  touched.value = true
  if (unitRate.value != null) {
    const raw = unitRate.value * (1 + (direction.value === 'above' ? pct : -pct) / 100)
    threshold.value = round(raw, decimals.value)
  }
}
const round = (n: number, d: number) => {
  const f = 10 ** d
  return Math.round(n * f) / f
}

// Seed the threshold from the live best rate once it arrives (if untouched).
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

// --- other fields ---
const scope = ref<'all' | 'favorites' | 'manual'>('all')
const minReserve = ref(true)
const ratingFilter = ref(true)
const channels = ref<Record<AlertChannel, boolean>>({ telegram: true, email: false, push: false })
const validity = ref<'7' | '30' | 'forever'>('30')

const { alerts } = useAlerts()

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

const conditionOpts = [
  { value: 'above', label: 'Курс выше' },
  { value: 'below', label: 'Ниже' },
]
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
    <div class="relative z-10 mx-auto flex min-h-[100dvh] max-w-[560px] flex-col">
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

      <div class="flex flex-col gap-[18px] px-4 py-5 md:px-0">
        <!-- pair -->
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Пара</div>
          <div class="flex items-center gap-2.5">
            <button
              type="button"
              class="flex flex-1 items-center gap-2.5 rounded-2xl border border-line bg-surface px-3 py-3 transition-colors hover:border-line-strong"
              @click="openPicker('give')"
            >
              <span
                class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-full text-[13px] font-bold"
                :style="{ backgroundColor: give.color, color: give.dark ? '#111' : '#fff' }"
                >{{ give.symbol }}</span
              >
              <span class="flex-1 truncate text-left text-[15px] font-semibold text-ink">{{
                give.name
              }}</span>
              <svg
                class="text-ink-faint"
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
            <svg
              class="flex-none text-ink-faint"
              width="18"
              height="18"
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
              class="flex flex-1 items-center gap-2.5 rounded-2xl border border-line bg-surface px-3 py-3 transition-colors hover:border-line-strong"
              @click="openPicker('get')"
            >
              <span
                class="flex h-[30px] w-[30px] flex-none items-center justify-center rounded-full text-[13px] font-bold"
                :style="{ backgroundColor: get.color, color: get.dark ? '#111' : '#fff' }"
                >{{ get.symbol }}</span
              >
              <span class="flex-1 truncate text-left text-[15px] font-semibold text-ink">{{
                get.name
              }}</span>
              <svg
                class="text-ink-faint"
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
          </div>
          <p v-if="!directionSlug" class="mt-2 text-xs text-warning">
            Для этой пары пока нет курсов — выберите крипту и фиат.
          </p>
        </div>

        <!-- condition -->
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Условие</div>
          <KSegmented v-model="direction" :options="conditionOpts">
            <template #icon="{ option, active }">
              <svg
                v-if="option.value === 'above'"
                width="15"
                height="15"
                viewBox="0 0 24 24"
                fill="none"
                :stroke="active ? '#2BC58C' : 'currentColor'"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M5 14l7-7 7 7" />
              </svg>
              <svg
                v-else
                width="15"
                height="15"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.2"
                stroke-linecap="round"
                stroke-linejoin="round"
              >
                <path d="M5 10l7 7 7-7" />
              </svg>
            </template>
          </KSegmented>
        </div>

        <!-- threshold -->
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Порог</div>
          <div
            class="flex items-baseline gap-2 rounded-2xl border border-brand bg-surface px-4 py-3.5"
          >
            <input
              v-model="thresholdStr"
              inputmode="decimal"
              class="tnum w-[6ch] min-w-0 flex-none bg-transparent text-3xl font-bold text-ink focus:outline-none"
              :style="{ width: `${Math.max(4, thresholdStr.length + 1)}ch` }"
            />
            <span class="text-lg font-semibold text-ink-faint">{{ unit }}</span>
            <span v-if="unitRate != null" class="tnum ml-auto text-[13px] text-ink-faint"
              >сейчас лучший
              <span class="text-ink-muted"
                >{{ fmtNumber(unitRate, decimals) }} {{ unit }}</span
              ></span
            >
          </div>
          <div class="mt-2.5 flex gap-2">
            <button
              v-for="p in PCTS"
              :key="p"
              type="button"
              class="flex-1 rounded-full py-2 text-[13px] font-semibold transition-colors"
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

        <!-- exchangers scope -->
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Обменники</div>
          <KSegmented v-model="scope" :options="scopeOpts" />
        </div>

        <!-- extra filters -->
        <div class="flex flex-col gap-2.5">
          <div
            class="flex items-center justify-between rounded-2xl border border-line bg-surface px-4 py-3.5"
          >
            <span class="text-sm text-ink">Только с резервом от 5M ₽</span>
            <KToggle v-model="minReserve" />
          </div>
          <div
            class="flex items-center justify-between rounded-2xl border border-line bg-surface px-4 py-3.5"
          >
            <span class="text-sm text-ink">Учитывать рейтинг ≥ 4.5</span>
            <KToggle v-model="ratingFilter" />
          </div>
        </div>

        <!-- channels -->
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Куда уведомить</div>
          <div class="flex flex-col gap-2.5">
            <div
              class="flex items-center gap-3 rounded-2xl border border-line bg-surface px-4 py-3"
            >
              <span
                class="flex h-[34px] w-[34px] flex-none items-center justify-center rounded-[10px] bg-brand text-white"
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
            </div>
            <div
              class="flex items-center gap-3 rounded-2xl border border-line bg-surface px-4 py-3"
            >
              <span
                class="flex h-[34px] w-[34px] flex-none items-center justify-center rounded-[10px] bg-[#3A4452] text-white"
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
                  <rect x="3" y="5" width="18" height="14" rx="2" />
                  <path d="m3 7 9 6 9-6" />
                </svg>
              </span>
              <span class="flex-1 text-[15px] font-semibold text-ink">Email</span>
              <KToggle v-model="channels.email" />
            </div>
            <div
              class="flex items-center gap-3 rounded-2xl border border-line bg-surface px-4 py-3"
            >
              <span
                class="flex h-[34px] w-[34px] flex-none items-center justify-center rounded-[10px] bg-[#3A4452] text-white"
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
                  <path d="M18 8a6 6 0 1 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
                  <path d="M10.5 20a1.8 1.8 0 0 0 3 0" />
                </svg>
              </span>
              <span class="flex-1 text-[15px] font-semibold text-ink">Push</span>
              <KToggle v-model="channels.push" />
            </div>
          </div>
        </div>

        <!-- validity -->
        <div>
          <div class="mb-2 text-xs uppercase tracking-[0.05em] text-ink-faint">Срок действия</div>
          <KSegmented v-model="validity" :options="validityOpts" />
        </div>

        <KButton block size="lg" class="!rounded-2xl" :disabled="!canCreate" @click="create">
          Создать алерт
        </KButton>
        <p class="pb-6 text-center text-xs leading-relaxed text-ink-faint">
          До <span class="tnum">10</span> активных алертов на бесплатном тарифе · сейчас
          <span class="tnum">{{ alerts.length }}</span>
        </p>
      </div>
    </div>

    <CurrencyPickerDrawer />
  </div>
</template>
