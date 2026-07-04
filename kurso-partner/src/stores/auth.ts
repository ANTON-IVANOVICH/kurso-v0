import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api, setAccessToken } from '../lib/api'
import type { Merchant } from '../types/models'

// Merchant session. Nothing is persisted in localStorage: the access token lives
// in memory (see api.ts) and the session survives reloads through the httpOnly
// refresh cookie — `restore()` exchanges it for a fresh access token on boot.

interface AuthPayload {
  token: string
  merchant: Merchant
}

function initialsFor(name: string): string {
  const parts = name.split(/[\s._-]+/).filter(Boolean)
  const chars = parts.length >= 2 ? parts[0][0] + parts[1][0] : name.slice(0, 2)
  return chars.toUpperCase()
}

export const useAuthStore = defineStore('auth', () => {
  const merchant = ref<Merchant | null>(null)
  const ready = ref(false) // true once the boot-time restore attempt has settled
  const isAuthenticated = computed(() => merchant.value !== null)
  const initials = computed(() =>
    merchant.value ? initialsFor(merchant.value.exchangerName) : 'CB',
  )

  /** Boot: try to revive the session from the refresh cookie. */
  async function restore() {
    try {
      const res = await api.post<AuthPayload>('/partner/auth/refresh')
      setAccessToken(res.token)
      merchant.value = res.merchant
    } catch {
      setAccessToken(null)
      merchant.value = null
    } finally {
      ready.value = true
    }
  }

  async function login(email: string, password: string, otp?: string) {
    // Throws ApiError on bad credentials (401) — the login view shows the message.
    const res = await api.post<AuthPayload>('/partner/auth/login', { email, password, otp })
    setAccessToken(res.token)
    merchant.value = res.merchant
  }

  async function logout() {
    try {
      await api.post('/partner/auth/logout')
    } catch {
      /* clearing locally regardless */
    }
    setAccessToken(null)
    merchant.value = null
  }

  return { merchant, ready, isAuthenticated, initials, restore, login, logout }
})
