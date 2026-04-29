-- GetMerchants: Retrieves paginated list of active merchants with search capability
-- Purpose: List all active merchants for management UI
-- Parameters:
--   $1: search_term - Optional text to filter merchants by name or email (NULL for no filter)
--   $2: limit - Maximum number of records to return (pagination limit)
--   $3: offset - Number of records to skip (pagination offset)
-- Returns:
--   All merchant fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted merchants (deleted_at IS NULL)
--   - Supports partial text matching on name and contact_email fields (case-insensitive ILIKE)
--   - Returns newest merchants first (created_at DESC)
--   - Provides total_count for client-side pagination calculations
--   - Uses window function COUNT(*) OVER() for efficient total count
-- name: GetMerchants :many
SELECT
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM merchants
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
        OR contact_email ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetMerchantsActive: Retrieves paginated list of active merchants (identical to GetMerchants)
-- Purpose: Maintains consistent API pattern with other active/trashed endpoints
-- Parameters:
--   $1: search_term - Optional filter text for name/email
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns:
--   Active merchant records with total_count
-- Business Logic:
--   - Same functionality as GetMerchants
--   - Exists for consistency in API design patterns
-- Note: Could be consolidated with GetMerchants if duplicate functionality is undesired
-- name: GetMerchantsActive :many
SELECT
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM merchants
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
        OR contact_email ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetMerchantsTrashed: Retrieves paginated list of soft-deleted merchants
-- Purpose: View and manage deleted merchants for potential restoration
-- Parameters:
--   $1: search_term - Optional text to filter trashed merchants
--   $2: limit - Maximum records per page
--   $3: offset - Records to skip
-- Returns:
--   Trashed merchant records with total_count
-- Business Logic:
--   - Only returns soft-deleted records (deleted_at IS NOT NULL)
--   - Maintains same search functionality as active merchant queries
--   - Preserves chronological sorting (newest first)
--   - Used in merchant recovery/audit interfaces
--   - Includes total_count for pagination in trash management UI
-- name: GetMerchantsTrashed :many
SELECT
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM merchants
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR name ILIKE '%' || $1 || '%'
        OR contact_email ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- CreateMerchant: Creates a new merchant account
-- Purpose: Register a new merchant in the system
-- Parameters:
--   $1: user_id - Associated user account ID
--   $2: name - Business name
--   $3: description - Business description
--   $4: address - Physical address
--   $5: contact_email - Business email
--   $6: contact_phone - Business phone
--   $7: status - Account status (active/inactive)
-- Returns: The created merchant record
-- Business Logic:
--   - Sets created_at timestamp automatically
--   - Requires all mandatory merchant fields
--   - Status defaults to 'active' unless specified otherwise
-- name: CreateMerchant :one
INSERT INTO
    merchants (
        user_id,
        name,
        description,
        address,
        contact_email,
        contact_phone,
        status
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at;

-- GetMerchantByID: Retrieves active merchant by ID
-- Purpose: Fetch merchant details for display/editing
-- Parameters:
--   $1: merchant_id - ID of merchant to retrieve
-- Returns: Full merchant record if found and active
-- Business Logic:
--   - Excludes soft-deleted records
--   - Returns single record or nothing
--   - Used for merchant profile viewing and editing
-- name: GetMerchantByID :one
SELECT
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at
FROM merchants
WHERE
    merchant_id = $1
    AND deleted_at IS NULL;

-- UpdateMerchant: Modifies merchant information
-- Purpose: Update merchant profile details
-- Parameters:
--   $1: merchant_id - Target merchant ID
--   $2: name - Updated business name
--   $3: description - Updated description
--   $4: address - Updated physical address
--   $5: contact_email - Updated email
--   $6: contact_phone - Updated phone
--   $7: status - Updated account status
-- Returns: Updated merchant record
-- Business Logic:
--   - Automatically updates updated_at timestamp
--   - Only affects active (non-deleted) records
--   - Validates all required fields
--   - Returns modified record for confirmation
-- name: UpdateMerchant :one
UPDATE merchants
SET
    name = $2,
    description = $3,
    address = $4,
    contact_email = $5,
    contact_phone = $6,
    status = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at;

-- name: UpdateMerchantStatus :one
UPDATE merchants
SET
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at;

-- TrashMerchant: Soft-deletes a merchant account
-- Purpose: Deactivate merchant without permanent deletion
-- Parameters:
--   $1: merchant_id - ID of merchant to deactivate
-- Returns: The soft-deleted merchant record
-- Business Logic:
--   - Sets deleted_at timestamp to current time
--   - Only processes currently active records
--   - Allows recovery via restore function
--   - Maintains referential integrity
-- name: TrashMerchant :one
UPDATE merchants
SET
    deleted_at = current_timestamp
WHERE
    merchant_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at,
    deleted_at;

-- RestoreMerchant: Recovers a soft-deleted merchant
-- Purpose: Reactivate a previously deactivated merchant
-- Parameters:
--   $1: merchant_id - ID of merchant to restore
-- Returns: The restored merchant record
-- Business Logic:
--   - Nullifies the deleted_at field
--   - Only works on previously deleted records
--   - Preserves all original merchant data
--   - Reactivates associated services
-- name: RestoreMerchant :one
UPDATE merchants
SET
    deleted_at = NULL
WHERE
    merchant_id = $1
    AND deleted_at IS NOT NULL
RETURNING
     merchant_id,
    user_id,
    name,
    description,
    address,
    contact_email,
    contact_phone,
    status,
    created_at,
    updated_at,
    deleted_at;

-- DeleteMerchantPermanently: Hard-deletes a merchant
-- Purpose: Completely remove merchant from database
-- Parameters:
--   $1: merchant_id - ID of merchant to delete
-- Business Logic:
--   - Permanent deletion of already soft-deleted records
--   - No return value (exec-only operation)
--   - Irreversible action - use with caution
--   - Should trigger cleanup of related records
-- name: DeleteMerchantPermanently :exec
DELETE FROM merchants
WHERE
    merchant_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllMerchants: Mass restoration of deleted merchants
-- Purpose: Recover all trashed merchants at once
-- Business Logic:
--   - Reactivates all soft-deleted merchants
--   - No parameters needed (bulk operation)
--   - Typically used during system recovery
--   - Maintains original merchant data
-- name: RestoreAllMerchants :exec
UPDATE merchants
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentMerchants: Purges all trashed merchants
-- Purpose: Clean up all soft-deleted merchant records
-- Business Logic:
--   - Irreversible bulk deletion operation
--   - Only affects already soft-deleted records
--   - Typically used during database maintenance
--   - Should be restricted to admin users
-- name: DeleteAllPermanentMerchants :exec
DELETE FROM merchants WHERE deleted_at IS NOT NULL;