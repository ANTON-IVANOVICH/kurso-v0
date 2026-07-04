import type { components } from '~/types/api'

// Contract types generated from api/openapi.yaml (`pnpm openapi`).
export type Currency = components['schemas']['Currency']
export type Direction = components['schemas']['Direction']
export type Exchanger = components['schemas']['Exchanger']
export type RateRow = components['schemas']['RateRow']
export type RatesResponse = components['schemas']['RatesResponse']

/** Base URL of kurso-api, from runtime config (NUXT_PUBLIC_API_BASE). */
export function useApiBase(): string {
  return useRuntimeConfig().public.apiBase as string
}
