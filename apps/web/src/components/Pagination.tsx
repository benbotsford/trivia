interface Props {
  page: number
  totalPages: number
  onPage: (page: number) => void
}

export default function Pagination({ page, totalPages, onPage }: Props) {
  if (totalPages <= 1) return null

  return (
    <div className="flex items-center justify-between pt-4">
      <button
        onClick={() => onPage(page - 1)}
        disabled={page === 1}
        className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        ← Prev
      </button>
      <span className="text-xs text-slate-400">
        Page {page} of {totalPages}
      </span>
      <button
        onClick={() => onPage(page + 1)}
        disabled={page === totalPages}
        className="rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
      >
        Next →
      </button>
    </div>
  )
}
