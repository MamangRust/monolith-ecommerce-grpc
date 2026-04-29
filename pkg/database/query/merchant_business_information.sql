-- GetMerchantsBusinessInformation: Retrieves all business info records with pagination and optional search
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Fetches all merchant business information regardless of merchant status (active or deleted)
--   - Supports case-insensitive partial search by merchant name
--   - Applies pagination using limit and offset
-- Returns: All business info records matching the search, with total_count for pagination
-- name: GetMerchantsBusinessInformation :many
SELECT
    mbi.merchant_business_info_id,
    mbi.merchant_id,
    mbi.business_type,
    mbi.tax_id,
    mbi.established_year,
    mbi.number_of_employees,
    mbi.website_url,
    mbi.created_at,
    mbi.updated_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_business_information mbi
    JOIN merchants m ON mbi.merchant_id = m.merchant_id
WHERE
    mbi.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantsBusinessInformationActive: Retrieves business info for active (non-deleted) merchants with pagination and optional search
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Fetches merchant business information for only active merchants (deleted_at IS NULL)
--   - Supports case-insensitive partial search by merchant name
--   - Applies pagination using limit and offset
-- Returns: Business info records for active merchants matching the search, with total_count
-- name: GetMerchantsBusinessInformationActive :many
SELECT
    mbi.merchant_business_info_id,
    mbi.merchant_id,
    mbi.business_type,
    mbi.tax_id,
    mbi.established_year,
    mbi.number_of_employees,
    mbi.website_url,
    mbi.created_at,
    mbi.updated_at,
    mbi.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_business_information mbi
    JOIN merchants m ON mbi.merchant_id = m.merchant_id
WHERE
    mbi.deleted_at IS NULL
    AND m.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantBusinessInformationTrashed: Retrieves business info for deleted merchants with pagination and optional search
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Fetches merchant business information for only deleted merchants (deleted_at IS NOT NULL)
--   - Supports case-insensitive partial search by merchant name
--   - Applies pagination using limit and offset
-- Returns: Business info records for deleted merchants matching the search, with total_count
-- name: GetMerchantsBusinessInformationTrashed :many
SELECT
    mbi.merchant_business_info_id,
    mbi.merchant_id,
    mbi.business_type,
    mbi.tax_id,
    mbi.established_year,
    mbi.number_of_employees,
    mbi.website_url,
    mbi.created_at,
    mbi.updated_at,
    mbi.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_business_information mbi
    JOIN merchants m ON mbi.merchant_id = m.merchant_id
WHERE
    mbi.deleted_at IS NOT NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantBusinessInformation: Retrieves a single business information record that is not soft-deleted
-- Parameters:
--   $1: merchant_business_info_id - ID of the business info
-- Business Logic:
--   - Returns the merchant business information where deleted_at IS NULL
-- Returns: A single merchant_business_information record
-- name: GetMerchantBusinessInformation :one
SELECT mbi.merchant_business_info_id, mbi.merchant_id, mbi.business_type, mbi.tax_id, mbi.established_year, mbi.number_of_employees, mbi.website_url, mbi.created_at, mbi.updated_at
FROM
    merchant_business_information mbi
WHERE
    merchant_business_info_id = $1
    AND deleted_at IS NULL;

-- CreateMerchantBusinessInformation: Inserts a new business info record
-- Purpose: Register extended business info for a merchant
-- Parameters:
--   $1: merchant_id
--   $2: business_type
--   $3: tax_id
--   $4: established_year
--   $5: number_of_employees
--   $6: website_url
-- Returns: The newly created business info record
-- Business Logic:
--   - Sets created_at timestamp automatically
--   - Validates all required fields
-- name: CreateMerchantBusinessInformation :one
INSERT INTO
    merchant_business_information (
        merchant_id,
        business_type,
        tax_id,
        established_year,
        number_of_employees,
        website_url
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    merchant_business_info_id,
    merchant_id,
    business_type,
    tax_id,
    established_year,
    number_of_employees,
    website_url,
    created_at,
    updated_at;

-- UpdateMerchantBusinessInformation: Updates existing business info
-- Purpose: Modify business information details
-- Parameters:
--   $1: merchant_business_info_id
--   $2: business_type
--   $3: tax_id
--   $4: established_year
--   $5: number_of_employees
--   $6: website_url
-- Returns: The updated business info record
-- Business Logic:
--   - Only affects active (non-deleted) records
--   - Automatically updates the updated_at timestamp
-- name: UpdateMerchantBusinessInformation :one
UPDATE merchant_business_information
SET
    business_type = $2,
    tax_id = $3,
    established_year = $4,
    number_of_employees = $5,
    website_url = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_business_info_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_business_info_id,
    merchant_id,
    business_type,
    tax_id,
    established_year,
    number_of_employees,
    website_url,
    created_at,
    updated_at;

-- TrashMerchantBusinessInformation: Soft-deletes a business info record
-- Purpose: Deactivate business info without removing from DB
-- Parameters:
--   $1: merchant_business_info_id
-- Returns: The soft-deleted record
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only affects active records
-- name: TrashMerchantBusinessInformation :one
UPDATE merchant_business_information
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    merchant_business_info_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_business_info_id,
    merchant_id,
    business_type,
    tax_id,
    established_year,
    number_of_employees,
    website_url,
    created_at,
    updated_at,
    deleted_at;

-- RestoreMerchantBusinessInformation: Restores a soft-deleted business info record
-- Purpose: Reactivate a previously trashed record
-- Parameters:
--   $1: merchant_business_info_id
-- Returns: The restored business info record
-- Business Logic:
--   - Clears deleted_at
--   - Only works for soft-deleted entries
-- name: RestoreMerchantBusinessInformation :one
UPDATE merchant_business_information
SET
    deleted_at = NULL
WHERE
    merchant_business_info_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    merchant_business_info_id,
    merchant_id,
    business_type,
    tax_id,
    established_year,
    number_of_employees,
    website_url,
    created_at,
    updated_at,
    deleted_at;

-- DeleteMerchantBusinessInformationPermanently: Hard-deletes a single record
-- Purpose: Completely remove soft-deleted business info
-- Parameters:
--   $1: merchant_business_info_id
-- Business Logic:
--   - Irreversible operation
--   - Only affects already soft-deleted records
-- name: DeleteMerchantBusinessInformationPermanently :exec
DELETE FROM merchant_business_information
WHERE
    merchant_business_info_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllMerchantBusinessInformation: Restores all soft-deleted records
-- Purpose: Bulk recovery operation
-- Business Logic:
--   - Clears deleted_at on all trashed records
--   - Useful during data recovery or admin resets
-- name: RestoreAllMerchantBusinessInformation :exec
UPDATE merchant_business_information
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentMerchantBusinessInformation: Hard-deletes all trashed records
-- Purpose: Clean up all soft-deleted business info entries
-- Business Logic:
--   - Irreversible bulk delete
--   - Used during database cleanup/maintenance
-- name: DeleteAllPermanentMerchantBusinessInformation :exec
DELETE FROM merchant_business_information
WHERE
    deleted_at IS NOT NULL;