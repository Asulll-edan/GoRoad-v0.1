'use client'

import { DataTable, Column } from '@/components/data-table'
import { useReports } from '@/lib/hooks/use-reports'
import type { Report } from '@/types'

export default function ReportsPage() {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useReports()
  const reports = data?.pages.flatMap(p => p.data) ?? []

  const columns: Column<Report>[] = [
    { key: 'type', label: 'Type', render: (r) => <span className="badge-blue">{r.type}</span> },
    { key: 'reason', label: 'Reason', render: (r) => <span className="max-w-xs truncate block">{r.reason}</span> },
    { key: 'reporter_id', label: 'Reporter', render: (r) => <code className="text-xs">{r.reporter_id.slice(0, 8)}...</code> },
    { key: 'reported_id', label: 'Reported', render: (r) => <code className="text-xs">{r.reported_id.slice(0, 8)}...</code> },
    { key: 'status', label: 'Status', render: (r) => {
      const map: Record<string, string> = { pending: 'badge-yellow', resolved: 'badge-green', dismissed: 'badge' }
      return <span className={map[r.status] || 'badge'}>{r.status}</span>
    }},
    { key: 'created_at', label: 'Date', render: (r) => new Date(r.created_at).toLocaleDateString() },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Reports</h1>
      <DataTable
        columns={columns}
        data={reports}
        keyExtractor={r => r.id}
        isLoading={isLoading || isFetchingNextPage}
        onLoadMore={() => fetchNextPage()}
        hasMore={!!hasNextPage}
      />
    </div>
  )
}
