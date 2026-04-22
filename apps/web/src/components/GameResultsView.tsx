'use client'

import { Fragment, useState } from 'react'
import type { GameResults, GamePlayerResult } from '@/types'

interface Props {
  results: GameResults
}

function correctCount(player: GamePlayerResult): number {
  return player.answers.filter((a) => a.is_correct).length
}

export default function GameResultsView({ results }: Props) {
  const [expanded, setExpanded] = useState<string | null>(null)

  if (results.players.length === 0) {
    return (
      <div className="overflow-hidden rounded-xl border border-dashed border-gray-200 bg-white/60 px-8 py-12 text-center text-slate-400">
        No answers recorded for this game.
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {/* Leaderboard table */}
      <div className="overflow-hidden rounded-xl border border-gray-100 bg-white shadow-sm">
        <div className="h-[3px] bg-brand-blue" />
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-gray-100 text-left text-xs font-medium uppercase tracking-wide text-slate-400">
              <th className="px-5 py-3 w-10">#</th>
              <th className="px-5 py-3">Player</th>
              <th className="px-5 py-3 text-right">Correct</th>
              <th className="px-5 py-3 text-right">Score</th>
              <th className="px-5 py-3 w-10" />
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-50">
            {results.players.map((player, idx) => {
              const isExpanded = expanded === player.player_id
              const correct = correctCount(player)
              const total = player.answers.length

              return (
                <Fragment key={player.player_id}>
                  <tr
                    className="hover:bg-slate-50 cursor-pointer transition-colors"
                    onClick={() =>
                      setExpanded(isExpanded ? null : player.player_id)
                    }
                  >
                    {/* Rank */}
                    <td className="px-5 py-3.5 text-slate-400 font-medium">
                      {idx === 0 ? '🥇' : idx === 1 ? '🥈' : idx === 2 ? '🥉' : idx + 1}
                    </td>

                    {/* Name */}
                    <td className="px-5 py-3.5 font-medium text-slate-700">
                      {player.display_name}
                    </td>

                    {/* Correct count */}
                    <td className="px-5 py-3.5 text-right tabular-nums text-slate-500">
                      {total > 0 ? `${correct} / ${total}` : '—'}
                    </td>

                    {/* Score */}
                    <td className="px-5 py-3.5 text-right tabular-nums font-semibold text-slate-800">
                      {player.total_score.toLocaleString()}
                    </td>

                    {/* Expand toggle */}
                    <td className="px-5 py-3.5 text-slate-300 text-xs">
                      {isExpanded ? '▲' : '▼'}
                    </td>
                  </tr>

                  {/* Expanded answer detail */}
                  {isExpanded && (
                    <tr>
                      <td colSpan={5} className="px-5 pb-4 pt-0">
                        <div className="rounded-lg border border-gray-100 bg-slate-50 overflow-hidden">
                          {player.answers.length === 0 ? (
                            <p className="px-4 py-3 text-xs text-slate-400">
                              No answers recorded.
                            </p>
                          ) : (
                            <table className="w-full text-xs">
                              <thead>
                                <tr className="border-b border-gray-100 text-left text-slate-400 uppercase tracking-wide">
                                  <th className="px-4 py-2">Answer</th>
                                  <th className="px-4 py-2 text-right">Points</th>
                                  <th className="px-4 py-2 w-20 text-center">Result</th>
                                </tr>
                              </thead>
                              <tbody className="divide-y divide-gray-100">
                                {player.answers.map((ans) => (
                                  <tr key={ans.question_id}>
                                    <td className="px-4 py-2 text-slate-600">
                                      {ans.answer || <span className="italic text-slate-300">no answer</span>}
                                    </td>
                                    <td className="px-4 py-2 text-right tabular-nums text-slate-500">
                                      {ans.is_correct ? `+${ans.points_awarded.toLocaleString()}` : '—'}
                                    </td>
                                    <td className="px-4 py-2 text-center">
                                      {ans.is_correct ? (
                                        <span className="rounded-full bg-green-100 px-2 py-0.5 text-green-700">
                                          ✓
                                        </span>
                                      ) : (
                                        <span className="rounded-full bg-red-50 px-2 py-0.5 text-red-400">
                                          ✗
                                        </span>
                                      )}
                                    </td>
                                  </tr>
                                ))}
                              </tbody>
                            </table>
                          )}
                        </div>
                      </td>
                    </tr>
                  )}
                </Fragment>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
