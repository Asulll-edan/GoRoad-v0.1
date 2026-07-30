'use client'

import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface UserGrowthChartProps {
  data: { date: string; count: number }[]
}

export function UserGrowthChart({ data }: UserGrowthChartProps) {
  return (
    <div className="card">
      <h3 className="font-semibold mb-4">User Growth</h3>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" />
          <XAxis dataKey="date" tick={{ fontSize: 12 }} stroke="#999" />
          <YAxis tick={{ fontSize: 12 }} stroke="#999" />
          <Tooltip />
          <Line type="monotone" dataKey="count" stroke="#1A73E8" strokeWidth={2} dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
