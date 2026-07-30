const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1'

async function fetchApi<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options?.headers },
    ...options,
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || `API error ${res.status}`)
  }
  return res.json()
}

export function buildQuery(params: Record<string, string | number | undefined>): string {
  const q = new URLSearchParams()
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined) q.set(k, String(v))
  }
  const s = q.toString()
  return s ? `?${s}` : ''
}

export const api = {
  dashboard: {
    stats: () => fetchApi<{ data: import('@/types').DashboardStats }>('/admin/dashboard'),
  },
  users: {
    list: (cursor?: string, limit = 50, search?: string) =>
      fetchApi<import('@/types').CursorPage<import('@/types').User>>(`/admin/users${buildQuery({ cursor, limit, q: search })}`),
    get: (id: string) => fetchApi<{ data: import('@/types').User }>(`/admin/users/${id}`),
  },
  rooms: {
    list: (cursor?: string, limit = 50, status?: string) =>
      fetchApi<import('@/types').CursorPage<import('@/types').Room>>(`/admin/rooms${buildQuery({ cursor, limit, status })}`),
    get: (id: string) => fetchApi<{ data: import('@/types').Room }>(`/admin/rooms/${id}`),
  },
  emergency: {
    list: (cursor?: string, limit = 50) =>
      fetchApi<import('@/types').CursorPage<import('@/types').EmergencyEvent>>(`/admin/emergency${buildQuery({ cursor, limit })}`),
  },
  reports: {
    list: (cursor?: string, limit = 50) =>
      fetchApi<import('@/types').CursorPage<import('@/types').Report>>(`/admin/reports${buildQuery({ cursor, limit })}`),
  },
  moderation: {
    list: (cursor?: string, limit = 50) =>
      fetchApi<import('@/types').CursorPage<import('@/types').ModerationAction>>(`/admin/moderation${buildQuery({ cursor, limit })}`),
  },
}
