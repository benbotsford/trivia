// Server-side API client for the /me endpoints.

import type { UserProfile } from '@/types'
import { apiFetch } from './client'

// getMe fetches the authenticated host's profile.
export async function getMe(): Promise<UserProfile> {
  return (await apiFetch('/me')).json()
}

// updateMe patches display_name and/or email.
// Only the fields present in `data` are changed; omitted fields are preserved.
export async function updateMe(data: {
  display_name?: string
  email?: string
}): Promise<UserProfile> {
  return (await apiFetch('/me', {
    method: 'PATCH',
    body: JSON.stringify(data),
  })).json()
}
