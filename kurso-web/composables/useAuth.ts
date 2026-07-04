import { computed } from 'vue'

// Real end-user session against kurso-api. Nuxt specifics:
// - The identity (`user`) lives in `useState`, resolved server-side from the
//   httpOnly refresh cookie (see plugins/auth.ts) so authed pages SSR correctly
//   with no hydration flash.
// - The ACCESS token is kept only in client memory (module variable) — never in
//   localStorage. The REFRESH token is the httpOnly `kurso_rt` cookie the server
//   sets and reads.

export interface AccountUser {
  id: string
  name: string
  email: string
  initials: string
}

interface ApiUser {
  id: string
  email: string
  name: string
}
interface AuthPayload {
  token: string
  user: ApiUser
}

// Client-only, in-memory access token.
let accessToken: string | null = null
export const setAccessToken = (t: string | null) => {
  accessToken = t
}
export const getAccessToken = () => accessToken

function initialsFor(name: string): string {
  const parts = name
    .trim()
    .split(/[\s.@_-]+/)
    .filter(Boolean)
  const chars = parts.length >= 2 ? parts[0][0] + parts[1][0] : name.slice(0, 2)
  return chars.toUpperCase()
}

export function toAccountUser(u: ApiUser): AccountUser {
  return { id: u.id, name: u.name, email: u.email, initials: initialsFor(u.name || u.email) }
}

export function useAuth() {
  const user = useState<AccountUser | null>('user', () => null)
  const isAuthenticated = computed(() => user.value !== null)
  const base = useApiBase()

  async function login(email: string, password: string) {
    const res = await $fetch<AuthPayload>('/api/v1/auth/login', {
      baseURL: base,
      method: 'POST',
      credentials: 'include',
      body: { email, password },
    })
    setAccessToken(res.token)
    user.value = toAccountUser(res.user)
  }

  async function register(email: string, password: string, name?: string) {
    const res = await $fetch<AuthPayload>('/api/v1/auth/register', {
      baseURL: base,
      method: 'POST',
      credentials: 'include',
      body: { email, password, name },
    })
    setAccessToken(res.token)
    user.value = toAccountUser(res.user)
  }

  async function logout() {
    try {
      await $fetch('/api/v1/auth/logout', { baseURL: base, method: 'POST', credentials: 'include' })
    } catch {
      /* clear locally regardless */
    }
    setAccessToken(null)
    user.value = null
  }

  // ensureToken returns a valid access token, refreshing from the httpOnly cookie
  // when the in-memory one is missing/expired (client-only — SSR has no token).
  async function ensureToken(): Promise<string | null> {
    if (accessToken) return accessToken
    try {
      const res = await $fetch<AuthPayload>('/api/v1/auth/refresh', {
        baseURL: base,
        method: 'POST',
        credentials: 'include',
      })
      setAccessToken(res.token)
      if (!user.value) user.value = toAccountUser(res.user)
      return res.token
    } catch {
      return null
    }
  }

  // authedGet fetches a bearer-protected endpoint, transparently refreshing once
  // on a 401. Client-side only (the access token lives in memory).
  async function authedGet<T>(path: string): Promise<T> {
    let token = getAccessToken() ?? (await ensureToken())
    const call = (t: string | null) =>
      $fetch<T>(path, {
        baseURL: base,
        credentials: 'include',
        headers: t ? { Authorization: `Bearer ${t}` } : {},
      })
    try {
      return await call(token)
    } catch (e) {
      if ((e as { status?: number; statusCode?: number })?.statusCode === 401) {
        token = await ensureToken()
        if (token) return await call(token)
      }
      throw e
    }
  }

  return { user, isAuthenticated, login, register, logout, ensureToken, authedGet }
}
