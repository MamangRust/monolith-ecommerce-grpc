-- GetMerchantCertificationsAndAwards: Retrieves all certification records with pagination and optional search
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Fetches certifications regardless of merchant status (active or deleted)
--   - Supports case-insensitive partial search by merchant name
--   - Applies pagination using limit and offset
-- Returns: Certification records matching the search, with total_count
-- name: GetMerchantCertificationsAndAwards :many
SELECT
    mca.merchant_certification_id,
    mca.merchant_id,
    mca.title,
    mca.description,
    mca.issued_by,
    mca.issue_date,
    mca.expiry_date,
    mca.certificate_url,
    mca.created_at,
    mca.updated_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_certifications_and_awards mca
    JOIN merchants m ON mca.merchant_id = m.merchant_id
WHERE
    mca.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantCertificationsAndAwardsActive: Retrieves certifications for active merchants with pagination and optional search
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Fetches certifications only for active merchants (deleted_at IS NULL)
--   - Supports case-insensitive partial search by merchant name
--   - Applies pagination using limit and offset
-- Returns: Active merchant certification records matching the search, with total_count
-- name: GetMerchantCertificationsAndAwardsActive :many
SELECT
    mca.merchant_certification_id,
    mca.merchant_id,
    mca.title,
    mca.description,
    mca.issued_by,
    mca.issue_date,
    mca.expiry_date,
    mca.certificate_url,
    mca.created_at,
    mca.updated_at,
    mca.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_certifications_and_awards mca
    JOIN merchants m ON mca.merchant_id = m.merchant_id
WHERE
    mca.deleted_at IS NULL
    AND m.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantCertificationsAndAwardsTrashed: Retrieves certifications for deleted merchants with pagination and optional search
-- Parameters:
--   $1: search - Keyword to filter merchant_name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Fetches certifications only for deleted merchants (deleted_at IS NOT NULL)
--   - Supports case-insensitive partial search by merchant name
--   - Applies pagination using limit and offset
-- Returns: Trashed merchant certification records matching the search, with total_count
-- name: GetMerchantCertificationsAndAwardsTrashed :many
SELECT
    mca.merchant_certification_id,
    mca.merchant_id,
    mca.title,
    mca.description,
    mca.issued_by,
    mca.issue_date,
    mca.expiry_date,
    mca.certificate_url,
    mca.created_at,
    mca.updated_at,
    mca.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count
FROM
    merchant_certifications_and_awards mca
    JOIN merchants m ON mca.merchant_id = m.merchant_id
WHERE
    mca.deleted_at IS NOT NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetMerchantCertificationOrAward: Retrieves a single merchant award or certification that is not soft-deleted
-- Parameters:
--   $1: merchant_certification_id - ID of the award/certification
-- Business Logic:
--   - Returns the merchant award or certification where deleted_at IS NULL
-- Returns: A single merchant_certifications_and_awards record
-- name: GetMerchantCertificationOrAward :one
SELECT mca.merchant_certification_id, mca.merchant_id, mca.title, mca.description, mca.issued_by, mca.issue_date, mca.expiry_date, mca.certificate_url, mca.created_at, mca.updated_at
FROM
    merchant_certifications_and_awards mca
WHERE
    merchant_certification_id = $1
    AND deleted_at IS NULL;

-- CreateMerchantCertificationOrAward: Inserts a new certification or award
-- Purpose: Register a new certification or award for a merchant
-- Parameters:
--   $1: merchant_id
--   $2: title
--   $3: description
--   $4: issued_by
--   $5: issue_date
--   $6: expiry_date
--   $7: certificate_url
-- Returns: The newly created certification record
-- Business Logic:
--   - Sets created_at automatically
--   - Requires merchant_id and title
-- name: CreateMerchantCertificationOrAward :one
INSERT INTO
    merchant_certifications_and_awards (
        merchant_id,
        title,
        description,
        issued_by,
        issue_date,
        expiry_date,
        certificate_url
    )
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    merchant_certification_id,
    merchant_id,
    title,
    description,
    issued_by,
    issue_date,
    expiry_date,
    certificate_url,
    created_at,
    updated_at;

-- UpdateMerchantCertificationOrAward: Updates an existing certification or award
-- Purpose: Modify certification or award details
-- Parameters:
--   $1: merchant_certification_id
--   $2: title
--   $3: description
--   $4: issued_by
--   $5: issue_date
--   $6: expiry_date
--   $7: certificate_url
-- Returns: The updated certification record
-- Business Logic:
--   - Only affects active (non-deleted) records
--   - Automatically updates the updated_at timestamp
-- name: UpdateMerchantCertificationOrAward :one
UPDATE merchant_certifications_and_awards
SET
    title = $2,
    description = $3,
    issued_by = $4,
    issue_date = $5,
    expiry_date = $6,
    certificate_url = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_certification_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_certification_id,
    merchant_id,
    title,
    description,
    issued_by,
    issue_date,
    expiry_date,
    certificate_url,
    created_at,
    updated_at;

-- TrashMerchantCertificationOrAward: Soft-deletes a certification or award
-- Purpose: Deactivate a record without permanent deletion
-- Parameters:
--   $1: merchant_certification_id
-- Returns: The soft-deleted certification record
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only affects active records
-- name: TrashMerchantCertificationOrAward :one
UPDATE merchant_certifications_and_awards
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    merchant_certification_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_certification_id,
    merchant_id,
    title,
    description,
    issued_by,
    issue_date,
    expiry_date,
    certificate_url,
    created_at,
    updated_at,
    deleted_at;

-- RestoreMerchantCertificationOrAward: Restores a soft-deleted certification
-- Purpose: Reactivate a previously deleted record
-- Parameters:
--   $1: merchant_certification_id
-- Returns: The restored record
-- Business Logic:
--   - Clears deleted_at
--   - Only works on previously soft-deleted entries
-- name: RestoreMerchantCertificationOrAward :one
UPDATE merchant_certifications_and_awards
SET
    deleted_at = NULL
WHERE
    merchant_certification_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    merchant_certification_id,
    merchant_id,
    title,
    description,
    issued_by,
    issue_date,
    expiry_date,
    certificate_url,
    created_at,
    updated_at,
    deleted_at;

-- DeleteMerchantCertificationOrAwardPermanently: Hard-deletes a certification record
-- Purpose: Permanently remove a soft-deleted record
-- Parameters:
--   $1: merchant_certification_id
-- Business Logic:
--   - Only affects already soft-deleted records
--   - Irreversible action
-- name: DeleteMerchantCertificationOrAwardPermanently :exec
DELETE FROM merchant_certifications_and_awards
WHERE
    merchant_certification_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllMerchantCertificationsAndAwards: Restores all soft-deleted certifications
-- Purpose: Bulk restore operation
-- Business Logic:
--   - Clears deleted_at field from all trashed records
--   - Useful for recovery or admin batch restore
-- name: RestoreAllMerchantCertificationsAndAwards :exec
UPDATE merchant_certifications_and_awards
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentMerchantCertificationsAndAwards: Hard-deletes all trashed certifications
-- Purpose: Bulk clean-up operation
-- Business Logic:
--   - Irreversible delete for all soft-deleted records
--   - Should be restricted to admin-level actions
-- name: DeleteAllPermanentMerchantCertificationsAndAwards :exec
DELETE FROM merchant_certifications_and_awards
WHERE
    deleted_at IS NOT NULL;