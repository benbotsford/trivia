-- name: CreateGame :one
INSERT INTO games (id, code, host_id, bank_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGameByID :one
SELECT * FROM games WHERE id = $1;

-- name: GetActiveGameByCode :one
SELECT * FROM games
WHERE code = $1
  AND status IN ('lobby', 'in_progress');

-- name: ListGamesByHost :many
SELECT * FROM games
WHERE host_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: StartGame :one
UPDATE games
SET status     = 'in_progress',
    started_at = now()
WHERE id = $1
  AND status = 'lobby'
RETURNING *;

-- name: AdvanceGameQuestion :one
UPDATE games
SET current_question_idx = current_question_idx + 1
WHERE id = $1
  AND status = 'in_progress'
RETURNING *;

-- name: EndGame :one
UPDATE games
SET status   = 'completed',
    ended_at = now()
WHERE id = $1
  AND status = 'in_progress'
RETURNING *;

-- name: CancelGame :one
UPDATE games
SET status   = 'cancelled',
    ended_at = now()
WHERE id = $1
  AND status IN ('lobby', 'in_progress')
RETURNING *;
