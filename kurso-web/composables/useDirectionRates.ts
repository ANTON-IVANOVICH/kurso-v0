import { computed, toValue, type MaybeRefOrGetter } from 'vue'
import type { RateRow } from './useApi'

// Live best rate for a direction slug, reusing the same Pinia-Colada query and
// SSE stream as the home calculator (deduped by key / ref-counted connection).
// `reversed` inverts the rate for a fiat→crypto selection, exactly like
// useCalculator. Powers alert "current rate" + the triggered screen's list.

export interface RankedRow {
  row: RateRow
  rate: number // effective get-per-give rate
}

export function useDirectionRates(
  slug: MaybeRefOrGetter<string>,
  reversed?: MaybeRefOrGetter<boolean>,
) {
  const { data, state } = useRatesQuery(slug)
  useRatesStream(slug)

  const ranked = computed<RankedRow[]>(() => {
    const invert = toValue(reversed) ?? false
    return (data.value?.rates ?? [])
      .map((row): RankedRow => {
        const raw = parseFloat(row.rate)
        return { row, rate: raw > 0 ? (invert ? 1 / raw : raw) : 0 }
      })
      .filter((r) => r.rate > 0)
      .sort((a, b) => b.rate - a.rate)
  })

  const best = computed<number | null>(() => ranked.value[0]?.rate ?? null)
  const bestRow = computed<RateRow | null>(() => ranked.value[0]?.row ?? null)
  const direction = computed(() => data.value?.direction ?? null)
  const loading = computed(
    () => !!toValue(slug) && data.value === null && state.value.status !== 'error',
  )

  return { ranked, best, bestRow, direction, loading }
}
