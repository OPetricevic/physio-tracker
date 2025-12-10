const API_BASE = import.meta.env.VITE_API_URL ?? 'http://localhost:3600'

type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'DELETE'

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
    throw new Error(text || `Request failed with status ${res.status}`)
  }
  // For 204 no content, return null as any
  if (res.status === 204) return null as T
  return (await res.json()) as T
}

export function getApiBase() {
  return API_BASE
}
