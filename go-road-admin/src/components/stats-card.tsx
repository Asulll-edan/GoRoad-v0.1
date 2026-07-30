import { clsx } from 'clsx'
import { LucideIcon } from 'lucide-react'

interface StatsCardProps {
  title: string
  value: string | number
  icon: LucideIcon
  trend?: { value: number; positive: boolean }
  className?: string
}

export function StatsCard({ title, value, icon: Icon, trend, className }: StatsCardProps) {
  return (
    <div className={clsx('card', className)}>
      <div className="flex items-start justify-between">
        <div>
          <p className="stat-label">{title}</p>
          <p className="stat-value mt-1">{value}</p>
          {trend && (
            <p className={clsx('text-sm mt-1', trend.positive ? 'text-green-600' : 'text-red-600')}>
              {trend.positive ? '↑' : '↓'} {Math.abs(trend.value)}%
            </p>
          )}
        </div>
        <div className="rounded-lg bg-[#1A73E8]/10 p-3">
          <Icon className="text-[#1A73E8]" size={24} />
        </div>
      </div>
    </div>
  )
}
