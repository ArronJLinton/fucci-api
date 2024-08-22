-- name: CreateUser :one
INSERT INTO users (firstname, lastname, phone)
VALUES ($1, $2, $3)
RETURNING *;