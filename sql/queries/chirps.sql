-- name: CreateChirps :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetChirps :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
ORDER BY created_at ASC;


-- name: GetChirpsByUserID :one
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE id = $1  
ORDER BY created_at ASC;