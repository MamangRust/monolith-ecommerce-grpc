-- GetReviewDetails: Retrieves all review details (active and trashed)
-- Parameters:
--   $1: search - Keyword to filter caption (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Returns both active and trashed records
--   - Applies case-insensitive partial search on caption
--   - Supports pagination
-- Returns: List of review details with total_count metadata
-- name: GetReviewDetails :many
SELECT
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM review_details
WHERE
    LOWER(COALESCE(caption, '')) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetReviewDetailsActive: Retrieves active (non-deleted) review details
-- Parameters:
--   $1: search - Keyword to filter caption (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Returns only records where deleted_at IS NULL
--   - Applies case-insensitive partial search on caption
--   - Supports pagination
-- Returns: List of active review details with total_count metadata
-- name: GetReviewDetailsActive :many
SELECT
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM review_details
WHERE
    deleted_at IS NULL
    AND LOWER(COALESCE(caption, '')) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetTrashedReviewDetails: Retrieves trashed (soft-deleted) review details
-- Parameters:
--   $1: search - Keyword to filter caption (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Returns only records where deleted_at IS NOT NULL
--   - Applies case-insensitive partial search on caption
--   - Supports pagination
-- Returns: List of trashed review details with total_count metadata
-- name: GetReviewDetailsTrashed :many
SELECT
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM review_details
WHERE
    deleted_at IS NOT NULL
    AND LOWER(COALESCE(caption, '')) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetReviewDetail: Retrieves a single review detail that is not soft-deleted
-- Parameters:
--   $1: review_detail_id - ID of the review detail
-- Business Logic:
--   - Returns the review detail where deleted_at IS NULL
-- Returns: A single review_details record
-- name: GetReviewDetail :one
SELECT
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at
FROM review_details
WHERE
    review_detail_id = $1
    AND deleted_at IS NULL;

-- GetReviewDetailTrashed: Retrieves a single review detail that is not soft-deleted
-- Parameters:
--   $1: review_detail_id - ID of the review detail
-- Business Logic:
--   - Returns the review detail where deleted_at IS NULL
-- Returns: A single review_details record
-- name: GetReviewDetailTrashed :one
SELECT *
FROM review_details
WHERE
    review_detail_id = $1
    AND deleted_at IS NOT NULL;

-- CreateReviewDetail: Inserts a new review detail
-- Parameters:
--   $1: review_id - Foreign key to reviews table
--   $2: type - Either 'photo' or 'video'
--   $3: url - URL of the media
--   $4: caption - Optional caption
-- Returns: The inserted review detail record
-- name: CreateReviewDetail :one
INSERT INTO
    review_details (review_id, type, url, caption)
VALUES ($1, $2, $3, $4)
RETURNING
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at;

-- UpdateReviewDetail: Updates an existing review detail
-- Parameters:
--   $1: type - New media type
--   $2: url - New media URL
--   $3: caption - Updated caption
--   $4: review_detail_id - Target record to update
-- Returns: The updated review detail record
-- name: UpdateReviewDetail :one
UPDATE review_details
SET
    type = $1,
    url = $2,
    caption = $3
WHERE
    review_detail_id = $4
RETURNING
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at;

-- TrashReviewDetail: Soft deletes a review detail
-- Parameters:
--   $1: review_detail_id - ID of the record to soft delete
-- Business Logic:
--   - Only applies soft-delete if record is not already trashed
-- Returns: The soft-deleted review detail
-- name: TrashReviewDetail :one
UPDATE review_details
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    review_detail_id = $1
    AND deleted_at IS NULL
RETURNING
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at,
    deleted_at;

-- RestoreReviewDetail: Restores a soft-deleted review detail
-- Parameters:
--   $1: review_detail_id - ID of the record to restore
-- Business Logic:
--   - Only restores if deleted_at IS NOT NULL
-- Returns: The restored review detail
-- name: RestoreReviewDetail :one
UPDATE review_details
SET
    deleted_at = NULL
WHERE
    review_detail_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    review_detail_id,
    review_id,
    type,
    url,
    caption,
    created_at,
    updated_at,
    deleted_at;

-- DeletePermanentReviewDetail: Permanently deletes a soft-deleted review detail
-- Parameters:
--   $1: review_detail_id - ID of the record to permanently delete
-- Business Logic:
--   - Only deletes if deleted_at IS NOT NULL
-- Returns: No return value (exec command)
-- name: DeletePermanentReviewDetail :exec
DELETE FROM review_details
WHERE
    review_detail_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllReviewDetails: Restores all soft-deleted review details
-- Parameters: None
-- Business Logic:
--   - Sets deleted_at to NULL for all trashed records
-- Returns: No return value (exec command)
-- name: RestoreAllReviewDetails :exec
UPDATE review_details
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentReviewDetails: Permanently deletes all trashed review details
-- Parameters: None
-- Business Logic:
--   - Deletes all records where deleted_at IS NOT NULL
-- Returns: No return value (exec command)
-- name: DeleteAllPermanentReviewDetails :exec
DELETE FROM review_details WHERE deleted_at IS NOT NULL;