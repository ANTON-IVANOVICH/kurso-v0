import { computed, toValue, type MaybeRefOrGetter } from 'vue'
import type { AsyncData } from 'nuxt/app'
import type { Currency, Direction, Exchanger } from './useApi'

// Catalogue reads via Nuxt's useAsyncData — server-rendered (SSR/ISR) and cached
// by key. Each degrades to an empty/null result on failure so a down API doesn't
// crash SSR. A compat `state` (Colada-shaped `{ status, data, error }`) is
// exposed so call sites read `state.value.status` unchanged.

async function safeList<T>(base: string, path: string): Promise<T[]> {
  try {
    return await $fetch<T[]>(path, { baseURL: base })
  } catch {
    return []
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function withState<T, R extends AsyncData<T, any>>(res: R) {
  const state = computed(() => ({
    status: res.status.value === 'idle' ? 'pending' : res.status.value,
    data: res.data.value,
    error: res.error.value,
  }))
  return Object.assign(res, { state })
}

export function useCurrenciesQuery() {
  const base = useApiBase()
  return withState(
    useAsyncData<Currency[]>('currencies', () => safeList<Currency>(base, '/api/v1/currencies'), {
      default: () => [],
    }),
  )
}

export function useDirectionsQuery() {
  const base = useApiBase()
  return withState(
    useAsyncData<Direction[]>('directions', () => safeList<Direction>(base, '/api/v1/directions'), {
      default: () => [],
    }),
  )
}

export function useExchangersQuery() {
  const base = useApiBase()
  return withState(
    useAsyncData<Exchanger[]>('exchangers', () => safeList<Exchanger>(base, '/api/v1/exchangers'), {
      default: () => [],
    }),
  )
}

export function useExchangerQuery(slug: MaybeRefOrGetter<string>) {
  const base = useApiBase()
  return withState(
    useAsyncData<Exchanger | null>(
      `exchanger-${toValue(slug)}`,
      () => {
        const s = toValue(slug)
        if (!s) return Promise.resolve(null)
        // null on 404/failure so the page can render a graceful "not found" state.
        return $fetch<Exchanger>(`/api/v1/exchangers/${s}`, { baseURL: base }).catch(() => null)
      },
      { watch: [() => toValue(slug)], default: () => null },
    ),
  )
}
