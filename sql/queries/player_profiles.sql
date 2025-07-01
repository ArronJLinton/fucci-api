-- name: CreatePlayerProfile :one
INSERT INTO player_profiles (
  user_id, team_id, position, age, country, height_cm,
  pace, shooting, passing, stamina, dribbling, defending, physical
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetPlayerProfile :one
SELECT * FROM player_profiles WHERE id = $1;

-- name: GetPlayerProfileByUser :one
SELECT * FROM player_profiles WHERE user_id = $1;

-- name: ListPlayerProfiles :many
SELECT * FROM player_profiles ORDER BY created_at DESC;

-- name: ListPlayerProfilesByTeam :many
SELECT pp.*, u.firstname, u.lastname 
FROM player_profiles pp
JOIN users u ON pp.user_id = u.id
WHERE pp.team_id = $1 
ORDER BY pp.position, u.lastname;

-- name: GetPlayersByTeam :many
SELECT pp.*, u.firstname, u.lastname 
FROM player_profiles pp
JOIN users u ON pp.user_id = u.id
WHERE pp.team_id = $1 
ORDER BY pp.position, u.lastname;

-- name: GetPlayersByLeague :many
SELECT pp.*, u.firstname, u.lastname 
FROM player_profiles pp
JOIN users u ON pp.user_id = u.id
JOIN teams t ON pp.team_id = t.id
WHERE t.league_id = $1 
ORDER BY t.name, pp.position, u.lastname;

-- name: ListFreeAgents :many
SELECT pp.*, u.firstname, u.lastname, u.email 
FROM player_profiles pp
JOIN users u ON pp.user_id = u.id
WHERE pp.team_id IS NULL
ORDER BY pp.created_at DESC;

-- name: UpdatePlayerProfile :one
UPDATE player_profiles 
SET 
  team_id = $2, position = $3, age = $4, country = $5, height_cm = $6,
  pace = $7, shooting = $8, passing = $9, stamina = $10, dribbling = $11, 
  defending = $12, physical = $13, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdatePlayerTeam :one
UPDATE player_profiles 
SET team_id = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: DeletePlayerProfile :exec
DELETE FROM player_profiles WHERE id = $1; 