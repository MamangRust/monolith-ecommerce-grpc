-- name: CreateResetToken :one
INSERT INTO
    reset_tokens (user_id, token, expiry_date)
VALUES ($1, $2, $3)
RETURNING
    *;

-- name: DeleteResetToken :exec
DELETE FROM reset_tokens WHERE user_id = $1;

-- name: GetResetToken :one
SELECT * FROM reset_tokens WHERE token = $1;