'use client'

import { useState } from 'react'
import { ChevronDown, ChevronUp, ChevronsUpDown, Search } from 'lucide-react'

export interface Column<T> {
  key: string
  label: string
  sortable?: boolean
  render?: (item: T) => React.ReactNode
}

interface DataTableProps<T> {
  columns: Column<T>[]
  data: T[]
  keyExtractor: (item: T) => string
  isLoading?: boolean
  onLoadMore?: () => void
  hasMore?: boolean
  searchPlaceholder?: string
  onSearch?: (q: string) => void
}

export function DataTable<T>({
  columns, data, keyExtractor, isLoading, onLoadMore, hasMore, searchPlaceholder, onSearch,
}: DataTableProps<T>) {
  const [sortKey, setSortKey] = useState<string | null>(null)
  const [sortDir, setSortDir] = useState<'asc' | 'desc'>('asc')

  const handleSort = (key: string) => {
    if (sortKey === key) {
      setSortDir(d => d === 'asc' ? 'desc' : 'asc')
    } else {
      setSortKey(key)
      setSortDir('asc')
    }
  }

  return (
    <div className="card p-0 overflow-hidden">
      {onSearch && (
        <div className="flex items-center gap-2 p-4 border-b">
          <Search size={18} className="text-gray-400" />
          <input
            className="input flex-1"
            placeholder={searchPlaceholder || 'Search...'}
            onChange={e => onSearch(e.target.value)}
          />
        </div>
      )}
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-gray-50">
              {columns.map(col => (
                <th key={col.key}
                  className={`px-4 py-3 text-left font-medium text-gray-500 ${col.sortable ? 'cursor-pointer select-none hover:text-gray-700' : ''}`}
                  onClick={() => col.sortable && handleSort(col.key)}
                >
                  <span className="flex items-center gap-1">
                    {col.label}
                    {col.sortable && (
                      sortKey === col.key
                        ? (sortDir === 'asc' ? <ChevronUp size={14} /> : <ChevronDown size={14} />)
                        : <ChevronsUpDown size={14} className="text-gray-300" />
                    )}
                  </span>
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {data.map(item => (
              <tr key={keyExtractor(item)} className="border-b last:border-0 hover:bg-gray-50">
                {columns.map(col => (
                  <td key={col.key} className="px-4 py-3">
                    {col.render ? col.render(item) : String((item as any)[col.key] ?? '')}
                  </td>
                ))}
              </tr>
            ))}
            {isLoading && (
              <tr><td colSpan={columns.length} className="text-center py-8 text-gray-400">Loading...</td></tr>
            )}
          </tbody>
        </table>
      </div>
      {hasMore && (
        <div className="flex justify-center p-4 border-t">
          <button className="btn-ghost" onClick={onLoadMore} disabled={isLoading}>
            {isLoading ? 'Loading...' : 'Load More'}
          </button>
        </div>
      )}
    </div>
  )
}
