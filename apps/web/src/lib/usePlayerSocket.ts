'use client'

import { useCallback } from 'react'
import { useGameSocket, SocketStatus } from './useGameSocket'
import type {
  GameStartedPayload,
  QuestionReleasedPayload,
  AnswerAcceptedPayload,
  RoundEndedPayload,
  RoundScoresPayload,
  LeaderboardPayload,
} from '@/types'

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export interface PlayerSocketHandlers {
  onGameStarted?: (payload: GameStartedPayload) => void
  onQuestionReleased?: (payload: QuestionReleasedPayload) => void
  onAnswerAccepted?: (payload: AnswerAcceptedPayload) => void
  onRoundEnded?: (payload: RoundEndedPayload) => void
  onRoundScores?: (payload: RoundScoresPayload) => void
  onRoundLeaderboard?: (payload: LeaderboardPayload) => void
  onGameEnded?: (payload: LeaderboardPayload) => void
  onOpen?: () => void
}

export interface UsePlayerSocketResult {
  status: SocketStatus

  // Player action — no-ops silently if the socket is not open.
  submitAnswer: (questionID: string, answer: string) => void
}

// ---------------------------------------------------------------------------
// Hook
// ---------------------------------------------------------------------------

/**
 * Typed WebSocket hook for the player game panel.
 *
 * Wraps useGameSocket with player-specific message dispatch and exposes
 * submitAnswer so PlayerGame.tsx never constructs raw message objects.
 *
 * Session token is read from sessionStorage by the caller and passed in —
 * this hook does not touch sessionStorage directly, keeping it testable.
 *
 * @param wsBase        WebSocket base URL, e.g. "ws://localhost:8080"
 * @param code          6-character game code
 * @param sessionToken  Player session token (from sessionStorage after joining)
 * @param handlers      Message handler callbacks — all optional
 */
export function usePlayerSocket(
  wsBase: string,
  code: string,
  sessionToken: string,
  handlers: PlayerSocketHandlers,
): UsePlayerSocketResult {
  const url = `${wsBase}/ws/${code}?session=${encodeURIComponent(sessionToken)}`

  const onMessage = useCallback((type: string, payload: unknown) => {
    switch (type) {
      case 'game_started':
        handlers.onGameStarted?.(payload as GameStartedPayload)
        break
      case 'question_released':
        handlers.onQuestionReleased?.(payload as QuestionReleasedPayload)
        break
      case 'answer_accepted':
        handlers.onAnswerAccepted?.(payload as AnswerAcceptedPayload)
        break
      case 'round_ended':
        handlers.onRoundEnded?.(payload as RoundEndedPayload)
        break
      case 'round_scores':
        handlers.onRoundScores?.(payload as RoundScoresPayload)
        break
      case 'round_leaderboard':
        handlers.onRoundLeaderboard?.(payload as LeaderboardPayload)
        break
      case 'game_ended':
        handlers.onGameEnded?.(payload as LeaderboardPayload)
        break
    }
  // handlers intentionally excluded — useGameSocket keeps onMessage current
  // via optionsRef, so we don't need a new callback when handler identities change.
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const { status, send } = useGameSocket(url, {
    onMessage,
    onOpen: handlers.onOpen,
  })

  // ---- Actions -------------------------------------------------------------

  const submitAnswer = useCallback(
    (questionID: string, answer: string) =>
      send('submit_answer', { question_id: questionID, answer }),
    [send],
  )

  return { status, submitAnswer }
}
