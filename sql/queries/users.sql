-- name: CreateUser :one
INSERT INTO users (firstname, lastname, email)
VALUES ($1, $2, $3)
RETURNING *;