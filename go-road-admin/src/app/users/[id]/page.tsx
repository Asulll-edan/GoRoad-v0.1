'use client'

import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'

export default function UserDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { data, isLoading } = useQuery({
    queryKey: ['admin', 'users', id],
    queryFn: () => api.users.get(id),
  })

  if (isLoading) return <div className="flex items-center justify-center h-64"><div className="animate-spin w-8 h-8 border-2 border-[#1A73E8] border-t-transparent rounded-full" /></div>
  if (!data) return <div className="text-center py-12 text-gray-500">User not found</div>

  const user = data.data

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{user.username}</h1>
      <div className="card">
        <dl className="grid grid-cols-2 gap-4">
          <div><dt className="text-sm text-gray-500">Email</dt><dd>{user.email}</dd></div>
          <div><dt className="text-sm text-gray-500">Role</dt><dd>{user.role}</dd></div>
          <div><dt className="text-sm text-gray-500">Phone</dt><dd>{user.phone || '-'}</dd></div>
          <div><dt className="text-sm text-gray-500">Verified</dt><dd>{user.is_verified ? 'Yes' : 'No'}</dd></div>
          <div><dt className="text-sm text-gray-500">Joined</dt><dd>{new Date(user.created_at).toLocaleDateString()}</dd></div>
          <div><dt className="text-sm text-gray-500">Total KM</dt><dd>{(user.total_distance_km || 0).toLocaleString()} km</dd></div>
          <div><dt className="text-sm text-gray-500">Motors</dt><dd>{user.motor_count || 0}</dd></div>
          <div><dt className="text-sm text-gray-500">Badges</dt><dd>{user.badge_count || 0}</dd></div>
        </dl>
      </div>
    </div>
  )
}
