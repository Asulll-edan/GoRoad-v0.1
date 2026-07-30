'use client'

import { DataTable, Column } from '@/components/data-table'
import { useInfiniteQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'
import type { ModerationAction } from '@/types'

export default function ModerationPage() {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useInfiniteQuery({
    queryKey: ['admin', 'moderation'],
    queryFn: ({ pageParam }) => api.moderation.list(pageParam as string | undefined),
    getNextPageParam: (last) => last.has_more ? last.next_cursor : undefined,
    initialPageParam: undefined as string | undefined,
    staleTime: 2 * 60 * 1000,
  })

  const actions = data?.pages.flatMap(p => p.data) ?? []

  const columns: Column<ModerationAction>[] = [
    { key: 'action', label: 'Action', render: (a) => <span className="badge-yellow">{a.action}</span> },
    { key: 'target_type', label: 'Target Type', render: (a) => <span className="badge-blue">{a.target_type}</span> },
    { key: 'target_id', label: 'Target', render: (a) => <code className="text-xs">{a.target_id.slice(0, 8)}...</code> },
    { key: 'reason', label: 'Reason', render: (a) => a.reason },
    { key: 'moderator_id', label: 'Moderator', render: (a) => <code className="text-xs">{a.moderator_id.slice(0, 8)}...</code> },
    { key: 'created_at', label: 'Date', render: (a) => new Date(a.created_at).toLocaleDateString() },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Moderation Log</h1>
      <DataTable
        columns={columns}
        data={actions}
        keyExtractor={a => a.id}
        isLoading={isLoading || isFetchingNextPage}
        onLoadMore={() => fetchNextPage()}
        hasMore={!!hasNextPage}
      />
    </div>
  )
}
