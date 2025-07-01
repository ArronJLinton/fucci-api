-- +goose Up

CREATE TABLE player_profiles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
  position VARCHAR(50) NOT NULL,
  age INTEGER NOT NULL CHECK (age >= 16 AND age <= 50),
  country VARCHAR(50) NOT NULL,
  height_cm INTEGER NOT NULL, -- TODO: Re-enable CHECK (height_cm >= 150 AND height_cm <= 220) after debugging
  pace INTEGER NOT NULL DEFAULT 50 CHECK (pace >= 0 AND pace <= 100),
  shooting INTEGER NOT NULL DEFAULT 50 CHECK (shooting >= 0 AND shooting <= 100),
  passing INTEGER NOT NULL DEFAULT 50 CHECK (passing >= 0 AND passing <= 100),
  stamina INTEGER NOT NULL DEFAULT 50 CHECK (stamina >= 0 AND stamina <= 100),
  dribbling INTEGER NOT NULL DEFAULT 50 CHECK (dribbling >= 0 AND dribbling <= 100),
  defending INTEGER NOT NULL DEFAULT 50 CHECK (defending >= 0 AND defending <= 100),
  physical INTEGER NOT NULL DEFAULT 50 CHECK (physical >= 0 AND physical <= 100),
  is_verified BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id)
);

-- +goose Down
DROP TABLE player_profiles; 