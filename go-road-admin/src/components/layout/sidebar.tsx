'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  LayoutDashboard, Users, DoorOpen, FileText, BarChart3,
  AlertTriangle, Shield, Settings, ChevronLeft, ChevronRight, Bike,
} from 'lucide-react'
import { useState } from 'react'
import { clsx } from 'clsx'

const navItems = [
  { href: '/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { href: '/users', label: 'Users', icon: Users },
  { href: '/rooms', label: 'Rooms', icon: DoorOpen },
  { href: '/reports', label: 'Reports', icon: FileText },
  { href: '/analytics', label: 'Analytics', icon: BarChart3 },
  { href: '/emergency', label: 'Emergency', icon: AlertTriangle },
  { href: '/moderation', label: 'Moderation', icon: Shield },
  { href: '/settings', label: 'Settings', icon: Settings },
]

export function Sidebar() {
  const pathname = usePathname()
  const [collapsed, setCollapsed] = useState(false)

  return (
    <aside className={clsx(
      'flex flex-col border-r bg-white transition-all duration-200',
      collapsed ? 'w-16' : 'w-60'
    )}>
      <div className="flex items-center gap-3 px-4 h-16 border-b">
        <Bike className="text-[#1A73E8]" size={28} />
        {!collapsed && <span className="font-bold text-lg">Go Road Admin</span>}
      </div>

      <nav className="flex-1 p-2 space-y-1">
        {navItems.map(({ href, label, icon: Icon }) => {
          const active = pathname.startsWith(href)
          return (
            <Link key={href} href={href}
              className={clsx(
                'flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors',
                active ? 'bg-[#1A73E8] text-white' : 'text-gray-600 hover:bg-gray-100'
              )}
              title={collapsed ? label : undefined}
            >
              <Icon size={20} />
              {!collapsed && <span>{label}</span>}
            </Link>
          )
        })}
      </nav>

      <button onClick={() => setCollapsed(!collapsed)}
        className="flex items-center justify-center h-12 border-t text-gray-400 hover:text-gray-600"
      >
        {collapsed ? <ChevronRight size={20} /> : <ChevronLeft size={20} />}
      </button>
    </aside>
  )
}
