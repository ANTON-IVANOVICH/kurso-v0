<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const {
  give,
  get,
  amount,
  directionSlug,
  hasRates,
  bestRow,
  unitRate,
  receiveAmount,
  exchangerCount,
  giveCode,
  giveUnit,
  receiveUnit,
  swap,
} = useCalculator()
const { openPicker } = useCurrencyPicker()
const { t, plural } = useI18n()

// --- editable give amount ---
// A real exchange-calculator field: it accepts ONLY digits and a single decimal
// separator, groups thousands live as you type, and silently drops anything else
// — letters and symbols can never appear. `amount` is the shared numeric truth.
const display = ref(fmtNumber(amount.value, currencyDecimals(giveCode.value)))
const focused = ref(false)

// Sanitize a raw field value to {formatted, value}: strip non-numerics, keep one
// dot, clamp decimals to the currency, drop leading zeros, group thousands.
function sanitize(raw: string, dec: number): { formatted: string; value: number } {
  const s = raw.replace(/[^\d.,]/g, '').replace(/,/g, '.')
  const firstDot = s.indexOf('.')
  let intPart =
    firstDot === -1 || dec === 0
      ? firstDot === -1
        ? s
        : s.slice(0, firstDot)
      : s.slice(0, firstDot)
  let fracPart = ''
  const hasDot = firstDot !== -1 && dec > 0
  if (hasDot)
    fracPart = s
      .slice(firstDot + 1)
      .replace(/\./g, '')
      .slice(0, dec)
  intPart = intPart.replace(/\D/g, '').replace(/^0+(?=\d)/, '')
  if (intPart === '' && !hasDot) return { formatted: '', value: 0 }
  const groupedInt = (intPart === '' ? '0' : intPart).replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  const formatted = hasDot ? `${groupedInt}.${fracPart}` : groupedInt
  const value =
    parseFloat(`${intPart === '' ? '0' : intPart}${hasDot ? `.${fracPart || '0'}` : ''}`) || 0
  return { formatted, value }
}
function countDigits(s: string): number {
  return (s.match(/\d/g) || []).length
}
function caretForDigits(formatted: string, n: number): number {
  if (n <= 0) return 0
  let seen = 0
  for (let i = 0; i < formatted.length; i++) {
    const c = formatted.charCodeAt(i)
    if (c >= 48 && c <= 57 && ++seen === n) return i + 1
  }
  return formatted.length
}

function onInput(e: Event) {
  const el = e.target as HTMLInputElement
  const dec = currencyDecimals(giveCode.value)
  const digitsBeforeCaret = countDigits(el.value.slice(0, el.selectionStart ?? el.value.length))
  const { formatted, value } = sanitize(el.value, dec)
  amount.value = value
  display.value = formatted
  // Force the DOM value even when Vue would skip the patch (bound value
  // unchanged), so a rejected character is never left visible; restore caret.
  el.value = formatted
  const caret = caretForDigits(formatted, digitsBeforeCaret)
  el.setSelectionRange(caret, caret)
}
function onFocus(e: Event) {
  focused.value = true
  // Drop the padded ".00" so the field is comfortable to edit.
  const { formatted } = sanitize(String(amount.value), currencyDecimals(giveCode.value))
  display.value = formatted
  ;(e.target as HTMLInputElement).value = formatted
}
function onBlur() {
  focused.value = false
  display.value = fmtNumber(amount.value, currencyDecimals(giveCode.value))
}

// Reflect external changes (swap, currency switch) when not actively editing.
watch([amount, giveCode], () => {
  if (!focused.value) display.value = fmtNumber(amount.value, currencyDecimals(giveCode.value))
})

// --- receive side ---
const receiveText = computed(() =>
  receiveAmount.value == null
    ? '—'
    : `${fmtNumber(receiveAmount.value, currencyDecimals(get.value.code))} ${receiveUnit.value}`,
)
const rateHint = computed(() => {
  if (unitRate.value == null) return t('home.calc.byRate')
  return t('home.calc.rateHint', {
    from: giveUnit.value,
    rate: fmtNumber(unitRate.value, rateDecimals(unitRate.value)),
    to: receiveUnit.value,
  })
})

// --- live freshness (client-only clock; avoids SSR hydration drift) ---
const now = ref(0)
let timer: ReturnType<typeof setInterval> | undefined
onMounted(() => {
  now.value = Date.now()
  timer = setInterval(() => (now.value = Date.now()), 1000)
})
onBeforeUnmount(() => clearInterval(timer))

const freshness = computed(() => {
  if (!now.value || !bestRow.value) return ''
  const ts = new Date(bestRow.value.fetchedAt).getTime()
  if (!Number.isFinite(ts)) return ''
  const sec = Math.max(0, Math.round((now.value - ts) / 1000))
  if (sec < 5) return t('home.calc.justNow')
  if (sec < 60) return t('home.calc.secAgo', { n: sec })
  const min = Math.floor(sec / 60)
  if (min < 60) return t('home.calc.minAgo', { n: min })
  return t('home.calc.hourAgo', { n: Math.floor(min / 60) })
})

const exchangersWord = computed(() => plural(exchangerCount.value, 'exchangers'))
</script>

<template>
  <div class="rounded-3xl border border-line bg-surface p-5">
    <div class="mb-4 font-semibold">{{ t('home.calc.title') }}</div>

    <div class="relative">
      <!-- Отдаю -->
      <div class="rounded-2xl border border-line bg-well p-[15px]">
        <div class="mb-[7px] text-xs text-ink-muted">{{ t('home.calc.give') }}</div>
        <div class="flex items-center justify-between gap-3">
          <input
            :value="display"
            inputmode="decimal"
            autocomplete="off"
            :aria-label="t('home.calc.amountAria')"
            class="tnum min-w-0 flex-1 bg-transparent text-[23px] font-bold text-ink focus:outline-none"
            @focus="onFocus"
            @input="onInput"
            @blur="onBlur"
          />
          <button
            type="button"
            class="inline-flex min-w-0 flex-none items-center gap-2 rounded-full bg-surface-raised py-1.5 pl-1.5 pr-3 transition-colors hover:bg-surface-chip"
            @click="openPicker('give')"
          >
            <span
              class="flex h-7 w-7 flex-none items-center justify-center rounded-full font-bold"
              :class="give.symbol.length > 1 ? 'text-[10px]' : 'text-[13px]'"
              :style="{ backgroundColor: give.color, color: give.dark ? '#111' : '#fff' }"
              >{{ give.symbol }}</span
            >
            <span class="min-w-0 truncate text-sm font-semibold">{{ give.name }}</span>
            <svg
              class="shrink-0 text-ink-faint"
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
      </div>

      <!-- swap -->
      <div class="relative z-10 -my-[13px] flex justify-center">
        <button
          type="button"
          class="flex h-11 w-11 items-center justify-center rounded-[15px] border-4 border-surface bg-surface-chip text-ink transition-colors hover:bg-surface-raised active:scale-95"
          :aria-label="t('home.calc.swapAria')"
          @click="swap"
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
            <path d="M7 4 L7 20 M7 20 L4 17 M7 20 L10 17 M17 20 L17 4 M17 4 L14 7 M17 4 L20 7" />
          </svg>
        </button>
      </div>

      <!-- Получаю -->
      <div class="rounded-2xl border border-line bg-well p-[15px]">
        <div class="mb-[7px] flex justify-between gap-3">
          <span class="text-xs text-ink-muted">{{ t('home.calc.receive') }}</span>
          <span class="truncate text-xs text-ink-faint">{{ rateHint }}</span>
        </div>
        <div class="flex items-center justify-between gap-3">
          <span
            class="tnum min-w-0 flex-1 truncate text-[23px] font-bold"
            :class="hasRates ? 'text-success' : 'text-ink-faint'"
            >{{ receiveText }}</span
          >
          <button
            type="button"
            class="inline-flex min-w-0 flex-none items-center gap-2 rounded-full bg-surface-raised py-1.5 pl-1.5 pr-3 transition-colors hover:bg-surface-chip"
            @click="openPicker('get')"
          >
            <span
              class="flex h-7 w-7 flex-none items-center justify-center rounded-full font-extrabold"
              :class="get.symbol.length > 1 ? 'text-[10px]' : 'text-[14px]'"
              :style="{ backgroundColor: get.color, color: get.dark ? '#111' : '#fff' }"
              >{{ get.symbol }}</span
            >
            <span class="min-w-0 truncate text-sm font-semibold">{{ get.name }}</span>
            <svg
              class="shrink-0 text-ink-faint"
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
      </div>
    </div>

    <div class="mt-4 flex items-center gap-2.5 text-[13px] text-ink-faint">
      <KStatusDot :tone="hasRates ? 'success' : 'neutral'" :pulse="hasRates" />
      <span v-if="hasRates" class="tnum">{{ exchangerCount }} {{ exchangersWord }}</span>
      <span v-else-if="!directionSlug">{{ t('home.calc.pairUnavailable') }}</span>
      <span v-else>{{ t('home.calc.updating') }}</span>
      <template v-if="freshness">
        <span class="text-line-strong">·</span>
        <span>{{ freshness }}</span>
      </template>
    </div>
  </div>
</template>
