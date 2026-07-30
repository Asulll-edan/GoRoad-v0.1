'use client'

import { useInfiniteQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'

export function useRooms(status?: string) {
  return useInfiniteQuery({
    queryKey: ['admin', 'rooms', status],
    queryFn: ({ pageParam }) => api.rooms.list(pageParam as string | undefined, 50, status),
    getNextPageParam: (last) => last.has_more ? last.next_cursor : undefined,
    initialPageParam: undefined as string | undefined,
    staleTime: 2 * 60 * 1000,
  })
}
