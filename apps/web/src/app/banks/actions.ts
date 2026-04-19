'use server'

// Server actions for question bank mutations.
//
// Each action now returns the mutated bank (or a success flag) so the calling
// client component can update its local state immediately — giving instant UI
// feedback without needing a full page re-fetch.
//
// revalidatePath is still called so that if the user navigates away and back,
// or opens the page in a new tab, they see fresh data from the server.

import { revalidatePath } from 'next/cache'
import { getSession } from '@/lib/session'
import * as banksApi from '@/lib/mock/banks'
import type { Bank } from '@/types'

export async function createBankAction(
  formData: FormData,
): Promise<{ bank?: Bank; error?: string }> {
  const session = getSession()
  const name = (formData.get('name') as string | null)?.trim()
  const description = (formData.get('description') as string | null)?.trim() || undefined

  if (!name) return { error: 'Name is required' }

  const bank = await banksApi.createBank(session.userId, name, description)
  revalidatePath('/banks')
  // Return the new bank so the client can append it to local state instantly.
  return { bank }
}

export async function updateBankAction(
  bankId: string,
  formData: FormData,
): Promise<{ bank?: Bank; error?: string }> {
  const session = getSession()
  const name = (formData.get('name') as string | null)?.trim()
  const description = (formData.get('description') as string | null)?.trim() || undefined

  if (!name) return { error: 'Name is required' }

  const bank = await banksApi.updateBank(bankId, session.userId, name, description)
  if (!bank) return { error: 'Bank not found' }

  revalidatePath('/banks')
  // Return the updated bank so the client can replace the stale card in local state.
  return { bank }
}

export async function deleteBankAction(
  bankId: string,
): Promise<{ success: boolean }> {
  const session = getSession()
  const success = await banksApi.deleteBank(bankId, session.userId)
  revalidatePath('/banks')
  return { success }
}
