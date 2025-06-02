-- name: CreateMedia :one
INSERT INTO media (match_id, media_url)
VALUES ($1, $2)
RETURNING *;

-- name: GetMediaByMatchId :many
SELECT * FROM media WHERE match_id = $1 ORDER BY created_at DESC;

-- name: DeleteMediaById :exec
DELETE FROM media WHERE id = $1;
