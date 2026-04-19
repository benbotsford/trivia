-- name: CreateQuestionBank :one
INSERT INTO question_banks (id, owner_id, name, description)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetQuestionBank :one
SELECT * FROM question_banks WHERE id = $1;

-- name: ListQuestionBanksByOwner :many
SELECT * FROM question_banks
WHERE owner_id = $1
ORDER BY created_at DESC;

-- name: UpdateQuestionBank :one
UPDATE question_banks
SET name        = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteQuestionBank :exec
DELETE FROM question_banks WHERE id = $1;
