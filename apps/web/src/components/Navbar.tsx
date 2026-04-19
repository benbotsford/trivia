// Navbar is a server component — it has no interactivity so it doesn't need
// 'use client'. Server components render on the server and send plain HTML,
// which is faster than shipping the component's JS to the browser.
import type { Session } from '@/types'

interface NavbarProps {
  session: Session
}

export default function Navbar({ session }: NavbarProps) {
  return (
    <header className="bg-brand-blue">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-6 py-4">
        {/* App name / logo */}
        <span className="text-lg font-semibold tracking-tight text-white">
          Trivia
        </span>

        {/* Stubbed user identity — replaced with Auth0 profile + logout button later */}
        <div className="flex items-center gap-3">
          {/* Initials avatar — inverted colors so it pops against the blue header */}
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-white text-xs font-semibold text-brand-blue">
            {session.name
              .split(' ')
              .map((w) => w[0])
              .join('')
              .toUpperCase()
              .slice(0, 2)}
          </div>
          <span className="hidden text-sm text-white/70 sm:block">{session.email}</span>
        </div>
      </div>
    </header>
  )
}
