-- name: CreateTeam :one
INSERT INTO teams (name, country, state)
VALUES ($1, $2, $3)
RETURNING *;