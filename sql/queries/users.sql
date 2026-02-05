-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, is_chirpy_red, hashed_password
FROM users
WHERE email = $1;

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2, email = $3, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;


-- name: GetUserByID :one
SELECT id, created_at, updated_at, email, is_chirpy_red
FROM users
WHERE id = $1;

-- name: UpdateUserChirpyRed :one
UPDATE users
SET is_chirpy_red = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, email, is_chirpy_red;