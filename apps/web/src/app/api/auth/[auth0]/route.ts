// Auth0 catch-all route handler for the Next.js App Router.
//
// This single file handles all Auth0 endpoints:
//   GET /api/auth/login    — redirect to Auth0 Universal Login
//   GET /api/auth/logout   — clear session and redirect to Auth0 logout
//   GET /api/auth/callback — exchange code for tokens and set session cookie
//   GET /api/auth/me       — return the current user profile as JSON
//
// Auth0 reads AUTH0_SECRET, AUTH0_BASE_URL, AUTH0_ISSUER_BASE_URL,
// AUTH0_CLIENT_ID, AUTH0_CLIENT_SECRET, and AUTH0_AUDIENCE automatically
// from environment variables — no explicit config object needed here.
import { handleAuth } from '@auth0/nextjs-auth0'

export const GET = handleAuth()
