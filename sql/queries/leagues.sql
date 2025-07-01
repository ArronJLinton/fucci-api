-- name: CreateLeague :one
INSERT INTO leagues (name, description, owner_id, country, level, logo_url, website, founded)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetLeague :one
SELECT * FROM leagues WHERE id = $1;

-- name: ListLeagues :many
SELECT * FROM leagues 
WHERE ($1::int IS NULL OR owner_id = $1)
  AND ($2::text IS NULL OR country = $2)
  AND ($3::int IS NULL OR level = $3)
ORDER BY name
LIMIT $4 OFFSET $5;

-- name: UpdateLeague :one
UPDATE leagues 
SET name = $2, description = $3, country = $4, level = $5, 
    logo_url = $6, website = $7, founded = $8, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteLeague :exec
DELETE FROM leagues WHERE id = $1; 