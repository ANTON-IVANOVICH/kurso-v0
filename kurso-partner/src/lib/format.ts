// Display formatting for the merchant cabinet. Numbers use a thin space as the
// thousands separator to line up under the tabular monospace face (`.tnum`).

const THIN = ' '

/** Integer with thin-space grouping: 84200 → "84 200". */
export function fmtInt(n: number): string {
  return Math.round(n)
    .toString()
    .replace(/\B(?=(\d{3})+(?!\d))/g, THIN)
}

/** Compact figure for metrics: 84200 → "84.2k", 2_400_000 → "2.4M". */
export function fmtCompact(n: number): string {
  const abs = Math.abs(n)
  if (abs >= 1_000_000) return trim1(n / 1_000_000) + 'M'
  if (abs >= 1_000) return trim1(n / 1_000) + 'k'
  return String(Math.round(n))
}

/** Ruble amount, compact: 412000 → "₽412k". */
export function fmtRub(n: number): string {
  return '₽' + fmtCompact(n)
}

/** Parse a decimal string (reserves come from the API as strings) to a number. */
export function toNum(s: string | null | undefined): number {
  if (!s) return 0
  const n = Number(s)
  return Number.isFinite(n) ? n : 0
}

/** Russian 3-form plural: pluralRu(n, 'обменник', 'обменника', 'обменников'). */
export function pluralRu(n: number, one: string, few: string, many: string): string {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod10 === 1 && mod100 !== 11) return one
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 10 || mod100 >= 20)) return few
  return many
}

function trim1(x: number): string {
  const r = Math.round(x * 10) / 10
  return Number.isInteger(r) ? String(r) : r.toFixed(1)
}

/** Relative "N сек/мин/ч/дн назад" from an ISO timestamp. */
export function timeAgo(iso: string | null | undefined): string {
  if (!iso) return '—'
  const then = new Date(iso).getTime()
  if (!Number.isFinite(then)) return '—'
  const sec = Math.max(0, Math.round((Date.now() - then) / 1000))
  if (sec < 60) return `${sec} сек назад`
  const min = Math.round(sec / 60)
  if (min < 60) return `${min} ${pluralRu(min, 'мин', 'мин', 'мин')} назад`
  const hr = Math.round(min / 60)
  if (hr < 24) return `${hr} ${pluralRu(hr, 'ч', 'ч', 'ч')} назад`
  const day = Math.round(hr / 24)
  return `${day} ${pluralRu(day, 'день', 'дня', 'дней')} назад`
}

/** Money in rubles with thin-space grouping: 14800 → "14 800 ₽". */
export function fmtRubFull(n: number): string {
  return fmtInt(n) + ' ₽'
}

/** Parse a decimal rate string to a grouped display, keeping up to 2 decimals. */
export function fmtRate(s: string | null | undefined): string {
  if (!s) return '— ₽'
  const n = Number(s)
  if (!Number.isFinite(n)) return '— ₽'
  const rounded = n >= 1000 ? Math.round(n) : Math.round(n * 100) / 100
  const [int, dec] = String(rounded).split('.')
  const grouped = int.replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  return (dec ? `${grouped}.${dec}` : grouped) + ' ₽'
}
