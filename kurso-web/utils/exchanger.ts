// Deterministic avatar (initials + colour) for an exchanger. Known seeded slugs
// keep their brand-ish colours; anything else hashes to a stable palette entry.
const KNOWN: Record<string, { initials: string; color: string }> = {
  cryptobridge: { initials: 'CB', color: '#3A4452' },
  netex24: { initials: 'N', color: '#5B3FA0' },
  '24paybank': { initials: '24', color: '#1F8A5B' },
  coino: { initials: 'Co', color: '#26A17B' },
  baksman: { initials: 'BM', color: '#8A5A2B' },
  bitx: { initials: 'BX', color: '#2E5C8A' },
}
const PALETTE = [
  '#3A4452',
  '#5B3FA0',
  '#1F8A5B',
  '#8A5A2B',
  '#26A17B',
  '#2E5C8A',
  '#7A3F8A',
  '#A0463F',
]

export function exchangerAvatar(slug: string, name: string): { initials: string; color: string } {
  const known = KNOWN[slug]
  if (known) return known
  let hash = 0
  for (const ch of slug) hash = (hash * 31 + ch.charCodeAt(0)) >>> 0
  const letters = name.replace(/[^\p{L}\p{N}]/gu, '')
  return {
    initials: (letters.slice(0, 2) || name.slice(0, 2)).toUpperCase(),
    color: PALETTE[hash % PALETTE.length],
  }
}

/** Currency badge glyph + colour for direction pills. */
const CURRENCY: Record<string, { symbol: string; color: string; dark?: boolean }> = {
  USDT: { symbol: '₮', color: '#26A17B' },
  BTC: { symbol: '₿', color: '#F7931A' },
  ETH: { symbol: 'Ξ', color: '#627EEA' },
  TON: { symbol: 'TON', color: '#0098EA' },
  TINKOFF: { symbol: 'Т', color: '#FFDD2D', dark: true },
  SBER: { symbol: 'С', color: '#1FAE54' },
  RUBCASH: { symbol: '₽', color: '#3A4452' },
}
export function currencyBadge(code: string): { symbol: string; color: string; dark?: boolean } {
  return CURRENCY[code] ?? { symbol: code.slice(0, 1), color: '#3A4452' }
}
