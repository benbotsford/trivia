// Shared TypeScript types for the trivia web app.
// These mirror the JSON shapes returned by the Go API —
// keep them in sync with the response types in apps/api/internal/game/types.go.

export interface Bank {
  id: string
  owner_id: string
  name: string
  description?: string
  created_at: string
  updated_at: string
}

export interface Session {
  userId: string
  email: string
  name: string
}

export interface UserProfile {
  id: string
  email?: string
  display_name?: string
  created_at: string
}

export interface MCChoice {
  text: string
  correct: boolean
}

export interface Question {
  id: string
  bank_id: string
  type: 'text' | 'multiple_choice'
  prompt: string
  points: number
  position: number
  accepted_answers?: string[]
  choices?: MCChoice[]
  created_at: string
  updated_at: string
}

export const MAX_PROMPT_LEN  = 500
export const MAX_CHOICE_LEN  = 200
export const MAX_ANSWER_LEN  = 150
export const MAX_CHOICES     = 6
export const MIN_CHOICES     = 2
export const MAX_ANSWERS     = 10

// --- Quiz types ---

export interface Quiz {
  id: string
  owner_id: string
  name: string
  description?: string
  created_at: string
  updated_at: string
}

export interface QuizRound {
  id: string
  quiz_id: string
  round_number: number
  title?: string
  questions: Question[]
  created_at: string
}

export interface QuizDetail extends Quiz {
  rounds: QuizRound[]
}

// --- Game types ---

export type GameStatus = 'lobby' | 'in_progress' | 'completed' | 'cancelled'

export interface Game {
  id: string
  code: string
  status: GameStatus
  bank_id?: string
  quiz_id?: string
  current_question_idx: number
  current_round_idx: number
  round_size: number
  created_at: string
}

export interface GamePlayer {
  id: string
  display_name: string
  score: number
}

export interface GameAnswerResult {
  question_id: string
  answer: string
  is_correct: boolean
  points_awarded: number
  submitted_at: string
}

export interface GamePlayerResult {
  player_id: string
  display_name: string
  total_score: number
  answers: GameAnswerResult[]
}

export interface GameResults {
  game_id: string
  code: string
  players: GamePlayerResult[]
}

export interface JoinGameResponse {
  game_code: string
  session_token: string
  display_name: string
}

// --- WebSocket message shapes ---

export interface LeaderboardEntry {
  rank: number
  display_name: string
  score: number
}

export interface QuestionPayload {
  id: string
  type: 'text' | 'multiple_choice'
  prompt: string
  points: number
  choices?: string[]
}

// question_released (new, quiz-based)
export interface QuestionReleasedPayload {
  pos_in_round: number   // 1-indexed position of this question within the round
  total_in_round: number // total questions in this round
  round: number          // current round number (1-indexed)
  total_rounds: number
  question: QuestionPayload
}

// round_ended — sent to all when host moves to review phase
export interface RoundEndedPayload {
  round: number
  total_rounds: number
}

// round_review — sent to host only
export interface RoundReviewAnswerEntry {
  player_id: string
  player_name: string
  answer: string
  correct: boolean
  overridden: boolean
}

export interface RoundReviewQuestion {
  question_id: string
  prompt: string
  correct_answers: string[]
  answers: RoundReviewAnswerEntry[]
}

export interface RoundReviewPayload {
  round: number
  total_rounds: number
  questions: RoundReviewQuestion[]
}

// round_scores — sent per-player (and a host variant)
export interface RoundScoreQuestion {
  question_id: string
  prompt: string
  correct_answers: string[]
  your_answer?: string    // absent on host variant
  correct?: boolean       // absent on host variant
  points_earned?: number  // absent on host variant
}

export interface RoundScoresPayload {
  round: number
  total_rounds: number
  questions: RoundScoreQuestion[]
  round_score?: number  // player total for this round
  is_host?: boolean
}

// round_leaderboard / game_ended
export interface LeaderboardPayload {
  entries: LeaderboardEntry[]
  round: number
  total_rounds: number
}

// answer_accepted — sent per-player
export interface AnswerAcceptedPayload {
  question_id: string
  correct: boolean
}

// scoreboard_update — sent to host
export interface ScoreboardUpdatePayload {
  question_id: string
  answer_count: number
}

// lobby_update
export interface LobbyPlayer {
  id: string
  display_name: string
}

export interface LobbyUpdatePayload {
  players: LobbyPlayer[]
}

// game_started
export interface GameStartedPayload {
  round?: number
  total_rounds?: number
  // legacy fields (bank-based)
  total?: number
  round_size?: number
}

// ---------------------------------------------------------------------------
// Message type constants
// ---------------------------------------------------------------------------
// Single source of truth for the string literals used on the wire.
// Import these instead of sprinkling raw strings across components and hooks.

// Server → client
export const MSG_LOBBY_UPDATE       = 'lobby_update'       as const
export const MSG_GAME_STARTED       = 'game_started'       as const
export const MSG_QUESTION_RELEASED  = 'question_released'  as const
export const MSG_SCOREBOARD_UPDATE  = 'scoreboard_update'  as const
export const MSG_ROUND_ENDED        = 'round_ended'        as const
export const MSG_ROUND_REVIEW       = 'round_review'       as const
export const MSG_ROUND_SCORES       = 'round_scores'       as const
export const MSG_ROUND_LEADERBOARD  = 'round_leaderboard'  as const
export const MSG_GAME_ENDED         = 'game_ended'         as const
export const MSG_ANSWER_ACCEPTED    = 'answer_accepted'    as const

// Client → server (host)
export const MSG_START_GAME         = 'start_game'         as const
export const MSG_RELEASE_QUESTION   = 'release_question'   as const
export const MSG_END_ROUND          = 'end_round'          as const
export const MSG_OVERRIDE_ANSWER    = 'override_answer'    as const
export const MSG_RELEASE_SCORES     = 'release_scores'     as const
export const MSG_START_NEXT_ROUND   = 'start_next_round'   as const
export const MSG_END_GAME           = 'end_game'           as const

// Client → server (player)
export const MSG_SUBMIT_ANSWER      = 'submit_answer'      as const

// ---------------------------------------------------------------------------
// Inbound message discriminated union
// ---------------------------------------------------------------------------
// Exhaustive union of every message the server can send. Useful for typed
// switch statements and ensures the switch in useGameSocket's onMessage is
// checked at compile time.

export type InboundMessage =
  | { type: typeof MSG_LOBBY_UPDATE;      payload: LobbyUpdatePayload }
  | { type: typeof MSG_GAME_STARTED;      payload: GameStartedPayload }
  | { type: typeof MSG_QUESTION_RELEASED; payload: QuestionReleasedPayload }
  | { type: typeof MSG_SCOREBOARD_UPDATE; payload: ScoreboardUpdatePayload }
  | { type: typeof MSG_ROUND_ENDED;       payload: RoundEndedPayload }
  | { type: typeof MSG_ROUND_REVIEW;      payload: RoundReviewPayload }
  | { type: typeof MSG_ROUND_SCORES;      payload: RoundScoresPayload }
  | { type: typeof MSG_ROUND_LEADERBOARD; payload: LeaderboardPayload }
  | { type: typeof MSG_GAME_ENDED;        payload: LeaderboardPayload }
  | { type: typeof MSG_ANSWER_ACCEPTED;   payload: AnswerAcceptedPayload }
