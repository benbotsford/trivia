// API client for games. Server-side only.

import type { Game, GamePlayer, GameResults } from '@/types'
import { apiFetch } from './client'

// listGames returns all games created by the authenticated host, newest first.
export async function listGames(): Promise<Game[]> {
  const res = await apiFetch('/games')
  return res.json()
}

// createGame creates a new game linked to a quiz (preferred) or bank (legacy).
export async function createGame(opts: { quizID?: string; bankID?: string; roundSize?: number }): Promise<Game> {
  const body: Record<string, unknown> = {}
  if (opts.quizID) body.quiz_id = opts.quizID
  if (opts.bankID) { body.bank_id = opts.bankID; body.round_size = opts.roundSize ?? 5 }
  const res = await apiFetch('/games', {
    method: 'POST',
    body: JSON.stringify(body),
  })
  return res.json()
}

// getGame fetches a single game by ID.
export async function getGame(id: string): Promise<Game> {
  const res = await apiFetch(`/games/${id}`)
  return res.json()
}

// listPlayers returns current players in a game (for the lobby UI).
export async function listPlayers(gameID: string): Promise<GamePlayer[]> {
  const res = await apiFetch(`/games/${gameID}/players`)
  return res.json()
}

// cancelGame transitions a lobby or in-progress game to 'cancelled'.
export async function cancelGame(gameID: string): Promise<void> {
  await apiFetch(`/games/${gameID}`, { method: 'DELETE' })
}

// getGameResults returns per-player answer history and final scores for a completed game.
export async function getGameResults(gameID: string): Promise<GameResults> {
  const res = await apiFetch(`/games/${gameID}/results`)
  return res.json()
}
