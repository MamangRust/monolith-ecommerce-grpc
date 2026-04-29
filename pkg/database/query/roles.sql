-- GetRoles: Retrieves all roles (active & trashed) with optional name search and pagination
-- Purpose: General listing of roles regardless of status
-- Parameters:
--   $1: Search query (role name, nullable)
--   $2: Limit (number of records per page)
--   $3: Offset (starting index for pagination)
-- Returns:
--   role_id, role_name, timestamps, and total_count (for pagination support)
-- Business Logic:
--   - Supports fuzzy search on role_name
--   - Includes both active and trashed roles
--   - Useful for admin panels with filters and pagination
-- name: GetRoles :many
SELECT
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    roles
WHERE
    $1::TEXT IS NULL OR role_name ILIKE '%' || $1 || '%'
ORDER BY
    created_at ASC
LIMIT $2 OFFSET $3;



-- GetRole: Retrieves role details by role_id
-- Purpose: Fetch a single role record (regardless of deleted status)
-- Parameters:
--   $1: Role ID
-- Returns:
--   role_id, role_name, and timestamps
-- name: GetRole :one
SELECT
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at
FROM
    roles
WHERE
    role_id = $1;


-- GetRoleByName: Retrieves role by exact role name
-- Purpose: Check role existence or fetch role info based on name
-- Parameters:
--   $1: Role name (exact match)
-- Returns:
--   role_id, role_name, and timestamps
-- name: GetRoleByName :one
SELECT
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at
FROM
    roles
WHERE
    role_name = $1;


-- GetUserRoles: Retrieves all roles assigned to a specific user
-- Purpose: Identify the access level(s) of a user
-- Parameters:
--   $1: User ID
-- Returns:
--   List of roles (id, name, timestamps)
-- name: GetUserRoles :many
SELECT
    r.role_id,
    r.role_name,
    r.created_at,
    r.updated_at,
    r.deleted_at
FROM
    roles r
JOIN
    user_roles ur ON ur.role_id = r.role_id
WHERE
    ur.user_id = $1
ORDER BY
    r.created_at ASC;


-- GetActiveRoles: Retrieves only active (non-deleted) roles with optional search and pagination
-- Purpose: Display roles that are currently usable in the system
-- Parameters:
--   $1: Search query (nullable)
--   $2: Limit
--   $3: Offset
-- Returns:
--   role_id, role_name, timestamps, and total_count
-- name: GetActiveRoles :many
SELECT
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    roles
WHERE
    deleted_at IS NULL
    AND ($1::TEXT IS NULL OR role_name ILIKE '%' || $1 || '%')
ORDER BY
    created_at ASC
LIMIT $2 OFFSET $3;

-- GetTrashedRoles: Retrieves only soft-deleted roles with optional search and pagination
-- Purpose: For trash/recycle bin management
-- Parameters:
--   $1: Search query (nullable)
--   $2: Limit
--   $3: Offset
-- Returns:
--   role_id, role_name, timestamps, and total_count
-- name: GetTrashedRoles :many
SELECT
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER() AS total_count
FROM
    roles
WHERE
    deleted_at IS NOT NULL
    AND ($1::TEXT IS NULL OR role_name ILIKE '%' || $1 || '%')
ORDER BY
    deleted_at DESC
LIMIT $2 OFFSET $3;



-- CreateRole: Inserts a new role into the system
-- Purpose: Add new role definitions (e.g., Admin, Cashier, etc.)
-- Parameters:
--   $1: Role name
-- Returns:
--   Newly created role's full data (including timestamps)
-- name: CreateRole :one
INSERT INTO roles (
    role_name,
    created_at,
    updated_at
) VALUES (
    $1,
    current_timestamp,
    current_timestamp
) RETURNING
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at;

-- UpdateRole: Updates role name by ID
-- Purpose: Modify role information (e.g., name correction)
-- Parameters:
--   $1: Role ID
--   $2: New role name
-- Returns:
--   Updated role's data
-- name: UpdateRole :one
UPDATE roles
SET
    role_name = $2,
    updated_at = current_timestamp
WHERE
    role_id = $1
RETURNING
    role_id,
    role_name,
    created_at,
    updated_at,
    deleted_at;


-- TrashRole: Soft-deletes a role (moves to trash)
-- Purpose: Mark role as deleted without removing it permanently
-- Parameters:
--   $1: Role ID
-- name: TrashRole :one
UPDATE roles
SET
    deleted_at = current_timestamp
WHERE
    role_id = $1
RETURNING *;


-- RestoreRole: Restores a previously trashed role
-- Purpose: Undelete a soft-deleted role
-- Parameters:
--   $1: Role ID
-- name: RestoreRole :one
UPDATE roles
SET
    deleted_at = NULL
WHERE
    role_id = $1
RETURNING *;


-- DeletePermanentRole: Permanently deletes a trashed role
-- Purpose: Remove role from DB after soft delete
-- Parameters:
--   $1: Role ID
-- name: DeletePermanentRole :exec
DELETE FROM roles
WHERE
    role_id = $1 AND deleted_at IS NOT NULL;


-- RestoreAllRoles: Restores all soft-deleted roles in bulk
-- Purpose: Bulk recovery of all trashed roles
-- Parameters: None
-- name: RestoreAllRoles :exec
UPDATE roles
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;


-- DeleteAllPermanentRoles: Permanently deletes all soft-deleted roles
-- Purpose: Bulk cleanup of trashed roles
-- Parameters: None
-- name: DeleteAllPermanentRoles :exec
DELETE FROM roles
WHERE
    deleted_at IS NOT NULL;
