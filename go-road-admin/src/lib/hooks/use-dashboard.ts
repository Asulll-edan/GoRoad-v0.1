'use client'

import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'

export function useDashboardStats() {
  return useQuery({
    queryKey: ['admin', 'dashboard'],
    queryFn: async () => {
      const res = await api.dashboard.stats()
      return res.data
    },
    staleTime: 5 * 60 * 1000,
    refetchInterval: 60 * 1000,
  })
}
