'use client'

import { useState } from 'react'
import { Save } from 'lucide-react'

export default function SettingsPage() {
  const [appName, setAppName] = useState('Go Road')
  const [maintenanceMode, setMaintenanceMode] = useState(false)
  const [maxRoomMembers, setMaxRoomMembers] = useState('50')

  return (
    <div className="space-y-6 max-w-2xl">
      <h1 className="text-2xl font-bold">Settings</h1>

      <div className="card space-y-6">
        <div>
          <label className="block text-sm font-medium mb-1">App Name</label>
          <input className="input" value={appName} onChange={e => setAppName(e.target.value)} />
        </div>

        <div>
          <label className="block text-sm font-medium mb-1">Default Max Room Members</label>
          <input className="input w-32" type="number" value={maxRoomMembers} onChange={e => setMaxRoomMembers(e.target.value)} />
        </div>

        <div className="flex items-center justify-between">
          <div>
            <p className="font-medium text-sm">Maintenance Mode</p>
            <p className="text-sm text-gray-500">Block all non-admin access</p>
          </div>
          <button
            className={`relative w-12 h-6 rounded-full transition-colors ${maintenanceMode ? 'bg-red-500' : 'bg-gray-300'}`}
            onClick={() => setMaintenanceMode(!maintenanceMode)}
          >
            <div className={`absolute top-0.5 w-5 h-5 bg-white rounded-full shadow transition-transform ${maintenanceMode ? 'translate-x-6' : 'translate-x-0.5'}`} />
          </button>
        </div>

        <button className="btn-primary" onClick={() => alert('Settings saved (mock)')}>
          <Save size={16} className="mr-2" /> Save Settings
        </button>
      </div>
    </div>
  )
}
