// Route protection middleware.
//
// When Auth0 is fully configured (AUTH0_SECRET is set), withMiddlewareAuthRequired
// redirects unauthenticated users to /api/auth/login on every matched route.
//
// When AUTH0_SECRET is absent (local dev with DEV_AUTH_TOKEN), the middleware
// is a passthrough — the host layout's own guard handles the dev bypass instead.
//
// The matcher excludes:
//  - /api/auth/*          Auth0 callback/logout/me endpoints — must be public
//  - /join and /play/*    Player-facing pages — no host login required
//  - Next.js internals    _next/static, _next/image, favicon.ico

import { withMiddlewareAuthRequired } from '@auth0/nextjs-auth0/edge'
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const isAuth0Configured = Boolean(process.env.AUTH0_SECRET)

export default isAuth0Configured
  ? withMiddlewareAuthRequired()
  : (_req: NextRequest) => NextResponse.next()

export const config = {
  matcher: [
    '/((?!api/auth|_next/static|_next/image|favicon\\.ico|join|play).*)',
  ],
}
