// Fetch client for the kurso-api admin endpoints.
//
// Token model (matches the backend): the ACCESS token lives only in memory here
// (never localStorage), and rides as a Bearer header. The REFRESH token is an
// httpOnly cookie the browser sends automatically (`credentials: 'include'`) and
// only the server can read. On a 401 the client transparently hits
// `/admin/auth/refresh` once (single-flight) and retries the original request;
// if that fails the session is gone and the caller gets the 401.

const BASE = (import.meta.env.VITE_API_BASE as string | undefined) ?? 'http://localhost:8080'

let accessToken: string | null = null

export function setAccessToken(token: string | null) {
  accessToken = token
}
export function getAccessToken() {
  return accessToken
}

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

function headers(hasBody: boolean, withAuth: boolean): HeadersInit {
  const h: Record<string, string> = {}
  if (hasBody) h['Content-Type'] = 'application/json'
  if (withAuth && accessToken) h['Authorization'] = `Bearer ${accessToken}`
  return h
}

function raw(method: string, path: string, body?: unknown, withAuth = true): Promise<Response> {
  return fetch(`${BASE}${path}`, {
    method,
    credentials: 'include', // send/receive the httpOnly refresh cookie
    headers: headers(body !== undefined, withAuth),
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
}

// Single-flight refresh: concurrent 401s share one refresh call.
let refreshing: Promise<boolean> | null = null

export function refreshAccess(): Promise<boolean> {
  if (!refreshing) {
    refreshing = raw('POST', '/admin/auth/refresh', undefined, false)
      .then(async (res) => {
        if (!res.ok) {
          accessToken = null
          return false
        }
        const data = (await res.json()) as { token: string }
        accessToken = data.token
        return true
      })
      .catch(() => {
        accessToken = null
        return false
      })
      .finally(() => {
        refreshing = null
      })
  }
  return refreshing
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const isAuthCall = path.startsWith('/admin/auth/')
  let res = await raw(method, path, body)

  // Access expired — refresh once and retry (but never for the auth calls
  // themselves, to avoid loops).
  if (res.status === 401 && !isAuthCall) {
    if (await refreshAccess()) res = await raw(method, path, body)
  }

  if (!res.ok) {
    let message = `${method} ${path} → ${res.status}`
    try {
      const data = (await res.json()) as { message?: string }
      if (data?.message) message = data.message
    } catch {
      /* non-JSON body */
    }
    throw new ApiError(res.status, message)
  }
  if (res.status === 204) return undefined as T
  return (await res.json()) as T
}

export const api = {
  base: BASE,
  get: <T>(path: string) => request<T>('GET', path),
  post: <T>(path: string, body?: unknown) => request<T>('POST', path, body),
  patch: <T>(path: string, body?: unknown) => request<T>('PATCH', path, body),
  del: <T>(path: string) => request<T>('DELETE', path),
}
