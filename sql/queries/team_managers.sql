-- name: CreateTeamManager :one
INSERT INTO team_managers (user_id, league_id, team_id, title, experience, bio)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetTeamManager :one
SELECT * FROM team_managers WHERE id = $1;

-- name: ListTeamManagers :many
SELECT * FROM team_managers 
WHERE ($1::uuid IS NULL OR league_id = $1)
  AND ($2::uuid IS NULL OR team_id = $2)
  AND ($3::uuid IS NULL OR user_id = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: UpdateTeamManager :one
UPDATE team_managers 
SET team_id = $2, title = $3, experience = $4, bio = $5, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeleteTeamManager :exec
DELETE FROM team_managers WHERE id = $1;

-- name: GetTeamManagersByLeague :many
SELECT * FROM team_managers WHERE league_id = $1 ORDER BY created_at DESC;

-- name: ListTeamsByManager :many
SELECT t.* FROM teams t
JOIN team_managers tm ON t.id = tm.team_id
WHERE tm.user_id = $1; 