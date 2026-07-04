import { computed, toValue, type MaybeRefOrGetter } from 'vue'

// Located exchangers (cash desks) for the map, with their rate for a direction.
// SSR-friendly via useAsyncData so the list renders even before Mapbox mounts.

export interface MapPoint {
  slug: string
  name: string
  lat: number
  lng: number
  address?: string | null
  city?: string | null
  hours?: string | null
  ratingAvg?: number | null
  reviewsCount: number
  partner: boolean
  rate?: string | null
}
interface MapResponse {
  direction: { slug: string; fromCode: string; toCode: string }
  points: MapPoint[]
}

export function useMapPoints(direction: MaybeRefOrGetter<string> = 'usdt-tinkoff') {
  const base = useApiBase()
  const slug = computed(() => toValue(direction))

  const res = useAsyncData<MapResponse | null>(
    `map-${slug.value}`,
    () =>
      $fetch<MapResponse>('/api/v1/map', { baseURL: base, query: { direction: slug.value } }).catch(
        () => null,
      ),
    { watch: [slug], default: () => null },
  )

  const points = computed<MapPoint[]>(() => res.data.value?.points ?? [])
  return {
    points,
    direction: computed(() => res.data.value?.direction ?? null),
    refresh: res.refresh,
  }
}
