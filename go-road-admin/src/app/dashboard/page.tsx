'use client'

import { Users, DoorOpen, Route, AlertTriangle, Bike, FileText } from 'lucide-react'
import { StatsCard } from '@/components/stats-card'
import { UserGrowthChart } from '@/components/charts/user-growth-chart'
import { RoomActivityChart } from '@/components/charts/room-activity-chart'
import { useDashboardStats } from '@/lib/hooks/use-dashboard'

export default function DashboardPage() {
  const { data, isLoading } = useDashboardStats()

  if (isLoading) return <div className="flex items-center justify-center h-64"><div className="animate-spin w-8 h-8 border-2 border-[#1A73E8] border-t-transparent rounded-full" /></div>
  if (!data) return <div className="text-center py-12 text-gray-500">Failed to load dashboard</div>

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <StatsCard title="Total Users" value={data.total_users} icon={Users} />
        <StatsCard title="Active 24h" value={data.active_users_24h} icon={Bike} />
        <StatsCard title="Total Rooms" value={data.total_rooms} icon={DoorOpen} />
        <StatsCard title="Active Rooms" value={data.active_rooms} icon={DoorOpen} />
        <StatsCard title="Total Routes" value={data.total_routes} icon={Route} />
        <StatsCard title="Total KM" value={`${(data.total_distance_km / 1000).toFixed(0)}k`} icon={Route} />
        <StatsCard title="Emergency 24h" value={data.emergency_events_24h} icon={AlertTriangle} className={data.emergency_events_24h > 0 ? 'border-red-200' : ''} />
        <StatsCard title="Pending Reports" value={data.reports_pending} icon={FileText} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <UserGrowthChart data={data.user_growth} />
        <RoomActivityChart data={data.room_activity} />
      </div>

      <div className="card">
        <h3 className="font-semibold mb-4">Top Rooms</h3>
        <div className="space-y-3">
          {data.top_rooms.map((room, i) => (
            <div key={room.id} className="flex items-center justify-between">
              <span className="text-sm"><span className="font-medium text-gray-400 mr-2">#{i + 1}</span>{room.name}</span>
              <span className="text-sm text-gray-500">{room.member_count} members</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
