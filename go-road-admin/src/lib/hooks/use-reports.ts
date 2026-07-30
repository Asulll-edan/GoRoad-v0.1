'use client'

import { useInfiniteQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'

export function useReports() {
  return useInfiniteQuery({
    queryKey: ['admin', 'reports'],
    queryFn: ({ pageParam }) => api.reports.list(pageParam as string | undefined),
    getNextPageParam: (last) => last.has_more ? last.next_cursor : undefined,
    initialPageParam: undefined as string | undefined,
    staleTime: 2 * 60 * 1000,
  })
}

export function useEmergencyEvents() {
  return useInfiniteQuery({
    queryKey: ['admin', 'emergency'],
    queryFn: ({ pageParam }) => api.emergency.list(pageParam as string | undefined),
    getNextPageParam: (last) => last.has_more ? last.next_cursor : undefined,
    initialPageParam: undefined as string | undefined,
    staleTime: 60 * 1000,
    refetchInterval: 30 * 1000,
  })
}
