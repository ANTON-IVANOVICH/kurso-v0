import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { setAuthToken } from '../lib/api'

// Admin session. The admin auth backend (`POST /admin/auth/login`, JWT + TOTP,
// RequireAdmin middleware) is not built yet, so this is a local scaffold: it
// accepts any email/password so the console is usable in development, persists a
// placeholder token, and keeps the exact surface (`login`, `logout`, bearer
// token wiring) the real endpoint will slot into. `login()` already accepts a
// TOTP code so the form is contract-ready.

const TOKEN_KEY = 'kurso_admin_token'
const USER_KEY = 'kurso_admin_user'

export interface AdminUser {
  name: string
  email: string
  role: 'superadmin' | 'moderator'
  initials: string
}

function initialsFor(email: string): string {
  const local = email.split('@')[0] ?? email
  const parts = local.split(/[.\-_]/).filter(Boolean)
  const chars = parts.length >= 2 ? parts[0][0] + parts[1][0] : local.slice(0, 2)
  return chars.toUpperCase()
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(null)
  const user = ref<AdminUser | null>(null)
  const isAuthenticated = computed(() => token.value !== null)

  /** Re-hydrate a persisted session on app boot. */
  function restore() {
    const t = localStorage.getItem(TOKEN_KEY)
    const u = localStorage.getItem(USER_KEY)
    if (t && u) {
      token.value = t
      user.value = JSON.parse(u) as AdminUser
      setAuthToken(t)
    }
  }

  async function login(email: string, _password: string, _otp?: string) {
    // Placeholder auth — swap for `api.post('/admin/auth/login', …)` once the
    // backend exists. We keep the signature (incl. TOTP) so nothing else changes.
    const t = `dev.${btoa(email).replace(/=/g, '')}.${Date.now().toString(36)}`
    const u: AdminUser = {
      name: email.split('@')[0] ?? 'admin',
      email,
      role: 'superadmin',
      initials: initialsFor(email),
    }
    token.value = t
    user.value = u
    setAuthToken(t)
    localStorage.setItem(TOKEN_KEY, t)
    localStorage.setItem(USER_KEY, JSON.stringify(u))
  }

  function logout() {
    token.value = null
    user.value = null
    setAuthToken(null)
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  }

  return { token, user, isAuthenticated, restore, login, logout }
})
