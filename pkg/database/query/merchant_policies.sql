-- GetMerchantPolicies: Retrieves all merchant policies regardless of merchant status
-- name: GetMerchantPolicies :many
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns: Merchant policy records with merchant name and total_count
SELECT
    mp.merchant_policy_id,
    mp.merchant_id,
    mp.policy_type,
    mp.title,
    mp.description,
    mp.created_at,
    mp.updated_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_policies mp
    JOIN merchants m ON mp.merchant_id = m.merchant_id
WHERE
    mp.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantPoliciesActive: Retrieves merchant policies for active merchants
-- name: GetMerchantPoliciesActive :many
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns: Active merchant policy records with merchant name and total_count
SELECT
    mp.merchant_policy_id,
    mp.merchant_id,
    mp.policy_type,
    mp.title,
    mp.description,
    mp.created_at,
    mp.updated_at,
    mp.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_policies mp
    JOIN merchants m ON mp.merchant_id = m.merchant_id
WHERE
    mp.deleted_at IS NULL
    AND m.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantPoliciesTrashed: Retrieves merchant policies for deleted merchants
-- name: GetMerchantPoliciesTrashed :many
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns: Trashed merchant policy records with merchant name and total_count
SELECT
    mp.merchant_policy_id,
    mp.merchant_id,
    mp.policy_type,
    mp.title,
    mp.description,
    mp.created_at,
    mp.updated_at,
    mp.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_policies mp
    JOIN merchants m ON mp.merchant_id = m.merchant_id
WHERE
    mp.deleted_at IS NOT NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantPolicy: Retrieves a single merchant policy that is not soft-deleted
-- Parameters:
--   $1: merchant_policy_id - ID of the merchant policy
-- Business Logic:
--   - Returns the merchant policy where deleted_at IS NULL
-- Returns: A single merchant_policies record
-- name: GetMerchantPolicy :one
SELECT mp.merchant_policy_id, mp.merchant_id, mp.policy_type, mp.title, mp.description, mp.created_at, mp.updated_at
FROM merchant_policies mp
WHERE
    merchant_policy_id = $1
    AND deleted_at IS NULL;

-- CreateMerchantPolicy: Inserts a new merchant policy
-- Purpose: Register a new policy for a merchant
-- Parameters:
--   $1: merchant_id
--   $2: policy_type
--   $3: title
--   $4: description
-- Returns: The newly created policy record
-- Business Logic:
--   - Sets created_at, updated_at automatically
-- name: CreateMerchantPolicy :one
INSERT INTO
    merchant_policies (
        merchant_id,
        policy_type,
        title,
        description
    )
VALUES ($1, $2, $3, $4)
RETURNING
    merchant_policy_id,
    merchant_id,
    policy_type,
    title,
    description,
    created_at,
    updated_at;

-- UpdateMerchantPolicy: Updates an existing merchant policy
-- Purpose: Modify policy details
-- Parameters:
--   $1: merchant_policy_id
--   $2: policy_type
--   $3: title
--   $4: description
-- Returns: The updated policy record
-- Business Logic:
--   - Only affects active (non-deleted) records
--   - Automatically updates updated_at timestamps
-- name: UpdateMerchantPolicy :one
UPDATE merchant_policies
SET
    policy_type = $2,
    title = $3,
    description = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_policy_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_policy_id,
    merchant_id,
    policy_type,
    title,
    description,
    created_at,
    updated_at;

-- TrashMerchantPolicy: Soft-deletes a merchant policy
-- Purpose: Temporarily deactivate a policy without permanent deletion
-- Parameters:
--   $1: merchant_policy_id
-- Returns: The soft-deleted policy record
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only affects active records
-- name: TrashMerchantPolicy :one
UPDATE merchant_policies
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    merchant_policy_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_policy_id,
    merchant_id,
    policy_type,
    title,
    description,
    created_at,
    updated_at,
    deleted_at;

-- RestoreMerchantPolicy: Restores a soft-deleted policy
-- Purpose: Reactivate a previously deleted policy
-- Parameters:
--   $1: merchant_policy_id
-- Returns: The restored policy record
-- Business Logic:
--   - Clears deleted_at
--   - Only works on previously soft-deleted entries
-- name: RestoreMerchantPolicy :one
UPDATE merchant_policies
SET
    deleted_at = NULL
WHERE
    merchant_policy_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    merchant_policy_id,
    merchant_id,
    policy_type,
    title,
    description,
    created_at,
    updated_at,
    deleted_at;

-- DeleteMerchantPolicyPermanently: Hard-deletes a policy record
-- Purpose: Permanently remove a soft-deleted policy
-- Parameters:
--   $1: merchant_policy_id
-- Business Logic:
--   - Irreversible delete
--   - Only affects already soft-deleted records
-- name: DeleteMerchantPolicyPermanently :exec
DELETE FROM merchant_policies
WHERE
    merchant_policy_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllMerchantPolicies: Restores all soft-deleted merchant policies
-- Purpose: Bulk recovery operation
-- Business Logic:
--   - Clears deleted_at on all trashed records
-- name: RestoreAllMerchantPolicies :exec
UPDATE merchant_policies
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteMerchantPolicyPermanently: Hard-deletes a policy record
-- Purpose: Permanently remove a soft-deleted policy
-- Parameters:
--   $1: merchant_policy_id
-- Business Logic:
--   - Irreversible delete
--   - Only affects already soft-deleted records
-- name: DeleteAllMerchantPolicyPermanently :exec
DELETE FROM merchant_policies WHERE deleted_at IS NOT NULL;