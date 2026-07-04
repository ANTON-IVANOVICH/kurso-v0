// Thin fetch client for the kurso-api. Reads hit the existing public catalogue
// endpoints; write verbs are wired here so that once the admin API
// (auth + catalogue mutations) lands, only the store actions need swapping.

const BASE = (import.meta.env.VITE_API_BASE as string | undefined) ?? 'http://localhost:8080'

let authToken: string | null = null

/** Set the bearer token attached to every subsequent request (or clear it). */
export function setAuthToken(token: string | null) {
  authToken = token
}

export class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

function headers(hasBody: boolean): HeadersInit {
  const h: Record<string, string> = {}
  if (hasBody) h['Content-Type'] = 'application/json'
  if (authToken) h['Authorization'] = `Bearer ${authToken}`
  return h
}

async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    method,
    headers: headers(body !== undefined),
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) {
    let message = `${method} ${path} → ${res.status}`
    try {
      const data = (await res.json()) as { message?: string }
      if (data?.message) message = data.message
    } catch {
      /* non-JSON error body — keep the generic message */
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
