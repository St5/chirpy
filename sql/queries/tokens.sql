-- name: CreateToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT u.* FROM refresh_tokens as rt
JOIN users as u ON rt.user_id = u.id
WHERE token = $1 AND expires_at > now() AND revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked_at = now()
WHERE token = $1;
