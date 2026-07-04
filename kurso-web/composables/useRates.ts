import { computed, onScopeDispose, toValue, watch, type MaybeRefOrGetter } from 'vue'
import type { RateRow, RatesResponse } from './useApi'

// Live rates for a direction. The initial (SSR) load comes through Nuxt's
// useAsyncData; live updates arrive via {@link useRatesStream} which patches a
// shared per-slug useState store. Reading rates from that store (rather than
// useAsyncData's own data) keeps every consumer of a slug in sync — including
// the calculator whose direction changes over time — and keeps SSE the only
// piece that mutates rates, no Pinia Colada cache involved.

type RatesStore = Record<string, RatesResponse | null>

const STORE_KEY = 'rates'

export function useRatesQuery(direction: MaybeRefOrGetter<string>) {
  const base = useApiBase()
  const slug = computed(() => toValue(direction))
  const store = useState<RatesStore>(STORE_KEY, () => ({}))

  const res = useAsyncData<RatesResponse | null>(
    `rates-${slug.value}`,
    async () => {
      const s = slug.value
      if (!s) return null
      // null on failure so a down/unreachable API degrades to the static fallback.
      const data = await $fetch<RatesResponse>(`/api/v1/rates/${s}`, { baseURL: base }).catch(
        () => null,
      )
      store.value = { ...store.value, [s]: data }
      return data
    },
    { watch: [slug], default: () => null },
  )

  const data = computed(() => (slug.value ? (store.value[slug.value] ?? null) : null))
  const state = computed(() => ({
    status: res.status.value === 'idle' ? 'pending' : res.status.value,
    data: data.value,
    error: res.error.value,
  }))

  return { data, state, pending: res.pending, error: res.error, refresh: res.refresh }
}

// One EventSource per direction slug, shared across every subscriber (calculator
// + card list watch the same direction). Ref-counted so the connection opens
// once and closes when the last consumer leaves.
interface StreamEntry {
  es: EventSource
  refs: number
}
const streams = new Map<string, StreamEntry>()

/**
 * Subscribes to the server's SSE rate stream and patches the shared rates store
 * for the matching slug, so the UI updates live without polling. Client-only.
 */
export function useRatesStream(direction: MaybeRefOrGetter<string>) {
  if (!import.meta.client) return

  const base = useApiBase()
  const store = useState<RatesStore>(STORE_KEY, () => ({}))
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
        const prev = store.value[slug]
        if (prev) store.value = { ...store.value, [slug]: { ...prev, rates: rows } }
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
