// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (firstname, lastname, email)
VALUES ($1, $2, $3)
RETURNING id, firstname, lastname, email, created_at, updated_at
`

type CreateUserParams struct {
	Firstname string
	Lastname  string
	Email     string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Firstname, arg.Lastname, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Lastname,
		&i.Email,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}