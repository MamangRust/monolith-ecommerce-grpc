-- AssignRoleToUser: Assigns a role to a user (creates a user-role relation)
-- Purpose: Role management for user access control
-- Parameters:
--   $1: User ID
--   $2: Role ID
-- Returns:
--   user_role_id, user_id, role_id, timestamps (incl. deleted_at for future status check)
-- Business Logic:
--   - Adds a new entry in the user_roles mapping table
--   - Timestamps created_at and updated_at auto-set to current
-- name: AssignRoleToUser :one
INSERT INTO user_roles (
    user_id, 
    role_id, 
    created_at, 
    updated_at
) VALUES (
    $1, 
    $2, 
    current_timestamp, 
    current_timestamp
) RETURNING 
    user_role_id, 
    user_id, 
    role_id, 
    created_at, 
    updated_at, 
    deleted_at;


-- RemoveRoleFromUser: Permanently removes a role from a user
-- Purpose: Hard delete of a user-role mapping (bypasses trash)
-- Parameters:
--   $1: User ID
--   $2: Role ID
-- Business Logic:
--   - Deletes the record instead of soft-deleting
--   - Use cautiously if audit/history is important
-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles
WHERE 
    user_id = $1 
    AND role_id = $2;


-- TrashUserRole: Soft deletes a user-role mapping (moves to trash)
-- Purpose: Temporarily disable a role assignment without permanent deletion
-- Parameters:
--   $1: user_role_id (primary key of the mapping)
-- Business Logic:
--   - Sets deleted_at timestamp, indicating the relation is inactive
-- name: TrashUserRole :exec
UPDATE user_roles
SET 
    deleted_at = current_timestamp
WHERE 
    user_role_id = $1;


-- RestoreUserRole: Restores a trashed user-role relation
-- Purpose: Reactivate a previously soft-deleted user-role
-- Parameters:
--   $1: user_role_id
-- Business Logic:
--   - Clears the deleted_at field to mark as active again
-- name: RestoreUserRole :exec
UPDATE user_roles
SET 
    deleted_at = NULL
WHERE 
    user_role_id = $1;


-- GetTrashedUserRoles: Retrieves all soft-deleted roles for a given user
-- Purpose: Review previously deleted role assignments for recovery or audit
-- Parameters:
--   $1: User ID
-- Returns:
--   user_role_id, user_id, role_id, role_name, timestamps
-- Business Logic:
--   - Joins with roles to show role name
--   - Orders by most recently trashed
-- name: GetTrashedUserRoles :many
SELECT 
    ur.user_role_id,
    ur.user_id,
    ur.role_id,
    r.role_name,
    ur.created_at,
    ur.updated_at,
    ur.deleted_at
FROM 
    user_roles ur
JOIN 
    roles r ON ur.role_id = r.role_id
WHERE 
    ur.user_id = $1
    AND ur.deleted_at IS NOT NULL
ORDER BY 
    ur.deleted_at DESC;
