// Game results page — server component.
// Fetches per-player answer history for a completed game.
import { notFound } from 'next/navigation'
import Link from 'next/link'
import { getGame, getGameResults } from '@/lib/api/games'
import GameResultsView from '@/components/GameResultsView'

interface Props {
  params: Promise<{ gameID: string }>
}

export async function generateMetadata({ params }: Props) {
  const { gameID } = await params
  try {
    const game = await getGame(gameID)
    return { title: `Results — ${game.code} — Quibble` }
  } catch {
    return { title: 'Game Results — Quibble' }
  }
}

export default async function GameResultsPage({ params }: Props) {
  const { gameID } = await params

  let results
  try {
    results = await getGameResults(gameID)
  } catch {
    notFound()
  }

  return (
    <main className="mx-auto max-w-4xl px-6 py-10">
      <div className="mb-8 flex items-center gap-4">
        <Link
          href="/games"
          className="text-sm text-slate-400 hover:text-slate-600 transition-colors"
        >
          ← Back to Games
        </Link>
        <h1 className="text-2xl font-semibold text-slate-800">
          Results —{' '}
          <span className="font-display tracking-widest text-brand-blue">
            {results.code}
          </span>
        </h1>
      </div>

      <GameResultsView results={results} />
    </main>
  )
}
