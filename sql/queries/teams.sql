-- name: CreateTeam :one
INSERT INTO teams (name, description, league_id, manager_id, logo_url, city, country, founded, stadium, capacity)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetTeam :one
SELECT * FROM teams WHERE id = $1;

-- name: ListTeams :many
SELECT * FROM teams 
WHERE ($1::uuid IS NULL OR league_id = $1)
  AND ($2::uuid IS NULL OR manager_id = $2)
ORDER BY name
LIMIT $3 OFFSET $4;

-- name: ListTeamsByLeague :many
SELECT * FROM teams WHERE league_id = $1 ORDER BY name;

-- name: UpdateTeam :one
UPDATE teams 
SET name = $2, description = $3, league_id = $4, manager_id = $5, 
    logo_url = $6, city = $7, country = $8, founded = $9, 
    stadium = $10, capacity = $11, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = $1;

-- name: GetTeamsByLeague :many
SELECT * FROM teams WHERE league_id = $1 ORDER BY name;