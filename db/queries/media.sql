-- name: CreateMedia :one
INSERT INTO media (user_id, filename, mime_type, size)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListMediaByUser :many
SELECT * FROM media WHERE user_id = $1 ORDER BY created_at DESC;
