-- name: CreateMedia :one
INSERT INTO media (public_media_id, user_id, filename, mime_type, size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListMediaByUser :many
SELECT media.* FROM media JOIN users ON media.user_id = users.id
WHERE users.public_id = $1;

-- name: GetMediaByID :one
SELECT * FROM media WHERE public_media_id = $1;

-- name: GetMediaNameByPublicID :one
SELECT filename FROM media WHERE public_media_id = $1;