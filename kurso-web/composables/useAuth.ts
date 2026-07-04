import { computed } from 'vue'

// Client session. The account backend (JWT / Telegram / Google) isn't built yet,
// so the session lives in an SSR-readable cookie: logging in persists a real
// session that survives reloads and is visible to the server (so the header and
// route guard render correctly without a hydration flash). Swap login() for the
// real /auth calls when the backend lands — nothing else changes.

export interface AccountUser {
  name: string
  initials: string
  plan: string
  email?: string
  via?: 'telegram' | 'google' | 'apple' | 'email'
}

function initialsFrom(name: string): string {
  const parts = name
    .trim()
    .split(/[\s.]+/)
    .filter(Boolean)
  const chars = parts.length >= 2 ? parts[0][0] + parts[1][0] : name.slice(0, 2)
  return chars.toUpperCase()
}

export function useAuth() {
  const user = useCookie<AccountUser | null>('kurso_session', {
    default: () => null,
    sameSite: 'lax',
    maxAge: 60 * 60 * 24 * 30,
  })
  const isAuthenticated = computed(() => user.value !== null)

  function login(payload: { name?: string; email?: string; via?: AccountUser['via'] } = {}) {
    const name = payload.name || payload.email?.split('@')[0] || 'Пользователь'
    user.value = {
      name,
      initials: initialsFrom(name),
      plan: 'Free',
      email: payload.email,
      via: payload.via,
    }
  }

  function logout() {
    user.value = null
  }

  return { user, isAuthenticated, login, logout }
}
