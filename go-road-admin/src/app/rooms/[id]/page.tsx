'use client'

import dynamic from 'next/dynamic'
import { useParams } from 'next/navigation'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api-client'

const RoomLiveMap = dynamic(() => import('@/components/room-live-map'), { ssr: false, loading: () => <div className="card h-64 flex items-center justify-center text-gray-400">Loading map...</div> })

export default function RoomDetailPage() {
  const { id } = useParams<{ id: string }>()
  const { data, isLoading } = useQuery({
    queryKey: ['admin', 'rooms', id],
    queryFn: () => api.rooms.get(id),
    staleTime: 5 * 60 * 1000,
  })

  if (isLoading) return <div className="flex items-center justify-center h-64"><div className="animate-spin w-8 h-8 border-2 border-[#1A73E8] border-t-transparent rounded-full" /></div>
  if (!data) return <div className="text-center py-12 text-gray-500">Room not found</div>

  const room = data.data

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">{room.name}</h1>
      <div className="card">
        <dl className="grid grid-cols-2 gap-4">
          <div><dt className="text-sm text-gray-500">Description</dt><dd>{room.description || '-'}</dd></div>
          <div><dt className="text-sm text-gray-500">Category</dt><dd>{room.category || '-'}</dd></div>
          <div><dt className="text-sm text-gray-500">Members</dt><dd>{room.member_count}/{room.max_members}</dd></div>
          <div><dt className="text-sm text-gray-500">Status</dt><dd>{room.status}</dd></div>
          <div><dt className="text-sm text-gray-500">Origin</dt><dd>{room.origin_city || '-'}</dd></div>
          <div><dt className="text-sm text-gray-500">Destination</dt><dd>{room.destination_city || '-'}</dd></div>
        </dl>
      </div>
      <RoomLiveMap />
    </div>
  )
}
