-- +goose Up

CREATE TABLE users (
  id UUID PRIMARY KEY,
  firstname TEXT NOT NULL,
  lastname TEXT NOT NULL,
  photo TEXT,
  phone TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users; 