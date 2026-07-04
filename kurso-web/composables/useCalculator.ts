import { computed } from 'vue'
import type { RateRow } from './useApi'

// Calculator currency codes → seeded direction-slug fragments.
const codeToSlug: Record<string, string> = {
  USDT: 'usdt',
  BTC: 'btc',
  ETH: 'eth',
  TON: 'ton',
  Тинькофф: 'tinkoff',
  Сбербанк: 'sber',
}

/** A rate row enriched with the amount the user actually receives for it. */
export interface RankedRate {
  row: RateRow
  /** Rate expressed as get-per-give (inverted for a fiat→crypto selection). */
  effectiveRate: number
  /** `amount` converted at this exchanger's rate. */
  receive: number
}

/**
 * Shared calculator state: the give/get currencies (via the picker), the amount,
 * and the live best rate for the resulting direction.
 *
 * Directions are seeded crypto→fiat with `rate` = fiat-per-crypto. A crypto→fiat
 * selection uses the rate directly; a fiat→crypto selection queries the same slug
 * and inverts the rate. crypto↔crypto and fiat↔fiat pairs have no direction and
 * yield an empty result.
 *
 * Both the calculator card and the home page's exchanger list consume this, so
 * the "best" rate in the calculator always matches the top card. The underlying
 * Pina Colada query is deduped by key and the SSE stream is ref-counted, so the
 * multiple call sites share one fetch and one connection.
 */
export function useCalculator() {
  const { give, get, swap: swapCurrencies } = useCurrencyPicker()
  const amount = useState<number>('calcAmount', () => 1000)

  const giveIsCrypto = computed(() => !isFiatCode(give.value.code))
  const getIsCrypto = computed(() => !isFiatCode(get.value.code))
  // Only crypto↔fiat pairs map to a seeded direction; `reversed` is the
  // fiat→crypto case, which reuses the crypto→fiat slug with an inverted rate.
  const reversed = computed(() => !giveIsCrypto.value && getIsCrypto.value)

  const directionSlug = computed(() => {
    const g = codeToSlug[give.value.code]
    const t = codeToSlug[get.value.code]
    if (!g || !t) return ''
    if (giveIsCrypto.value && !getIsCrypto.value) return `${g}-${t}` // crypto→fiat
    if (reversed.value) return `${t}-${g}` // fiat→crypto (queried inverted)
    return '' // crypto↔crypto / fiat↔fiat: unsupported
  })

  const { data: ratesData, state } = useRatesQuery(directionSlug)
  useRatesStream(directionSlug)

  const rates = computed<RateRow[]>(() => ratesData.value?.rates ?? [])

  // Every exchanger's rate as get-per-give, plus the resulting receive amount,
  // sorted best (largest receive) first. This ordering is correct for both
  // directions: for fiat→crypto the inverted rate flips which exchanger wins.
  const ranked = computed<RankedRate[]>(() => {
    const invert = reversed.value
    const amt = amount.value
    return rates.value
      .map((row): RankedRate => {
        const raw = parseFloat(row.rate)
        const effectiveRate = Number.isFinite(raw) && raw > 0 ? (invert ? 1 / raw : raw) : 0
        return { row, effectiveRate, receive: amt * effectiveRate }
      })
      .filter((r) => r.effectiveRate > 0)
      .sort((a, b) => b.effectiveRate - a.effectiveRate)
  })

  const best = computed<RankedRate | null>(() => ranked.value[0] ?? null)
  const bestRow = computed<RateRow | null>(() => best.value?.row ?? null)
  const unitRate = computed<number | null>(() => best.value?.effectiveRate ?? null)
  const receiveAmount = computed<number | null>(() => best.value?.receive ?? null)
  const exchangerCount = computed(() => ranked.value.length)
  const hasRates = computed(() => ranked.value.length > 0)

  const giveCode = computed(() => give.value.code)
  const getCode = computed(() => get.value.code)
  const giveUnit = computed(() => currencyUnit(give.value.code))
  const receiveUnit = computed(() => currencyUnit(get.value.code))

  /** Swap give/get, carrying the current receive amount over to the new give. */
  function swap() {
    const carried = receiveAmount.value
    swapCurrencies()
    if (carried != null && Number.isFinite(carried)) {
      const factor = 10 ** currencyDecimals(give.value.code) // give is now the old get
      amount.value = Math.round(carried * factor) / factor
    }
  }

  return {
    give,
    get,
    amount,
    directionSlug,
    reversed,
    hasRates,
    state,
    rates,
    ranked,
    best,
    bestRow,
    unitRate,
    receiveAmount,
    exchangerCount,
    giveCode,
    getCode,
    giveUnit,
    receiveUnit,
    swap,
  }
}
