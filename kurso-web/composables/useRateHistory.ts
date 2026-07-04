import { computed, toValue, type MaybeRefOrGetter } from 'vue'

// Best-rate history for a direction, from the public sparkline endpoint. Used by
// the alert builder to draw the 7-day chart and compute the "forecast" from real
// observed rates (how often the threshold was reached). SSR-friendly via
// useAsyncData; refetches when the direction or window changes.

export interface RateHistoryPoint {
  t: string
  rate: string
}
interface RateHistoryResponse {
  days: number
  points: RateHistoryPoint[]
}

export function useRateHistory(
  direction: MaybeRefOrGetter<string>,
  days: MaybeRefOrGetter<number> = 7,
) {
  const base = useApiBase()
  const slug = computed(() => toValue(direction))
  const window = computed(() => toValue(days))

  const res = useAsyncData<RateHistoryResponse | null>(
    `rate-history-${slug.value}-${window.value}`,
    async () => {
      const s = slug.value
      if (!s) return null
      return await $fetch<RateHistoryResponse>(`/api/v1/rates/${s}/history`, {
        baseURL: base,
        query: { days: window.value },
      }).catch(() => null)
    },
    { watch: [slug, window], default: () => null },
  )

  // Numeric series (chronological). Not inverted — the builder inverts for a
  // fiat→crypto selection to match the displayed unit.
  const series = computed<number[]>(() =>
    (res.data.value?.points ?? [])
      .map((p) => parseFloat(p.rate))
      .filter((n) => Number.isFinite(n) && n > 0),
  )

  return { series, points: computed(() => res.data.value?.points ?? []), refresh: res.refresh }
}
