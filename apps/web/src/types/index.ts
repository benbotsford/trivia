// Shared TypeScript types for the trivia web app.
// These mirror the JSON shapes returned by the Go API —
// keep them in sync with the bankResponse type in apps/api/internal/game/types.go.

export interface Bank {
  id: string
  owner_id: string
  name: string
  description?: string   // optional — may be absent from the JSON entirely
  created_at: string     // ISO 8601 date string, e.g. "2026-04-19T12:00:00Z"
  updated_at: string
}

export interface Session {
  userId: string
  email: string
  name: string
}
