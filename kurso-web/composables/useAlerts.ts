import { computed } from 'vue'

// User price alerts. Real client state: the builder creates them, they persist
// (localStorage via the persist plugin) and survive reloads. Each alert stores
// its direction so the account/alert screens compute the CURRENT rate and
// triggered status live from the rates API — there is no stored/fake price.
// Server-side matching + Telegram delivery arrive with the account backend.

export type AlertChannel = 'telegram' | 'email' | 'push'

export interface AlertBadge {
  symbol: string
  color: string
  dark?: boolean
}

export interface Alert {
  id: string
  directionSlug: string // seeded rate direction, e.g. "usdt-tinkoff"
  reversed: boolean // fiat→crypto selection (rate inverted)
  pair: string // display, e.g. "USDT → Тинькофф"
  badge: AlertBadge // give-side coin/bank badge
  unit: string // rate unit, e.g. "₽"
  direction: 'above' | 'below'
  threshold: number
  channels: AlertChannel[]
  paused: boolean
  createdAt: number
  // Builder preferences (persisted; applied once the server-side matcher lands).
  scope?: 'all' | 'favorites' | 'manual'
  minReserve?: boolean
  ratingFilter?: boolean
  validity?: '7' | '30' | 'forever'
}

/** Free-tier cap on active alerts. */
export const ALERT_LIMIT = 10

/** Given a live current rate, resolve the alert's runtime status. */
export function alertStatus(a: Alert, current: number | null): 'active' | 'triggered' | 'paused' {
  if (a.paused) return 'paused'
  if (current == null) return 'active'
  const hit = a.direction === 'above' ? current >= a.threshold : current <= a.threshold
  return hit ? 'triggered' : 'active'
}

/** Progress (0..1) of current toward threshold, for the progress bars. */
export function alertProgress(a: Alert, current: number | null): number {
  if (current == null) return 0
  if (a.direction === 'above') {
    if (current >= a.threshold) return 1
    return Math.max(0, Math.min(1, current / a.threshold))
  }
  if (current <= a.threshold) return 1
  // below: closer as current falls toward threshold (cap the window at +25%)
  return Math.max(0, Math.min(1, a.threshold / current))
}

export function useAlerts() {
  const alerts = useState<Alert[]>('alerts', () => [])

  const activeCount = computed(() => alerts.value.filter((a) => !a.paused).length)

  function add(a: Omit<Alert, 'id' | 'createdAt' | 'paused'>): string {
    const now = Date.now()
    const id = `alr_${now.toString(36)}`
    alerts.value = [{ ...a, id, paused: false, createdAt: now }, ...alerts.value]
    return id
  }
  function pause(id: string) {
    const a = alerts.value.find((x) => x.id === id)
    if (a) a.paused = true
  }
  function resume(id: string) {
    const a = alerts.value.find((x) => x.id === id)
    if (a) a.paused = false
  }
  function remove(id: string) {
    alerts.value = alerts.value.filter((x) => x.id !== id)
  }
  function find(id: string) {
    return alerts.value.find((x) => x.id === id)
  }

  return { alerts, activeCount, add, pause, resume, remove, find }
}
