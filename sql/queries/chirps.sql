-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (gen_random_uuid(), now(), now(), $1, $2)
Returning *;

-- name: ResetAllChirps :exec
DELETE FROM chirps;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetAllChirpsDesc :many
SELECT * FROM chirps ORDER BY created_at DESC;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirpByID :exec
DELETE FROM chirps WHERE id = $1;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY $2;