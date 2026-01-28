-- name: CreateMedia :one
INSERT INTO media (public_media_id, user_id, filename, mime_type, size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListMediaByUser :many
SELECT * FROM media WHERE user_id = $1 ORDER BY created_at DESC;

-- name: GetMediaByID :one
SELECT * FROM media WHERE public_media_id = $1;