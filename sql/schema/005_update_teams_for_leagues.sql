-- +goose Up

-- Drop existing teams table and recreate with new structure
DROP TABLE IF EXISTS teams CASCADE;

CREATE TABLE teams (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) NOT NULL,
  league_id UUID REFERENCES leagues(id),
  state VARCHAR(50),
  country VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down

-- Recreate the original teams table structure
DROP TABLE IF EXISTS teams CASCADE;

CREATE TABLE teams (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  country VARCHAR(50) NOT NULL,
  state VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
); 