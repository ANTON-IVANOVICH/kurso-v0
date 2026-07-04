import { computed, reactive } from 'vue'
import { defineQuery, useQuery, useQueryCache } from '@pinia/colada'
import { api, ApiError } from '../lib/api'
import type { Currency, Direction, Exchanger } from '../types/models'

// Catalogue reads via Pinia Colada — the same async-cache layer kurso-web uses.
// Queries are deduped by key and shared across every component, so the sidebar
// badge, dashboard and tables all read one cached fetch. `defineQuery` makes each
// a singleton composable.

export const useExchangersQuery = defineQuery(() =>
  useQuery({ key: ['exchangers'], query: () => api.get<Exchanger[]>('/api/v1/exchangers') }),
)
export const useCurrenciesQuery = defineQuery(() =>
  useQuery({ key: ['currencies'], query: () => api.get<Currency[]>('/api/v1/currencies') }),
)
export const useDirectionsQuery = defineQuery(() =>
  useQuery({ key: ['directions'], query: () => api.get<Direction[]>('/api/v1/directions') }),
)

/** Aggregated catalogue view with the shape the admin pages consume. */
export function useCatalog() {
  const ex = useExchangersQuery()
  const cu = useCurrenciesQuery()
  const di = useDirectionsQuery()
  const cache = useQueryCache()

  const exchangers = computed(() => ex.data.value ?? [])
  const currencies = computed(() => cu.data.value ?? [])
  const directions = computed(() => di.data.value ?? [])

  const loading = computed(() => ex.isLoading.value || cu.isLoading.value || di.isLoading.value)
  const error = computed(() => {
    const e = ex.error.value ?? cu.error.value ?? di.error.value
    if (!e) return null
    return e instanceof ApiError ? e.message : 'Не удалось загрузить каталог'
  })

  const partnersCount = computed(() => exchangers.value.filter((e) => e.partner).length)
  const activeCount = computed(() => exchangers.value.filter((e) => e.status === 'active').length)

  const findExchanger = (slug: string) => exchangers.value.find((e) => e.slug === slug)

  /** Local-only patch to an exchanger (updates the query cache) pending the
      admin write API. */
  function patchExchanger(slug: string, patch: Partial<Exchanger>) {
    const list = cache.getQueryData<Exchanger[]>(['exchangers'])
    if (list) {
      cache.setQueryData(
        ['exchangers'],
        list.map((e) => (e.slug === slug ? { ...e, ...patch } : e)),
      )
    }
  }

  function refresh() {
    void ex.refetch()
    void cu.refetch()
    void di.refetch()
  }

  // reactive() unwraps the refs so `catalog.exchangers` reads as the array in
  // both templates and script — a drop-in for the old Pinia store's ergonomics.
  return reactive({
    exchangers,
    currencies,
    directions,
    loading,
    error,
    partnersCount,
    activeCount,
    findExchanger,
    patchExchanger,
    refresh,
  })
}
