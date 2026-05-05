// API client for question banks. Server-side only.

import type { Bank } from '@/types'
import { apiFetch } from './client'

// listBanks fetches all banks owned by the authenticated user.
// The API derives the owner from the Bearer token — no user ID needed here.
export async function listBanks(): Promise<Bank[]> {
  const res = await apiFetch('/banks')
  return res.json()
}

// createBank creates a new bank and returns the created record.
export async function createBank(name: string, description?: string): Promise<Bank> {
  const res = await apiFetch('/banks', {
    method: 'POST',
    body: JSON.stringify({ name, description }),
  })
  return res.json()
}

// updateBank replaces a bank's name and description, returning the updated record.
export async function updateBank(id: string, name: string, description?: string): Promise<Bank> {
  const res = await apiFetch(`/banks/${id}`, {
    method: 'PUT',
    body: JSON.stringify({ name, description }),
  })
  return res.json()
}

// deleteBank permanently removes a bank.
// The Go API returns 204 No Content on success — no body to parse.
export async function deleteBank(id: string): Promise<void> {
  await apiFetch(`/banks/${id}`, { method: 'DELETE' })
}

// getBank fetches a single bank by ID.
export async function getBank(id: string): Promise<Bank> {
  const res = await apiFetch(`/banks/${id}`)
  return res.json()
}
