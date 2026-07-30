'use client'

import { useState } from 'react'
import Link from 'next/link'
import { DataTable, Column } from '@/components/data-table'
import { useUsers } from '@/lib/hooks/use-users'
import type { User } from '@/types'

export default function UsersPage() {
  const [search, setSearch] = useState('')
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } = useUsers(search)

  const users = data?.pages.flatMap(p => p.data) ?? []

  const columns: Column<User>[] = [
    { key: 'username', label: 'Username', sortable: true, render: (u) => (
      <Link href={`/users/${u.id}`} className="text-[#1A73E8] hover:underline font-medium">{u.username}</Link>
    )},
    { key: 'email', label: 'Email' },
    { key: 'role', label: 'Role', render: (u) => (
      <span className={u.role === 'admin' ? 'badge-red' : u.role === 'moderator' ? 'badge-yellow' : 'badge-blue'}>{u.role}</span>
    )},
    { key: 'is_verified', label: 'Verified', render: (u) => u.is_verified ? <span className="badge-green">Yes</span> : <span className="badge">No</span> },
    { key: 'motor_count', label: 'Motors' },
    { key: 'badge_count', label: 'Badges' },
    { key: 'created_at', label: 'Joined', render: (u) => new Date(u.created_at).toLocaleDateString() },
  ]

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Users</h1>
      <DataTable
        columns={columns}
        data={users}
        keyExtractor={u => u.id}
        isLoading={isLoading || isFetchingNextPage}
        onLoadMore={() => fetchNextPage()}
        hasMore={!!hasNextPage}
        searchPlaceholder="Search by username or email..."
        onSearch={setSearch}
      />
    </div>
  )
}
