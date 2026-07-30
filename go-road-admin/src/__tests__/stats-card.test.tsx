import { render, screen } from '@testing-library/react'
import { Users } from 'lucide-react'
import { StatsCard } from '@/components/stats-card'

describe('StatsCard', () => {
  it('renders title and value', () => {
    render(<StatsCard title="Total Users" value={1234} icon={Users} />)
    expect(screen.getByText('Total Users')).toBeInTheDocument()
    expect(screen.getByText('1234')).toBeInTheDocument()
  })

  it('renders trend when provided', () => {
    render(<StatsCard title="Users" value={100} icon={Users} trend={{ value: 12, positive: true }} />)
    expect(screen.getByText(/↑/)).toBeInTheDocument()
  })

  it('renders negative trend', () => {
    render(<StatsCard title="Users" value={100} icon={Users} trend={{ value: 5, positive: false }} />)
    expect(screen.getByText(/↓/)).toBeInTheDocument()
  })
})
