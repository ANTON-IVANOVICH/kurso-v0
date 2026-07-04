import { toAccountUser, setAccessToken, type AccountUser } from '~/composables/useAuth'

interface ApiUser {
  id: string
  email: string
  name: string
}

// Resolves the session so authed pages render correctly on the server:
// - On SSR, forward the browser's httpOnly refresh cookie to the API's /session
//   (non-rotating) to learn the current user; the identity is transferred to the
//   client via the Nuxt payload (no hydration mismatch, no token in the HTML).
// - On the client, if already logged in (from SSR), fetch a fresh access token
//   into memory for authenticated requests.
export default defineNuxtPlugin(async () => {
  const user = useState<AccountUser | null>('user', () => null)
  const base = useApiBase()

  if (import.meta.server) {
    const headers = useRequestHeaders(['cookie'])
    try {
      const res = await $fetch<{ user: ApiUser | null }>('/api/v1/auth/session', {
        baseURL: base,
        headers,
      })
      if (res.user) user.value = toAccountUser(res.user)
    } catch {
      /* anonymous */
    }
  } else if (user.value) {
    try {
      const res = await $fetch<{ token: string }>('/api/v1/auth/refresh', {
        baseURL: base,
        method: 'POST',
        credentials: 'include',
      })
      setAccessToken(res.token)
    } catch {
      user.value = null
    }
  }
})
