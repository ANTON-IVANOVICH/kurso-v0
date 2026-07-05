// Clickout history — the real exchanger jumps the user made. `record()` is called
// from every "Перейти" button right before the outbound /go redirect, and the
// entry persists (localStorage). The account overview/history screens read it.
// Starts empty; fills as the user actually uses the site.

export interface HistoryEntry {
  id: string
  pair: string
  exchanger: string
  slug: string
  directionSlug: string
  amount: string // formatted receive at click time, e.g. "81.20 ₽"
  at: number // epoch ms
}

const MAX = 50

export function useHistory() {
  const history = useState<HistoryEntry[]>('history', () => [])

  function record(entry: Omit<HistoryEntry, 'id' | 'at'>) {
    const at = Date.now()
    history.value = [{ ...entry, id: `h_${at.toString(36)}`, at }, ...history.value].slice(0, MAX)
  }

  return { history, record }
}

/** Coarse "N ago" label for a past timestamp (client-side). */
export function timeAgo(at: number, now = Date.now()): string {
  const { t } = useI18n()
  const s = Math.max(0, Math.round((now - at) / 1000))
  if (s < 60) return t('timeAgo.justNow')
  const m = Math.floor(s / 60)
  if (m < 60) return t('timeAgo.minutes', { n: m })
  const h = Math.floor(m / 60)
  if (h < 24) return t('timeAgo.hours', { n: h })
  const d = Math.floor(h / 24)
  if (d < 7) return t('timeAgo.days', { n: d })
  const w = Math.floor(d / 7)
  if (w < 5) return t('timeAgo.weeks', { n: w })
  const mo = Math.floor(d / 30)
  return t('timeAgo.months', { n: mo })
}
