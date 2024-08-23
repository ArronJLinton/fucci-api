-- +goose Up

CREATE TABLE teams (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  country VARCHAR(50) NOT NULL,
  state VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE teams; 