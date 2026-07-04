// Shared filter/sort state for the home exchanger list. FilterBar writes it and
// the home page reads it, so the chips actually filter the results. Only filters
// backed by real data are offered (partner flag, reserve threshold, sort).
export type HomeFilter = 'all' | 'partners' | 'reserve5m'
export type HomeSort = 'best' | 'rating'

export const RESERVE_THRESHOLD = 5_000_000

export function useHomeFilters() {
  const filter = useState<HomeFilter>('homeFilter', () => 'all')
  const sort = useState<HomeSort>('homeSort', () => 'best')
  return { filter, sort }
}
