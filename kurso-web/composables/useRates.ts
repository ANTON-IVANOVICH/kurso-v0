import { useQuery, useQueryCache } from '@pinia/colada'
import { onScopeDispose, toValue, watch, type MaybeRefOrGetter } from 'vue'
import type { RateRow, RatesResponse } from './useApi'

/**
 * Current rates for a direction, best-first. The initial (and SSR) load comes
 * through Pina Colada; live updates arrive via {@link useRatesStream}.
 */
export function useRatesQuery(direction: MaybeRefOrGetter<string>) {
  const base = useApiBase()
  return useQuery({
    key: () => ['rates', toValue(direction)],
    // Resolve to null on failure so a down/unreachable API degrades to the
    // static fallback instead of crashing SSR.
    query: async (): Promise<RatesResponse | null> => {
      try {
        return await $fetch<RatesResponse>(`/api/v1/rates/${toValue(direction)}`, { baseURL: base })
      } catch {
        return null
      }
    },
    enabled: () => !!toValue(direction),
  })
}

// One EventSource per direction slug, shared across every subscriber (the
// calculator and the card list both watch the same direction). Ref-counted so
// the connection is opened once and closed only when the last consumer leaves.
interface StreamEntry {
  es: EventSource
  refs: number
}
const streams = new Map<string, StreamEntry>()

/**
 * Subscribes to the server's SSE rate stream and patches the Colada cache for
 * the matching `useRatesQuery`, so the UI updates live without polling.
 * Client-only; cleans up on scope dispose and when the direction changes.
 */
export function useRatesStream(direction: MaybeRefOrGetter<string>) {
  if (!import.meta.client) return

  const base = useApiBase()
  const cache = useQueryCache()
  let currentSlug = ''

  const release = () => {
    if (!currentSlug) return
    const entry = streams.get(currentSlug)
    if (entry && --entry.refs <= 0) {
      entry.es.close()
      streams.delete(currentSlug)
    }
    currentSlug = ''
  }

  const acquire = (slug: string) => {
    if (slug === currentSlug) return
    release()
    if (!slug) return

    let entry = streams.get(slug)
    if (!entry) {
      const es = new EventSource(`${base}/api/v1/rates/${slug}/stream`)
      es.addEventListener('rates', (ev) => {
        let rows: RateRow[]
        try {
          rows = JSON.parse((ev as MessageEvent).data)
        } catch {
          return
        }
        // Patch only once the query has populated (keeps the direction metadata).
        const prev = cache.getQueryData<RatesResponse>(['rates', slug])
        if (prev) cache.setQueryData(['rates', slug], { ...prev, rates: rows })
      })
      // On error the browser reconnects automatically; nothing to do.
      entry = { es, refs: 0 }
      streams.set(slug, entry)
    }
    entry.refs++
    currentSlug = slug
  }

  watch(() => toValue(direction), acquire, { immediate: true })
  onScopeDispose(release)
}
