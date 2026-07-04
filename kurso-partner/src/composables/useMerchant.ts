import { defineQuery, useQuery, useQueryCache } from '@pinia/colada'
import { api } from '../lib/api'
import type {
  Billing,
  Dashboard,
  MerchantComplaint,
  MerchantProfile,
  MerchantRate,
  MerchantReview,
  Traffic,
} from '../types/models'

// Cabinet reads via Pinia Colada — the SPA's async cache. Queries are deduped by
// key and shared across components (sidebar badges + pages read one fetch).
// `defineQuery` makes each a singleton composable; mutations invalidate the key.

export const useDashboardQuery = defineQuery(() =>
  useQuery({ key: ['dashboard'], query: () => api.get<Dashboard>('/partner/dashboard') }),
)

export const useRatesQuery = defineQuery(() =>
  useQuery({ key: ['rates'], query: () => api.get<MerchantRate[]>('/partner/rates') }),
)

export const useReviewsQuery = defineQuery(() =>
  useQuery({ key: ['reviews'], query: () => api.get<MerchantReview[]>('/partner/reviews') }),
)

export const useProfileQuery = defineQuery(() =>
  useQuery({ key: ['profile'], query: () => api.get<MerchantProfile>('/partner/profile') }),
)

export const useComplaintsQuery = defineQuery(() =>
  useQuery({
    key: ['complaints'],
    query: () => api.get<MerchantComplaint[]>('/partner/complaints'),
  }),
)

export const useBillingQuery = defineQuery(() =>
  useQuery({ key: ['billing'], query: () => api.get<Billing>('/partner/billing') }),
)

/** Traffic query keyed by the selected window so switching tabs refetches. */
export function useTrafficQuery(days: () => number) {
  return useQuery({
    key: () => ['traffic', days()],
    query: () => api.get<Traffic>(`/partner/traffic?days=${days()}`),
  })
}

/** Invalidate helper so mutation handlers can refresh affected caches. */
export function useCabinetCache() {
  const cache = useQueryCache()
  return {
    invalidate(key: string[]) {
      void cache.invalidateQueries({ key })
    },
  }
}
