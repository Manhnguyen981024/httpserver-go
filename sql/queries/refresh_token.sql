-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at, created_at, updated_at, revoked_at)
VALUES ($1, $2, $3, NOW(), NOW(), NULL)
RETURNING user_id, token, expires_at, created_at, updated_at, revoked_at;

-- name: GetRefreshTokenByUserId :one
SELECT user_id, token, expires_at, created_at, updated_at, revoked_at
FROM refresh_tokens
WHERE user_id = $1
    AND revoked_at IS NULL
LIMIT 1;

-- name: GetRefreshTokenByToken :one
SELECT user_id, token, expires_at, created_at, updated_at, revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;