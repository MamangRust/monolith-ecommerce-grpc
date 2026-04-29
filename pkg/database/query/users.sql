-- GetUsers: Retrieves paginated list of active users with search capability
-- Purpose: List all active users for management UI
-- Parameters:
--   $1: search_term - Optional text to filter users by name or email (NULL for no filter)
--   $2: limit - Maximum number of records to return (pagination limit)
--   $3: offset - Number of records to skip (pagination offset)
-- Returns:
--   All user fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted users (deleted_at IS NULL)
--   - Supports partial text matching on firstname, lastname, and email fields (case-insensitive)
--   - Returns newest users first (created_at DESC)
--   - Provides total_count for client-side pagination
--   - Uses window function COUNT(*) OVER() for efficient total count
-- name: GetUsers :many
SELECT
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM users
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR firstname ILIKE '%' || $1 || '%'
        OR lastname ILIKE '%' || $1 || '%'
        OR email ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetUsersActive: Retrieves paginated list of active users (identical to GetUsers)
-- Purpose: Maintains consistent API pattern with other active/trashed endpoints
-- Parameters:
--   $1: search_term - Optional filter text for name/email
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns:
--   Active user records with total_count
-- Business Logic:
--   - Same functionality as GetUsers
--   - Exists for consistency in API design patterns
-- Note: Could be consolidated with GetUsers if duplicate functionality is undesired
-- name: GetUsersActive :many
SELECT
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM users
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR firstname ILIKE '%' || $1 || '%'
        OR lastname ILIKE '%' || $1 || '%'
        OR email ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetUserTrashed: Retrieves paginated list of soft-deleted users
-- Purpose: View and manage deleted users for potential restoration
-- Parameters:
--   $1: search_term - Optional text to filter trashed users
--   $2: limit - Maximum records per page
--   $3: offset - Records to skip
-- Returns:
--   Trashed user records with total_count
-- Business Logic:
--   - Only returns soft-deleted records (deleted_at IS NOT NULL)
--   - Maintains same search functionality as active user queries
--   - Preserves chronological sorting (newest first)
--   - Used in user recovery/audit interfaces
--   - Includes total_count for pagination in trash management UI
-- name: GetUserTrashed :many
SELECT
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM users
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR firstname ILIKE '%' || $1 || '%'
        OR lastname ILIKE '%' || $1 || '%'
        OR email ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetUserByID: Retrieves active user by ID
-- Purpose: Fetch specific user details
-- Parameters:
--   $1: user_id - ID of user to retrieve
-- Returns: Full user record if found and active
-- Business Logic:
--   - Excludes deleted users
--   - Used for user profile viewing/editing
--   - Primary lookup for user management
-- name: GetUserByID :one
SELECT
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at
FROM users
WHERE
    user_id = $1
    AND deleted_at IS NULL;

-- GetUserByIdAndVerify: Retrieves active and verified user by ID
-- Purpose: Fetch user by ID with validation
-- Parameters:
--   $1: user_id - ID of the user to retrieve
-- Returns: Full user record if found and active
-- Business Logic:
--   - Excludes deleted users
--   - Use for cross-service validation and ownership checks
-- name: GetUserByIdAndVerify :one
SELECT user_id, email
FROM users
WHERE
    user_id = $1
    AND deleted_at IS NULL
    AND is_verified = true;

-- GetUserByEmail: Retrieves active user by email
-- Purpose: Lookup user by email address (for authentication)
-- Parameters:
--   $1: email - Exact email address to search for
-- Returns: User record if found and active
-- Business Logic:
--   - Case-sensitive exact match on email
--   - Excludes deleted users
--   - Used during login/authentication flows
--   - Helps prevent duplicate accounts
-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- GetUserByEmailAndVerify: Retrieves active and verified user by email
-- Purpose: Used for login/auth verification
-- Parameters:
--   $1: email - Email address
-- Returns: User record if found, active and valid
-- Business Logic:
--   - Excludes deleted users
--   - Used during authentication flow
-- name: GetUserByEmailAndVerify :one
SELECT
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at
FROM users
WHERE
    email = $1
    AND deleted_at IS NULL
    AND is_verified = true;

-- name: GetUserByVerificationCode :one
-- Purpose: Fetch a user based on their verification code.
-- Parameters:
--   $1: verification_code - The verification code of the user to fetch.
-- Returns:
--   - User record matching the provided verification code.
-- Business Logic:
--   - Filters the users table to find a user based on their verification code.
-- name: GetUserByVerificationCode :one
SELECT
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at
FROM users
WHERE
    verification_code = $1;

-- name: GetUserByEmailWithPassword :one
SELECT user_id, email, password
FROM users
WHERE
    email = $1
    AND deleted_at IS NULL;

-- CreateUser: Creates a new user account
-- Purpose: Register a new user in the system
-- Parameters:
--   $1: firstname - User's first name
--   $2: lastname - User's last name
--   $3: email - User's email address (must be unique)
--   $4: password - Hashed password string
-- Returns: The complete created user record
-- Business Logic:
--   - Sets created_at and updated_at timestamps automatically
--   - Requires all mandatory user fields
--   - Email must be unique across the system
--   - Password should be pre-hashed before insertion
-- name: CreateUser :one
INSERT INTO
    users (
        firstname,
        lastname,
        email,
        password,
        verification_code,
        is_verified,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        current_timestamp,
        current_timestamp
    )
RETURNING
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at;

-- UpdateUser: Modifies user account information
-- Purpose: Update user profile details
-- Parameters:
--   $1: user_id - ID of user to update
--   $2: firstname - Updated first name
--   $3: lastname - Updated last name
--   $4: email - Updated email address
--   $5: password - New hashed password (optional)
-- Returns: Updated user record
-- Business Logic:
--   - Auto-updates updated_at timestamp
--   - Only modifies active (non-deleted) users
--   - Validates email uniqueness
--   - Password field optional (can maintain existing)
-- name: UpdateUser :one
UPDATE users
SET
    firstname = $2,
    lastname = $3,
    email = $4,
    password = $5,
    updated_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
RETURNING
    *;

-- name: UpdateUserPassword :one
UPDATE users
SET
    password = $2,
    updated_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
RETURNING
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at;

-- name: UpdateUserIsVerified :one
UPDATE users
SET
    is_verified = $2,
    updated_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
RETURNING
    user_id,
    firstname,
    lastname,
    email,
    password,
    created_at,
    updated_at;

-- TrashUser: Soft-deletes a user account
-- Purpose: Deactivate user without permanent deletion
-- Parameters:
--   $1: user_id - ID of user to deactivate
-- Returns: The soft-deleted user record
-- Business Logic:
--   - Sets deleted_at timestamp to current time
--   - Only processes currently active users
--   - Preserves all user data for potential restoration
--   - Prevents login while deleted
-- name: TrashUser :one
UPDATE users
SET
    deleted_at = current_timestamp
WHERE
    user_id = $1
    AND deleted_at IS NULL
RETURNING
    user_id,
    firstname,
    lastname,
    email,
    created_at,
    updated_at,
    deleted_at;

-- RestoreUser: Recovers a soft-deleted user
-- Purpose: Reactivate a previously deactivated user
-- Parameters:
--   $1: user_id - ID of user to restore
-- Returns: The restored user record
-- Business Logic:
--   - Nullifies the deleted_at field
--   - Only works on previously deleted users
--   - Restores full account access
--   - Maintains all original user data
-- name: RestoreUser :one
UPDATE users
SET
    deleted_at = NULL
WHERE
    user_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    user_id,
    firstname,
    lastname,
    email,
    created_at,
    updated_at,
    deleted_at;

-- DeleteUserPermanently: Hard-deletes a user account
-- Purpose: Completely remove user from database
-- Parameters:
--   $1: user_id - ID of user to delete
-- Business Logic:
--   - Permanent deletion of already soft-deleted users
--   - No return value (exec-only operation)
--   - Irreversible action - use with caution
--   - Should trigger cleanup of related records
-- name: DeleteUserPermanently :exec
DELETE FROM users WHERE user_id = $1 AND deleted_at IS NOT NULL;

-- RestoreAllUsers: Mass restoration of deleted users
-- Purpose: Recover all trashed users at once
-- Business Logic:
--   - Reactivates all soft-deleted users
--   - No parameters needed (bulk operation)
--   - Typically used during system recovery
--   - Maintains all original user data
-- name: RestoreAllUsers :exec
UPDATE users
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentUsers: Purges all trashed users
-- Purpose: Clean up all soft-deleted user records
-- Business Logic:
--   - Irreversible bulk deletion operation
--   - Only affects already soft-deleted users
--   - Typically used during database maintenance
--   - Should be restricted to admin users
-- name: DeleteAllPermanentUsers :exec
DELETE FROM users WHERE deleted_at IS NOT NULL;