-- name: CreateUser :one
INSERT INTO users (id, auth0_sub, email, display_name)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByAuth0Sub :one
SELECT * FROM users WHERE auth0_sub = $1;

-- name: UpsertUserByAuth0Sub :one
INSERT INTO users (id, auth0_sub, email, display_name)
VALUES ($1, $2, $3, $4)
ON CONFLICT (auth0_sub) DO UPDATE SET
    email        = EXCLUDED.email,
    display_name = EXCLUDED.display_name
RETURNING *;

-- name: UpdateUserProfile :one
UPDATE users
SET email        = $2,
    display_name = $3
WHERE id = $1
RETURNING *;
