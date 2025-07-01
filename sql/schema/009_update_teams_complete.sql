-- +goose Up

-- Add missing fields to teams table
ALTER TABLE teams 
ADD COLUMN IF NOT EXISTS description TEXT,
ADD COLUMN IF NOT EXISTS manager_id UUID REFERENCES team_managers(id),
ADD COLUMN IF NOT EXISTS logo_url VARCHAR(255),
ADD COLUMN IF NOT EXISTS city VARCHAR(100),
ADD COLUMN IF NOT EXISTS founded INTEGER,
ADD COLUMN IF NOT EXISTS stadium VARCHAR(100),
ADD COLUMN IF NOT EXISTS capacity INTEGER;

-- Remove the height_cm check constraint from player_profiles
ALTER TABLE player_profiles DROP CONSTRAINT IF EXISTS player_profiles_height_cm_check; -- TODO: Re-enable this constraint after debugging

-- +goose Down

-- Remove the added columns
ALTER TABLE teams 
DROP COLUMN IF EXISTS description,
DROP COLUMN IF EXISTS manager_id,
DROP COLUMN IF EXISTS logo_url,
DROP COLUMN IF EXISTS city,
DROP COLUMN IF EXISTS founded,
DROP COLUMN IF EXISTS stadium,
DROP COLUMN IF EXISTS capacity;

-- Add the height_cm check constraint back to player_profiles
ALTER TABLE player_profiles ADD CONSTRAINT IF NOT EXISTS player_profiles_height_cm_check CHECK (height_cm >= 150 AND height_cm <= 220); 