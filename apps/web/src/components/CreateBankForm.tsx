'use client'

import { useRef, useTransition } from 'react'
import type { Bank } from '@/types'
import { createBankAction } from '@/app/banks/actions'

interface CreateBankFormProps {
  onClose: () => void
  // onCreate is called with the newly created bank so the parent can update
  // local state immediately rather than waiting for a page re-fetch.
  onCreate: (bank: Bank) => void
}

export default function CreateBankForm({ onClose, onCreate }: CreateBankFormProps) {
  const formRef = useRef<HTMLFormElement>(null)
  const [isPending, startTransition] = useTransition()

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const formData = new FormData(e.currentTarget)

    startTransition(async () => {
      const result = await createBankAction(formData)
      if (result.bank) {
        onCreate(result.bank) // update parent state instantly
      }
      // onClose is handled by the parent in the onCreate callback in BanksView,
      // but call it here too as a safety net in case onCreate isn't provided.
    })
  }

  return (
    <div className="mb-6 rounded-lg border border-gray-200 bg-white p-5 shadow-sm">
      <h2 className="mb-4 text-base font-semibold text-gray-900">New Question Bank</h2>

      <form ref={formRef} onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="name" className="mb-1 block text-sm font-medium text-gray-700">
            Name <span className="text-brand-red">*</span>
          </label>
          <input
            id="name"
            name="name"
            type="text"
            required
            autoFocus
            placeholder="e.g. General Knowledge"
            className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm placeholder-gray-400 focus:border-brand-blue focus:outline-none focus:ring-1 focus:ring-brand-blue"
          />
        </div>

        <div>
          <label htmlFor="description" className="mb-1 block text-sm font-medium text-gray-700">
            Description{' '}
            <span className="font-normal text-gray-400">(optional)</span>
          </label>
          <input
            id="description"
            name="description"
            type="text"
            placeholder="What's this bank about?"
            className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm placeholder-gray-400 focus:border-brand-blue focus:outline-none focus:ring-1 focus:ring-brand-blue"
          />
        </div>

        <div className="flex justify-end gap-2 pt-1">
          <button
            type="button"
            onClick={onClose}
            disabled={isPending}
            className="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isPending}
            className="rounded-md bg-brand-blue px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
          >
            {isPending ? 'Creating…' : 'Create Bank'}
          </button>
        </div>
      </form>
    </div>
  )
}
