'use client'

import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface RoomActivityChartProps {
  data: { date: string; count: number }[]
}

export function RoomActivityChart({ data }: RoomActivityChartProps) {
  return (
    <div className="card">
      <h3 className="font-semibold mb-4">Room Activity</h3>
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
          <XAxis dataKey="date" tick={{ fontSize: 12 }} stroke="#999" />
          <YAxis tick={{ fontSize: 12 }} stroke="#999" />
          <Tooltip />
          <Bar dataKey="count" fill="#1A73E8" radius={[4, 4, 0, 0]} />
        </BarChart>
      </ResponsiveContainer>
    </div>
  )
}
