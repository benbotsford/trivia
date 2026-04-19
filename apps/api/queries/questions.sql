-- name: CreateQuestion :one
INSERT INTO questions (id, bank_id, type, prompt, correct_answer, choices, points, position)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetQuestion :one
SELECT * FROM questions WHERE id = $1;

-- name: ListQuestionsByBank :many
SELECT * FROM questions
WHERE bank_id = $1
ORDER BY position ASC, created_at ASC;

-- name: CountQuestionsInBank :one
SELECT COUNT(*)::int AS count FROM questions WHERE bank_id = $1;

-- name: UpdateQuestion :one
UPDATE questions
SET type           = $2,
    prompt         = $3,
    correct_answer = $4,
    choices        = $5,
    points         = $6,
    position       = $7
WHERE id = $1
RETURNING *;

-- name: ReorderQuestion :one
UPDATE questions
SET position = $2
WHERE id = $1
RETURNING *;

-- name: DeleteQuestion :exec
DELETE FROM questions WHERE id = $1;
