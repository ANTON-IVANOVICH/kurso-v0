// Currency-aware number formatting shared by the calculator and the exchanger
// cards. All helpers are pure and deterministic so they render identically on
// the server and the client (no locale/`toLocaleString` hydration drift).

const CRYPTO_CODES = new Set(['USDT', 'BTC', 'ETH', 'TON'])

/** True for the fiat "bank" options (Тинькофф, Сбербанк) — everything non-crypto. */
export function isFiatCode(code: string): boolean {
  return !CRYPTO_CODES.has(code)
}

/** Sensible number of decimals to display an amount in the given currency. */
export function currencyDecimals(code: string): number {
  if (isFiatCode(code)) return 0
  if (code === 'BTC') return 6
  if (code === 'ETH') return 4
  return 2
}

/** Unit suffix shown next to an amount: ₽ for banks, the ticker for crypto. */
export function currencyUnit(code: string): string {
  return isFiatCode(code) ? '₽' : code
}

/** Decimals for a per-unit rate — more precision for sub-1 (inverted) rates. */
export function rateDecimals(rate: number): number {
  if (rate >= 1000) return 0
  if (rate >= 1) return 2
  return 4
}

/** Group thousands with a thin space and fix decimals. `−` for negatives. */
export function fmtNumber(n: number, decimals: number): string {
  if (!Number.isFinite(n)) return '0'
  const [int, frac] = Math.abs(n).toFixed(decimals).split('.')
  const grouped = int.replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  const body = frac ? `${grouped}.${frac}` : grouped
  return n < 0 ? `−${body}` : body
}

/** Compact magnitude: 24.8M, 312K, 980 (no currency sign). */
export function fmtCompact(n: number): string {
  if (!Number.isFinite(n)) return '0'
  if (n >= 1e6) return (n / 1e6).toFixed(1).replace(/\.0$/, '') + 'M'
  if (n >= 1e3) return Math.round(n / 1e3) + 'K'
  return fmtNumber(n, 0)
}

/** Russian plural selector: pluralRu(n, 'обменник', 'обменника', 'обменников'). */
export function pluralRu(n: number, one: string, few: string, many: string): string {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod10 === 1 && mod100 !== 11) return one
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return few
  return many
}
