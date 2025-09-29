-- name: CreateMedia :one
INSERT INTO media (user_id, slug, filename, mime_type, size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMediaBySlug :one
SELECT * FROM media WHERE slug = $1 LIMIT 1;

-- name: ListMediaByUser :many
SELECT * FROM media WHERE user_id = $1 ORDER BY created_at DESC;
