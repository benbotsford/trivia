// Session helper — thin wrapper over @auth0/nextjs-auth0's getSession().
//
// When Auth0 is configured (AUTH0_SECRET is set), reads the encrypted session
// cookie and maps the Auth0 user object to the app's Session type.
//
// When AUTH0_SECRET is absent (local dev with DEV_AUTH_TOKEN), returns null —
// callers should treat null + DEV_AUTH_TOKEN as "authenticated as dev user".
// The host layout handles this case explicitly.
//
// Call from server components and server actions only.

import type { Session } from '@/types'

const isAuth0Configured = Boolean(process.env.AUTH0_SECRET)

export async function getSession(): Promise<Session | null> {
  if (!isAuth0Configured) return null

  // Dynamically import so the Auth0 SDK never initializes in dev mode.
  const { getSession: auth0GetSession } = await import('@auth0/nextjs-auth0')
  const session = await auth0GetSession()
  if (!session) return null

  return {
    userId: session.user.sub,
    email: session.user.email ?? '',
    name: session.user.name ?? session.user.email ?? '',
  }
}
