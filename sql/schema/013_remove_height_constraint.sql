-- +goose Up
-- Remove the height_cm check constraint from player_profiles
ALTER TABLE player_profiles DROP CONSTRAINT IF EXISTS player_profiles_height_cm_check;

-- +goose Down
-- Add the height_cm check constraint back to player_profiles
ALTER TABLE player_profiles ADD CONSTRAINT IF NOT EXISTS player_profiles_height_cm_check CHECK (height_cm >= 150 AND height_cm <= 220); 