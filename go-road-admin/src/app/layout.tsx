'use client'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'
import { Sidebar } from '@/components/layout/sidebar'
import './globals.css'

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: { queries: { retry: 1, refetchOnWindowFocus: false } },
  }))

  return (
    <html lang="id">
      <body>
        <QueryClientProvider client={queryClient}>
          <div className="flex h-screen">
            <Sidebar />
            <main className="flex-1 overflow-auto p-6 bg-[var(--muted)]">
              {children}
            </main>
          </div>
        </QueryClientProvider>
      </body>
    </html>
  )
}
