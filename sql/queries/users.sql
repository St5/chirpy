-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
Returning *;

-- name: ResetAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET email = $1, hashed_password = $2, updated_at = now()
WHERE id = $3
RETURNING *;

-- name: UpdateChirpyRedByUserID :exec
UPDATE users SET is_chirpy_red = $1, updated_at = now()
WHERE id = $2;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;