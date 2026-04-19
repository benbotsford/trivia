// Root layout — wraps every page in the app.
// In Next.js App Router, layout.tsx at the app/ root is the outermost shell:
// it renders once and persists across navigations (the children swap out,
// the layout stays mounted). This is where global styles, fonts, and shared
// UI like the navbar live.
import type { Metadata } from 'next'
import './globals.css'
import Navbar from '@/components/Navbar'
import { getSession } from '@/lib/session'

export const metadata: Metadata = {
  title: 'Trivia',
  description: 'Host live trivia games',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const session = getSession()

  return (
    <html lang="en">
      <body className="min-h-screen bg-gray-50 text-gray-900 antialiased">
        <Navbar session={session} />
        {/* Main content area — each page renders into {children} */}
        <main className="mx-auto max-w-5xl px-6 py-8">{children}</main>
      </body>
    </html>
  )
}
