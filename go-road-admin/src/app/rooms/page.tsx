'use client'

import { useState } from 'react'
import Link from 'next/link'
import { DataTable, Column } from '@/components/data-table'
import { useRooms } from '@/lib/hooks/use-rooms'
import type { Room } from '@/types'

export default function RoomsPage() {
  const [status, setStatus] = useState('')
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useRooms(status || undefined)

  const rooms = data?.pages.flatMap(p => p.data) ?? []

  const columns: Column<Room>[] = [
    { key: 'name', label: 'Name', sortable: true, render: (r) => (
      <Link href={`/rooms/${r.id}`} className="text-[#1A73E8] hover:underline font-medium">{r.name}</Link>
    )},
    { key: 'origin_city', label: 'Origin', render: (r) => r.origin_city || '-' },
    { key: 'destination_city', label: 'Destination', render: (r) => r.destination_city || '-' },
    { key: 'member_count', label: 'Members', sortable: true },
    { key: 'status', label: 'Status', render: (r) => {
      const map: Record<string, string> = { active: 'badge-green', completed: 'badge-blue', cancelled: 'badge-red' }
      return <span className={map[r.status] || 'badge'}>{r.status}</span>
    }},
    { key: 'is_private', label: 'Private', render: (r) => r.is_private ? '🔒' : '🌍' },
    { key: 'created_at', label: 'Created', render: (r) => new Date(r.created_at).toLocaleDateString() },
  ]

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Rooms</h1>
        <select className="input w-40" value={status} onChange={e => setStatus(e.target.value)}>
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="completed">Completed</option>
          <option value="cancelled">Cancelled</option>
        </select>
      </div>
      <DataTable
        columns={columns}
        data={rooms}
        keyExtractor={r => r.id}
        isLoading={isLoading || isFetchingNextPage}
        onLoadMore={() => fetchNextPage()}
        hasMore={!!hasNextPage}
      />
    </div>
  )
}
