import { ref } from 'vue'

// Affiliate (partner) overview for the current user. The endpoint is bearer-
// protected, so this loads client-side via the authed fetch helper. Real stats:
// clicks come from clickouts attributed to the user's ref code, registrations
// from referred sign-ups.

export interface ReferralTag {
  tag: string
  clicks: number
}
export interface ReferralPoint {
  day: string
  clicks: number
}
export interface PartnerOverview {
  code: string
  revsharePct: number
  cookieDays: number
  clicks: number
  registrations: number
  estimatedRub: number
  series: ReferralPoint[]
  byTag: ReferralTag[]
}

export function usePartner() {
  const { authedGet } = useAuth()
  const data = ref<PartnerOverview | null>(null)
  const pending = ref(true)
  const error = ref(false)

  async function load() {
    pending.value = true
    error.value = false
    try {
      data.value = await authedGet<PartnerOverview>('/api/v1/auth/partner')
    } catch {
      error.value = true
    } finally {
      pending.value = false
    }
  }

  return { data, pending, error, load }
}
