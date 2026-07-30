import { render, screen } from '@testing-library/react'
import { DataTable, Column } from '@/components/data-table'

interface TestItem {
  id: string
  name: string
}

describe('DataTable', () => {
  const columns: Column<TestItem>[] = [
    { key: 'name', label: 'Name', sortable: true },
    { key: 'id', label: 'ID' },
  ]

  const data: TestItem[] = [
    { id: '1', name: 'Item A' },
    { id: '2', name: 'Item B' },
  ]

  it('renders headers', () => {
    render(<DataTable columns={columns} data={data} keyExtractor={i => i.id} />)
    expect(screen.getByText('Name')).toBeInTheDocument()
    expect(screen.getByText('ID')).toBeInTheDocument()
  })

  it('renders data rows', () => {
    render(<DataTable columns={columns} data={data} keyExtractor={i => i.id} />)
    expect(screen.getByText('Item A')).toBeInTheDocument()
    expect(screen.getByText('Item B')).toBeInTheDocument()
  })

  it('shows loading state', () => {
    render(<DataTable columns={columns} data={[]} keyExtractor={i => i.id} isLoading />)
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })

  it('shows load more button when hasMore', () => {
    const onLoadMore = jest.fn()
    render(<DataTable columns={columns} data={data} keyExtractor={i => i.id} hasMore onLoadMore={onLoadMore} />)
    expect(screen.getByText('Load More')).toBeInTheDocument()
  })
})
