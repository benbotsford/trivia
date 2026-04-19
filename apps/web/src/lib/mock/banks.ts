// Mock bank store — in-memory, resets on server restart.
//
// This module stands in for real API calls while Auth0 is not yet configured.
// The function signatures (listBanks, createBank, updateBank, deleteBank) are
// intentionally identical to what a real API client module would export, so
// swapping it out is a one-line import change in actions.ts.
//
// To replace with real API calls:
//   1. Create src/lib/api/banks.ts that fetches /api/banks/* with a real Bearer token
//   2. In src/app/banks/actions.ts, change:
//        import * as banksApi from '@/lib/mock/banks'
//      to:
//        import * as banksApi from '@/lib/api/banks'

import type { Bank } from '@/types'

// Module-level Map acts as an in-memory database table.
// Because Next.js runs server code in a single long-lived Node.js process during
// development, this Map persists across requests (but resets on `npm run dev` restart).
const store = new Map<string, Bank>([
  [
    'bank-seed-1',
    {
      id: 'bank-seed-1',
      owner_id: 'dev-user-00000000-0000-0000-0000-000000000001',
      name: 'General Knowledge',
      description: 'A mix of trivia across all topics.',
      created_at: new Date(Date.now() - 86_400_000 * 3).toISOString(), // 3 days ago
      updated_at: new Date(Date.now() - 86_400_000 * 3).toISOString(),
    },
  ],
  [
    'bank-seed-2',
    {
      id: 'bank-seed-2',
      owner_id: 'dev-user-00000000-0000-0000-0000-000000000001',
      name: 'Science & Nature',
      description: undefined, // no description — tests the optional field path
      created_at: new Date(Date.now() - 86_400_000).toISOString(), // 1 day ago
      updated_at: new Date(Date.now() - 86_400_000).toISOString(),
    },
  ],
])

// listBanks returns all banks owned by the given user, newest first.
export async function listBanks(ownerId: string): Promise<Bank[]> {
  return [...store.values()]
    .filter((b) => b.owner_id === ownerId)
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
}

// createBank inserts a new bank and returns it.
// crypto.randomUUID() is available globally in Node 19+ and in all modern browsers.
export async function createBank(
  ownerId: string,
  name: string,
  description?: string,
): Promise<Bank> {
  const id = crypto.randomUUID()
  const now = new Date().toISOString()
  const bank: Bank = { id, owner_id: ownerId, name, description, created_at: now, updated_at: now }
  store.set(id, bank)
  return bank
}

// updateBank replaces the name and description of an existing bank.
// Returns null if the bank doesn't exist or isn't owned by ownerId.
export async function updateBank(
  id: string,
  ownerId: string,
  name: string,
  description?: string,
): Promise<Bank | null> {
  const bank = store.get(id)
  if (!bank || bank.owner_id !== ownerId) return null
  const updated: Bank = { ...bank, name, description, updated_at: new Date().toISOString() }
  store.set(id, updated)
  return updated
}

// deleteBank removes a bank. Returns false if not found or not owned by ownerId.
export async function deleteBank(id: string, ownerId: string): Promise<boolean> {
  const bank = store.get(id)
  if (!bank || bank.owner_id !== ownerId) return false
  store.delete(id)
  return true
}
