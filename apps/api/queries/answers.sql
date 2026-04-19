-- name: RecordAnswer :one
INSERT INTO answers (id, game_id, question_id, player_id, answer, is_correct, points_awarded)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetAnswer :one
SELECT * FROM answers
WHERE game_id = $1
  AND question_id = $2
  AND player_id = $3;

-- name: ListAnswersForGame :many
SELECT * FROM answers
WHERE game_id = $1
ORDER BY submitted_at ASC;

-- name: ListAnswersForQuestion :many
SELECT * FROM answers
WHERE game_id = $1
  AND question_id = $2
ORDER BY submitted_at ASC;

-- name: ListAnswersForPlayer :many
SELECT * FROM answers
WHERE player_id = $1
ORDER BY submitted_at ASC;
