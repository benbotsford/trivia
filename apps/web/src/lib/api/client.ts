// Shared server-side API client.
//
// All API modules import apiFetch from here instead of rolling their own.
// Token resolution order:
//   1. DEV_AUTH_TOKEN — set in .env.local for local development. When present,
//      requests are sent with this token and the Go API's dev bypass handles them
//      without contacting Auth0. Never set this in production.
//   2. Auth0 access token — retrieved from the encrypted session cookie via
//      getAccessToken(). The token carries the configured AUTH0_AUDIENCE so the
//      Go API can validate it against Auth0's JWKS.
//
// This module is server-only. It must never be imported from a client component.

import { getAccessToken } from '@auth0/nextjs-auth0'

const API_BASE = (process.env.API_URL ?? 'http://localhost:8080').replace(/\/$/, '')

// getAuthToken returns the bearer token for the current request.
// Exported for use by server components that need the token for purposes
// other than HTTP API calls (e.g. passing it to a WebSocket via query param).
export async function getAuthToken(): Promise<string> {
  const devToken = process.env.DEV_AUTH_TOKEN
  if (devToken) return devToken

  const { accessToken } = await getAccessToken()
  if (!accessToken) throw new Error('No Auth0 access token in session')
  return accessToken
}

// apiFetch is a thin wrapper around fetch that:
//  - Prefixes the path with the Go API base URL
//  - Resolves and attaches a Bearer token (dev bypass or Auth0 access token)
//  - Sets cache: 'no-store' so Next.js never serves a stale API response
//  - Throws a descriptive error on non-2xx responses, parsing the Go API's
//    {"error":"..."} JSON body when available
export async function apiFetch(path: string, options: RequestInit = {}): Promise<Response> {
  const token = await getAuthToken()
  const url = `${API_BASE}${path}`

  const res = await fetch(url, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
      ...options.headers,
    },
    cache: 'no-store',
  })

  if (!res.ok) {
    const body = await res.text().catch(() => '')
    let message = `API ${options.method ?? 'GET'} ${path} → ${res.status}`
    try {
      const json = JSON.parse(body)
      if (json?.error) message = json.error
    } catch {
      if (body) message = body
    }
    throw new Error(message)
  }

  return res
}
