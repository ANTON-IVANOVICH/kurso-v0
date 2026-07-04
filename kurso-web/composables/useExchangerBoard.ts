import { toValue, type MaybeRefOrGetter } from 'vue'
import type { Direction, RatesResponse, RateRow } from './useApi'

/** One direction the exchanger trades, with its live market position. */
export interface BoardEntry {
  direction: Direction
  rows: RateRow[] // full ranking for the direction, best-first
  myRow: RateRow
  myRate: number
  rank: number // 1-based position among all exchangers (by rate)
  total: number
  marketAvg: number
  spreadPct: number // (myRate − marketAvg) / marketAvg × 100
}

/**
 * Aggregates the exchanger's real market position across every direction it
 * trades, from the public rates endpoint (one request per direction, run in
 * parallel). Powers the analytics + directions sections without any fabricated
 * numbers. Uses `useAsyncData` so it renders on the server and hydrates.
 */
export function useExchangerBoard(slug: MaybeRefOrGetter<string>) {
  const base = useApiBase()
  return useAsyncData<BoardEntry[]>(
    `exchanger-board-${toValue(slug)}`,
    async () => {
      const s = toValue(slug)
      if (!s) return []
      const directions = await $fetch<Direction[]>('/api/v1/directions', { baseURL: base }).catch(
        () => [] as Direction[],
      )
      const entries = await Promise.all(
        directions.map(async (direction): Promise<BoardEntry | null> => {
          const resp = await $fetch<RatesResponse>(`/api/v1/rates/${direction.slug}`, {
            baseURL: base,
          }).catch(() => null)
          if (!resp?.rates?.length) return null
          const rows = resp.rates
          const idx = rows.findIndex((r) => r.exchangerSlug === s)
          if (idx === -1) return null
          const myRow = rows[idx]
          const myRate = parseFloat(myRow.rate)
          if (!Number.isFinite(myRate)) return null
          const nums = rows.map((r) => parseFloat(r.rate)).filter((n) => Number.isFinite(n))
          const marketAvg = nums.reduce((a, b) => a + b, 0) / (nums.length || 1)
          return {
            direction,
            rows,
            myRow,
            myRate,
            rank: idx + 1,
            total: rows.length,
            marketAvg,
            spreadPct: marketAvg ? ((myRate - marketAvg) / marketAvg) * 100 : 0,
          }
        }),
      )
      return entries.filter((e): e is BoardEntry => e !== null)
    },
    { watch: [() => toValue(slug)] },
  )
}
