// Hand-written mirrors of the kurso-api merchant-cabinet JSON contracts
// (see kurso-api internal/adapter/http/merchant.go + partner_auth.go).

export interface Merchant {
  id: string
  email: string
  role: 'owner' | 'manager' | 'viewer'
  exchangerSlug: string
  exchangerName: string
}

export interface MerchantMetrics {
  ratesActive: number
  ratesTotal: number
  ratesStale: number
  clicksToday: number
  clicksYesterday: number
  ratingAvg: number
  reviewsCount: number
  reviewsUnanswered: number
}

export interface TrafficPoint {
  day: string
  clicks: number
}

export interface Dashboard {
  metrics: MerchantMetrics
  traffic: TrafficPoint[]
}

export type FeedStatus = 'ok' | 'delayed' | 'down'

export interface MerchantRate {
  directionId: string
  directionSlug: string
  fromCode: string
  toCode: string
  rate: string | null
  reserve: string | null
  fetchedAt: string | null
  feed: FeedStatus
}

export interface MerchantReview {
  id: string
  author: string
  rating: number
  title?: string | null
  body: string
  createdAt: string
  reply?: string | null
  replyAt?: string | null
}

export interface TrafficDirection {
  directionSlug: string
  fromCode: string
  toCode: string
  clicks: number
}

export interface Traffic {
  days: number
  total: number
  series: TrafficPoint[]
  directions: TrafficDirection[]
}

export interface MerchantProfile {
  slug: string
  name: string
  description: string | null
  websiteUrl: string | null
  logoUrl: string | null
  status: 'active' | 'paused' | 'banned'
  isVerified: boolean
  ratingAvg: string | null
  reviewsCount: number
  onSince: number
  assets: string[]
  reserveTotal: string | null
}

export interface MerchantComplaint {
  id: string
  reviewId: string
  reason: string
  details?: string | null
  status: 'open' | 'reviewed' | 'dismissed'
  createdAt: string
  reviewBody: string
  author: string
  rating: number
}

export interface Payout {
  id: string
  periodStart: string
  periodEnd: string
  clicksCount: number
  amount: string
  currency: string
  status: 'pending' | 'paid' | 'cancelled'
  paidAt?: string | null
}

export interface Billing {
  perClick: number
  currentMonth: { clicks: number; estimated: number }
  payouts: Payout[]
}
