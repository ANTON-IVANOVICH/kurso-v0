// Display formatting for the admin console. Numbers use a thin space as the
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
