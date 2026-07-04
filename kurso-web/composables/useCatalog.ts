import { useQuery } from '@pinia/colada'
import { toValue, type MaybeRefOrGetter } from 'vue'
import type { Currency, Direction, Exchanger } from './useApi'

// Catalogue reads via Pina Colada. Nuxt handles the SSR fetch + hydration; the
// queries are cached and shared across components by their key. Each resolves to
// an empty list on failure so a down API degrades gracefully rather than
// crashing SSR.

async function safeList<T>(base: string, path: string): Promise<T[]> {
  try {
    return await $fetch<T[]>(path, { baseURL: base })
  } catch {
    return []
  }
}

export function useCurrenciesQuery() {
  const base = useApiBase()
  return useQuery({
    key: ['currencies'],
    query: () => safeList<Currency>(base, '/api/v1/currencies'),
  })
}

export function useDirectionsQuery() {
  const base = useApiBase()
  return useQuery({
    key: ['directions'],
    query: () => safeList<Direction>(base, '/api/v1/directions'),
  })
}

export function useExchangersQuery() {
  const base = useApiBase()
  return useQuery({
    key: ['exchangers'],
    query: () => safeList<Exchanger>(base, '/api/v1/exchangers'),
  })
}

export function useExchangerQuery(slug: MaybeRefOrGetter<string>) {
  const base = useApiBase()
  return useQuery({
    key: () => ['exchanger', toValue(slug)],
    // null on 404/failure so the page can render a graceful "not found" state.
    query: async (): Promise<Exchanger | null> => {
      try {
        return await $fetch<Exchanger>(`/api/v1/exchangers/${toValue(slug)}`, { baseURL: base })
      } catch {
        return null
      }
    },
    enabled: () => !!toValue(slug),
  })
}
