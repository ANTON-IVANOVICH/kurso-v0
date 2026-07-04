<script setup lang="ts">
type State = 'best' | 'normal' | 'stale'
type DeltaTone = 'success' | 'danger' | 'muted'

interface Exchanger {
  name: string
  initials: string
  avatarColor: string
  rating: string
  reviews: string
  time: string
  reserve: string
  extraTag?: string
  partner?: boolean
  note?: string
  state: State
  receive: string
  rate: string
  delta: string
  deltaTone: DeltaTone
  href?: string
}

// Static fallback shown during first render and when the API is unreachable.
const fallbackExchangers: Exchanger[] = [
  {
    name: 'CryptoBridge',
    initials: 'CB',
    avatarColor: '#3A4452',
    rating: '4.9',
    reviews: '1203',
    time: '10 мин',
    reserve: '24.8M ₽',
    extraTag: 'AML ✓',
    partner: true,
    state: 'best',
    receive: '81 200 ₽',
    rate: '81.20 ₽',
    delta: 'Лучший курс',
    deltaTone: 'success',
  },
  {
    name: 'NetEx24',
    initials: 'N',
    avatarColor: '#5B3FA0',
    rating: '4.9',
    reviews: '2104',
    time: '15 мин',
    reserve: '31.2M ₽',
    extraTag: 'без KYC',
    state: 'normal',
    receive: '80 950 ₽',
    rate: '80.95 ₽',
    delta: '−250 ₽',
    deltaTone: 'danger',
  },
  {
    name: '24Paybank',
    initials: '24',
    avatarColor: '#1F8A5B',
    rating: '4.7',
    reviews: '560',
    time: '20 мин',
    reserve: '8.4M ₽',
    state: 'normal',
    receive: '80 740 ₽',
    rate: '80.74 ₽',
    delta: '−460 ₽',
    deltaTone: 'danger',
  },
  {
    name: 'BaksMan',
    initials: 'BM',
    avatarColor: '#8A5A2B',
    rating: '4.4',
    reviews: '310',
    time: '30 мин',
    reserve: '2.1M ₽',
    note: 'курс устарел',
    state: 'stale',
    receive: '79 100 ₽',
    rate: '79.10 ₽',
    delta: '−2 100 ₽',
    deltaTone: 'muted',
  },
]

// --- live rates (shared with the calculator via useCalculator), with fallback ---
// `ranked` is the current direction's exchangers, best-first, already converted
// at the entered amount — so the top card always matches the calculator's best.
const { ranked, getCode, receiveUnit, directionSlug } = useCalculator()
const { filter, sort } = useHomeFilters()
const { t } = useI18n()
const apiBase = useApiBase()

// Apply the FilterBar's filter + sort to the ranked list (real data only).
const filteredRanked = computed(() => {
  let list = ranked.value
  if (filter.value === 'partners') list = list.filter((r) => r.row.partner)
  else if (filter.value === 'reserve5m')
    list = list.filter((r) => parseFloat(r.row.reserve ?? '0') >= RESERVE_THRESHOLD)
  if (sort.value === 'rating')
    list = [...list].sort((a, b) => (b.row.ratingAvg ?? 0) - (a.row.ratingAvg ?? 0))
  return list
})

const slugMeta: Record<string, { initials: string; color: string }> = {
  cryptobridge: { initials: 'CB', color: '#3A4452' },
  netex24: { initials: 'N', color: '#5B3FA0' },
  '24paybank': { initials: '24', color: '#1F8A5B' },
  baksman: { initials: 'BM', color: '#8A5A2B' },
  coino: { initials: 'Co', color: '#26A17B' },
  bitx: { initials: 'BX', color: '#2E5C8A' },
}
const palette = ['#3A4452', '#5B3FA0', '#1F8A5B', '#8A5A2B', '#26A17B', '#2E5C8A']

function fmtReserve(s: string | null | undefined): string {
  const n = s ? parseFloat(s) : 0
  if (n >= 1e6) return (n / 1e6).toFixed(1).replace(/\.0$/, '') + 'M ₽'
  if (n >= 1e3) return Math.round(n / 1e3) + 'K ₽'
  return fmtNumber(n, 0) + ' ₽'
}

const liveExchangers = computed<Exchanger[] | null>(() => {
  // null → no live data at all (fall back to the static demo list).
  if (ranked.value.length === 0) return null
  const list = filteredRanked.value
  const unit = receiveUnit.value
  const dec = currencyDecimals(getCode.value)
  // "Best" (Лучший курс) = highest receive in the shown set, independent of sort.
  const maxReceive = list.length ? Math.max(...list.map((r) => r.receive)) : 0
  return list.map(({ row, effectiveRate, receive }, i): Exchanger => {
    const meta = slugMeta[row.exchangerSlug]
    const isBest = receive === maxReceive
    return {
      name: row.exchangerName,
      initials: meta?.initials ?? row.exchangerName.slice(0, 2).toUpperCase(),
      avatarColor: meta?.color ?? palette[i % palette.length],
      rating: row.ratingAvg != null ? row.ratingAvg.toFixed(1) : '—',
      reviews: String(row.reviewsCount),
      time: t('home.card.eta'),
      reserve: fmtReserve(row.reserve),
      partner: row.partner || undefined,
      state: isBest ? 'best' : 'normal',
      receive: `${fmtNumber(receive, dec)} ${unit}`,
      rate: `${fmtNumber(effectiveRate, rateDecimals(effectiveRate))} ${unit}`,
      delta: isBest ? t('home.card.bestRate') : `−${fmtNumber(maxReceive - receive, dec)} ${unit}`,
      deltaTone: isBest ? 'success' : 'danger',
      href: `${apiBase}/go/${row.exchangerSlug}${directionSlug.value ? `?direction=${directionSlug.value}` : ''}`,
    }
  })
})

// null → fallback; [] → live data exists but nothing matches the active filter.
const displayExchangers = computed(() => liveExchangers.value ?? fallbackExchangers)

useSeoMeta({
  title: 'Kurso — курсы криптообменников',
  description: 'Сравните курсы обмена USDT, BTC и других валют у проверенных обменников.',
})
</script>

<template>
  <div class="mx-auto max-w-[1200px] px-4 py-6 md:px-6">
    <div class="grid grid-cols-1 gap-6 md:grid-cols-[340px_1fr] md:items-start">
      <aside class="md:sticky md:top-6">
        <ExchangeCalculator />
      </aside>

      <div class="min-w-0">
        <FilterBar />
        <div class="flex flex-col gap-3">
          <ExchangerCard v-for="ex in displayExchangers" :key="ex.name" v-bind="ex" />
          <div
            v-if="displayExchangers.length === 0"
            class="rounded-2xl border border-line bg-surface py-16 text-center text-sm text-ink-faint"
          >
            {{ t('home.nothingFound') }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
