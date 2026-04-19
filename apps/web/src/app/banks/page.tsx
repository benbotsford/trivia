// Banks page — server component.
// Because this is a server component (no 'use client' directive), it runs only
// on the server. It can be async and call the data layer directly — no useEffect,
// no loading spinner for the initial render. The result is streamed to the browser
// as HTML before any JavaScript runs.
import { getSession } from '@/lib/session'
import { listBanks } from '@/lib/mock/banks'
import BanksView from '@/components/BanksView'

export const metadata = { title: 'Question Banks — Trivia' }

export default async function BanksPage() {
  const session = getSession()
  // Fetch banks on the server — no network round-trip from the browser.
  // When auth is real, getSession() will read an encrypted cookie and listBanks()
  // will call the Go API with a Bearer token, all server-side.
  const banks = await listBanks(session.userId)

  // Pass the data down to the client component that handles interactivity.
  return <BanksView banks={banks} />
}
