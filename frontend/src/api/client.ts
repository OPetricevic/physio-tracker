// API base:
// - In dev we proxy `/api` to the Go backend via Vite (see vite.config.ts).
// - If VITE_API_URL is set (e.g., http://localhost:3600/api), we use it.
// - Otherwise default to empty string; call sites already include full paths like `/api/...`.
const API_BASE = (import.meta.env.VITE_API_URL || '').replace(/\/$/, '')

type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'DELETE' | 'PUT'

export async function apiRequest<T>(
  path: string,
  options: {
    method?: HttpMethod
    token?: string | null
    body?: unknown
  } = {},
): Promise<T> {
  const { method = 'GET', token, body } = options
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  const res = await fetch(`${API_BASE}${path}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) {
    const text = await res.text().catch(() => '')
    const err: any = new Error(text || `Request failed with status ${res.status}`)
    err.status = res.status
    throw err
  }
  // For 204 no content, return null as any
  if (res.status === 204) return null as T
  return (await res.json()) as T
}

export function getApiBase() {
  return API_BASE
}
