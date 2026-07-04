import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api, setAccessToken } from '../lib/api'

// Admin session. Nothing is persisted in localStorage: the access token lives in
// memory (see api.ts) and the session survives reloads through the httpOnly
// refresh cookie — `restore()` exchanges it for a fresh access token on boot.

export interface AdminUser {
  name: string
  email: string
  role: 'superadmin' | 'moderator'
  initials: string
}

interface AuthPayload {
  token: string
  admin: { id: string; email: string; role: 'superadmin' | 'moderator' }
}

function initialsFor(email: string): string {
  const local = email.split('@')[0] ?? email
  const parts = local.split(/[.\-_]/).filter(Boolean)
  const chars = parts.length >= 2 ? parts[0][0] + parts[1][0] : local.slice(0, 2)
  return chars.toUpperCase()
}

function userFrom(admin: AuthPayload['admin']): AdminUser {
  return {
    name: admin.email.split('@')[0] ?? 'admin',
    email: admin.email,
    role: admin.role,
    initials: initialsFor(admin.email),
  }
}

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AdminUser | null>(null)
  const ready = ref(false) // true once the boot-time restore attempt has settled
  const isAuthenticated = computed(() => user.value !== null)

  /** Boot: try to revive the session from the refresh cookie. */
  async function restore() {
    try {
      const res = await api.post<AuthPayload>('/admin/auth/refresh')
      setAccessToken(res.token)
      user.value = userFrom(res.admin)
    } catch {
      setAccessToken(null)
      user.value = null
    } finally {
      ready.value = true
    }
  }

  async function login(email: string, password: string, otp?: string) {
    // Throws ApiError on bad credentials (401) — the login view shows the message.
    const res = await api.post<AuthPayload>('/admin/auth/login', { email, password, otp })
    setAccessToken(res.token)
    user.value = userFrom(res.admin)
  }

  async function logout() {
    try {
      await api.post('/admin/auth/logout')
    } catch {
      /* clearing locally regardless */
    }
    setAccessToken(null)
    user.value = null
  }

  return { user, ready, isAuthenticated, restore, login, logout }
})
