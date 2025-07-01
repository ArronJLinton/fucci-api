-- +goose Up

CREATE TABLE verifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  player_profile_id UUID NOT NULL REFERENCES player_profiles(id) ON DELETE CASCADE,
  verifier_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(player_profile_id, verifier_user_id)
);

-- +goose Down
DROP TABLE verifications; 