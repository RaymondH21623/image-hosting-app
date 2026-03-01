-- name: CreateUser :one
INSERT INTO users (public_id, username, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT id, username, email, created_at FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByEmailAuth :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, username, email, created_at FROM users ORDER BY created_at DESC;
