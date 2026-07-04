import { toValue, type MaybeRefOrGetter } from 'vue'

// Reviews for an exchanger. READS use Nuxt's useAsyncData so they render on the
// server (SSR/ISR, good for SEO); only the write ACTIONS below are plain $fetch
// mutations. (Pinia Colada is reserved for the admin SPA / client actions.)

export interface ReviewItem {
  id: string
  author: string
  rating: number
  title?: string
  body: string
  status: string
  createdAt: string
}
export interface RatingSummary {
  average: number
  count: number
  histogram: number[] // [1★, 2★, 3★, 4★, 5★]
}
export interface ReviewsResponse {
  summary: RatingSummary
  reviews: ReviewItem[]
}

const empty = (): ReviewsResponse => ({
  summary: { average: 0, count: 0, histogram: [0, 0, 0, 0, 0] },
  reviews: [],
})

export function useReviewsData(slug: MaybeRefOrGetter<string>) {
  const base = useApiBase()
  return useAsyncData<ReviewsResponse>(
    `reviews-${toValue(slug)}`,
    () =>
      $fetch<ReviewsResponse>(`/api/v1/exchangers/${toValue(slug)}/reviews`, {
        baseURL: base,
      }).catch(empty),
    { watch: [() => toValue(slug)], default: empty },
  )
}

export interface NewReviewInput {
  author: string
  rating: number
  body: string
  title?: string
}

/** Submit a review; returns the created review (status published|pending). */
export function submitReview(slug: string, input: NewReviewInput) {
  const base = useApiBase()
  return $fetch<ReviewItem>(`/api/v1/exchangers/${slug}/reviews`, {
    baseURL: base,
    method: 'POST',
    body: input,
  })
}

/** File a complaint against a review. */
export function reportReview(id: string, reason: string, details?: string) {
  const base = useApiBase()
  return $fetch(`/api/v1/reviews/${id}/report`, {
    baseURL: base,
    method: 'POST',
    body: { reason, details },
  })
}
