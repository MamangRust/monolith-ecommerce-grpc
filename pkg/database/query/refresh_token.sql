-- CreateRefreshToken: Creates a new refresh token
-- Purpose: Generate a refresh token for user authentication
-- Parameters:
--   $1: user_id - ID of the user this token belongs to
--   $2: token - The actual refresh token string
--   $3: expiration - Expiration timestamp of the token
-- Returns: The created refresh token record (excluding sensitive fields if any)
-- Business Logic:
--   - Sets both created_at and updated_at to current timestamp
--   - Used in JWT refresh token rotation
--   - Typically created during login/auth flows
-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expiration, created_at, updated_at)
VALUES ($1, $2, $3, current_timestamp, current_timestamp)
RETURNING refresh_token_id, user_id, token, expiration, created_at, updated_at, deleted_at;

-- FindRefreshTokenByToken: Retrieves active refresh token by token string
-- Purpose: Validate and lookup refresh token
-- Parameters:
--   $1: token - The refresh token string to find
-- Returns: The refresh token record if found and active
-- Business Logic:
--   - Only returns non-deleted tokens
--   - Used during token refresh operations
--   - Helps prevent token reuse
-- name: FindRefreshTokenByToken :one
SELECT refresh_token_id, user_id, token, expiration, created_at, updated_at, deleted_at
FROM refresh_tokens
WHERE token = $1 AND deleted_at IS NULL;

-- FindRefreshTokenByUserId: Retrieves latest active refresh token for user
-- Purpose: Get current valid refresh token for a user
-- Parameters:
--   $1: user_id - ID of the user to find token for
-- Returns: The most recent refresh token for the user
-- Business Logic:
--   - Returns only active (non-deleted) tokens
--   - Orders by creation date (newest first)
--   - Used for token management and validation
--   - Limits to 1 result to get latest token
-- name: FindRefreshTokenByUserId :one
SELECT
    refresh_token_id,
    user_id,
    token,
    expiration,
    created_at,
    updated_at,
    deleted_at
FROM
    refresh_tokens
WHERE
    user_id = $1 AND deleted_at IS NULL
ORDER BY
    created_at DESC
LIMIT 1;

-- UpdateRefreshTokenByUserId: Updates refresh token for a user
-- Purpose: Rotate/refresh token for a user
-- Parameters:
--   $1: user_id - ID of the user to update token for
--   $2: token - New token string
--   $3: expiration - New expiration timestamp
-- Business Logic:
--   - Updates token and expiration fields
--   - Sets updated_at to current time
--   - Only modifies active tokens
--   - Used during token rotation flows
-- name: UpdateRefreshTokenByUserId :one
UPDATE refresh_tokens
SET token = $2, expiration = $3, updated_at = current_timestamp
WHERE user_id = $1 AND deleted_at IS NULL
RETURNING *;

-- DeleteRefreshToken: Permanently deletes a refresh token
-- Purpose: Invalidate a specific refresh token
-- Parameters:
--   $1: token - The token string to delete
-- Business Logic:
--   - Hard deletes the token record
--   - Used during logout/token invalidation
--   - Prevents token reuse after deletion
-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens
WHERE token = $1;

-- DeleteRefreshTokenByUserId: Permanently deletes all tokens for a user
-- Purpose: Invalidate all refresh tokens for a user
-- Parameters:
--   $1: user_id - ID of the user whose tokens to delete
-- Business Logic:
--   - Hard deletes all tokens for the user
--   - Used during password reset or account lock
--   - Ensures complete session invalidation
-- name: DeleteRefreshTokenByUserId :exec
DELETE FROM refresh_tokens
WHERE user_id = $1;
