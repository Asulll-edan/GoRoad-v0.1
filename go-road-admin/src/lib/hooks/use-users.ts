'use client'

import { useInfiniteQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'

export function useUsers(search?: string) {
  return useInfiniteQuery({
    queryKey: ['admin', 'users', search],
    queryFn: ({ pageParam }) => api.users.list(pageParam as string | undefined, 50, search),
    getNextPageParam: (last) => last.has_more ? last.next_cursor : undefined,
    initialPageParam: undefined as string | undefined,
    staleTime: 2 * 60 * 1000,
  })
}
