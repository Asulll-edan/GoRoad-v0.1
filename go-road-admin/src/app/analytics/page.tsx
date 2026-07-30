'use client'

import dynamic from 'next/dynamic'
import { useDashboardStats } from '@/lib/hooks/use-dashboard'

const UserGrowthChart = dynamic(() => import('@/components/charts/user-growth-chart').then(m => ({ default: m.UserGrowthChart })), { ssr: false })
const RoomActivityChart = dynamic(() => import('@/components/charts/room-activity-chart').then(m => ({ default: m.RoomActivityChart })), { ssr: false })

export default function AnalyticsPage() {
  const { data, isLoading } = useDashboardStats()

  if (isLoading) return <div className="flex items-center justify-center h-64"><div className="animate-spin w-8 h-8 border-2 border-[#1A73E8] border-t-transparent rounded-full" /></div>
  if (!data) return <div className="text-center py-12 text-gray-500">No data</div>

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold">Analytics</h1>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <UserGrowthChart data={data.user_growth} />
        <RoomActivityChart data={data.room_activity} />
      </div>
    </div>
  )
}
