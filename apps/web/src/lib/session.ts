// Stub session module.
//
// While Auth0 is not yet configured, every request is treated as coming from
// a hardcoded dev user. When Auth0 is wired up, replace getSession() with a
// call to @auth0/nextjs-auth0's getSession() (which reads the encrypted
// session cookie set during the OAuth callback).
//
// Nothing else in the codebase imports Auth0 directly — all auth reads go
// through this module, so the swap is a single-file change.

import type { Session } from '@/types'

// A stable fake user ID that matches the owner_id seeded in the mock bank store.
const DEV_USER_ID = 'dev-user-00000000-0000-0000-0000-000000000001'

export const DEV_SESSION: Session = {
  userId: DEV_USER_ID,
  email: 'dev@example.com',
  name: 'Dev User',
}

// getSession returns the current user's session.
// In dev: always returns the hardcoded DEV_SESSION.
// In prod (after Auth0 is set up): call getSession() from @auth0/nextjs-auth0
// and map its fields to our Session type.
export function getSession(): Session {
  return DEV_SESSION
}
