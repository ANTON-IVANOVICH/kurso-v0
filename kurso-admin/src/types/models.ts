// Hand-written mirrors of the kurso-api JSON contracts the admin actually reads.
// The generated `api.d.ts` only covers a couple of health probes; these match
// the real catalogue payloads (see kurso-api mappers.go).

export type ExchangerStatus = 'active' | 'paused' | 'banned'

export interface Exchanger {
  id: string
  slug: string
  name: string
  status: ExchangerStatus
  partner: boolean
  isVerified: boolean
  reviewsCount: number
  ratingAvg?: number | null
  websiteUrl?: string | null
  logoUrl?: string | null
  description?: string | null
  reserveTotal?: string | null
  directionsCount: number
  assets: string[]
  onSince: number
}

export type CurrencyKind = 'crypto' | 'fiat' | 'cash'

export interface Currency {
  id: string
  code: string
  name: string
  kind: CurrencyKind
  network?: string | null
  iconUrl?: string | null
}

export interface Direction {
  id: string
  slug: string
  fromCurrencyId: string
  toCurrencyId: string
  fromCode: string
  fromName: string
  toCode: string
  toName: string
  isPopular: boolean
}
