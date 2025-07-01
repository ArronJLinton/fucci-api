-- name: CreateVerification :one
INSERT INTO verifications (player_profile_id, verifier_user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetVerification :one
SELECT * FROM verifications WHERE id = $1;

-- name: ListVerificationsByPlayer :many
SELECT * FROM verifications WHERE player_profile_id = $1;

-- name: CountVerificationsByPlayer :one
SELECT COUNT(*) FROM verifications WHERE player_profile_id = $1;

-- name: DeleteVerification :exec
DELETE FROM verifications WHERE id = $1;

-- name: UpdatePlayerVerificationStatus :exec
UPDATE player_profiles 
SET is_verified = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1; 