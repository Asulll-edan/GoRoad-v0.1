'use client'

import { DataTable, Column } from '@/components/data-table'
import { useEmergencyEvents } from '@/lib/hooks/use-reports'
import type { EmergencyEvent } from '@/types'

export default function EmergencyPage() {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useEmergencyEvents()
  const events = data?.pages.flatMap(p => p.data) ?? []

  const columns: Column<EmergencyEvent>[] = [
    { key: 'type', label: 'Type', render: (e) => <span className="badge-red">{e.type}</span> },
    { key: 'description', label: 'Description', render: (e) => e.description || '-' },
    { key: 'user_id', label: 'User', render: (e) => <code className="text-xs">{e.user_id.slice(0, 8)}...</code> },
    { key: 'room_id', label: 'Room', render: (e) => <code className="text-xs">{e.room_id.slice(0, 8)}...</code> },
    { key: 'status', label: 'Status', render: (e) => {
      const map: Record<string, string> = { active: 'badge-red', resolved: 'badge-green', acknowledged: 'badge-yellow' }
      return <span className={map[e.status] || 'badge'}>{e.status}</span>
    }},
    { key: 'created_at', label: 'Time', render: (e) => new Date(e.created_at).toLocaleString() },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Emergency Events</h1>
      <DataTable
        columns={columns}
        data={events}
        keyExtractor={e => e.id}
        isLoading={isLoading || isFetchingNextPage}
        onLoadMore={() => fetchNextPage()}
        hasMore={!!hasNextPage}
      />
    </div>
  )
}
