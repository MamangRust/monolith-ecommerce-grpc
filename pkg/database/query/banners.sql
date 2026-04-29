-- GetBanners: Retrieves all banners (active and trashed) with optional search and pagination
-- Parameters:
--   $1: search - Keyword to filter banner name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Returns all banners regardless of soft-delete status
--   - Applies case-insensitive partial match on banner name
--   - Supports pagination
-- Returns: List of all banners with total_count metadata
-- name: GetBanners :many
SELECT
    b.banner_id,
    b.name,
    b.start_date,
    b.end_date,
    b.start_time,
    b.end_time,
    b.is_active,
    b.created_at,
    b.updated_at,
    COUNT(*) OVER () AS total_count
FROM banners b
WHERE
    LOWER(name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetBannersActive: Retrieves active banners (not soft-deleted)
-- Parameters:
--   $1: search - Keyword to filter banner name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Returns banners where deleted_at IS NULL
--   - Applies case-insensitive partial match on banner name
--   - Supports pagination
-- Returns: List of active banners with total_count metadata
-- name: GetBannersActive :many
SELECT
    b.banner_id,
    b.name,
    b.start_date,
    b.end_date,
    b.start_time,
    b.end_time,
    b.is_active,
    b.created_at,
    b.updated_at,
    b.deleted_at,
    COUNT(*) OVER () AS total_count
FROM banners b
WHERE
    deleted_at IS NULL
    AND LOWER(name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetTrashedBanners: Retrieves soft-deleted (trashed) banners
-- Parameters:
--   $1: search - Keyword to filter banner name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Business Logic:
--   - Returns banners where deleted_at IS NOT NULL
--   - Applies case-insensitive partial match on banner name
--   - Supports pagination
-- Returns: List of trashed banners with total_count metadata
-- name: GetBannersTrashed :many
SELECT
    b.banner_id,
    b.name,
    b.start_date,
    b.end_date,
    b.start_time,
    b.end_time,
    b.is_active,
    b.created_at,
    b.updated_at,
    b.deleted_at,
    COUNT(*) OVER () AS total_count
FROM banners b
WHERE
    deleted_at IS NOT NULL
    AND LOWER(name) LIKE LOWER(CONCAT('%', $1::text, '%'))
LIMIT $2
OFFSET
    $3;

-- GetBanner: Retrieves a single banner that is not soft-deleted
-- Parameters:
--   $1: banner_id - ID of the banner
-- Business Logic:
--   - Returns the banner where deleted_at IS NULL
-- Returns: A single banners record
-- name: GetBanner :one
SELECT b.banner_id, b.name, b.start_date, b.end_date, b.start_time, b.end_time, b.is_active, b.created_at, b.updated_at
FROM banners b
WHERE
    banner_id = $1
    AND deleted_at IS NULL;

-- CreateBanner: Inserts a new banner
-- name: CreateBanner :one
-- Parameters:
--   $1: name
--   $2: start_date
--   $3: end_date
--   $4: start_time
--   $5: end_time
--   $6: is_active
-- Returns: The newly created banner
INSERT INTO
    banners (
        name,
        start_date,
        end_date,
        start_time,
        end_time,
        is_active
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    banner_id,
    name,
    start_date,
    end_date,
    start_time,
    end_time,
    is_active,
    created_at,
    updated_at;

-- UpdateBanner: Updates an existing banner
-- name: UpdateBanner :one
-- Parameters:
--   $1: banner_id
--   $2: name
--   $3: start_date
--   $4: end_date
--   $5: start_time
--   $6: end_time
--   $7: is_active
-- Returns: The updated banner
UPDATE banners
SET
    name = $2,
    start_date = $3,
    end_date = $4,
    start_time = $5,
    end_time = $6,
    is_active = $7,
    updated_at = CURRENT_TIMESTAMP
WHERE
    banner_id = $1
    AND deleted_at IS NULL
RETURNING
    banner_id,
    name,
    start_date,
    end_date,
    start_time,
    end_time,
    is_active,
    created_at,
    updated_at;

-- TrashBanner: Soft deletes a banner
-- name: TrashBanner :one
-- Parameters:
--   $1: banner_id
-- Returns: The trashed banner
UPDATE banners
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    banner_id = $1
    AND deleted_at IS NULL
RETURNING
    banner_id,
    name,
    start_date,
    end_date,
    start_time,
    end_time,
    is_active,
    created_at,
    updated_at,
    deleted_at;

-- RestoreBanner: Restores a soft-deleted banner
-- name: RestoreBanner :one
-- Parameters:
--   $1: banner_id
-- Returns: The restored banner
UPDATE banners
SET
    deleted_at = NULL
WHERE
    banner_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    banner_id,
    name,
    start_date,
    end_date,
    start_time,
    end_time,
    is_active,
    created_at,
    updated_at,
    deleted_at;

-- DeleteBannerPermanently: Permanently delete a single trashed banner
-- name: DeleteBannerPermanently :exec
-- Parameters:
--   $1: banner_id
DELETE FROM banners
WHERE
    banner_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllBanners: Restore all trashed banners
-- name: RestoreAllBanners :exec
UPDATE banners SET deleted_at = NULL WHERE deleted_at IS NOT NULL;

-- DeleteAllPermanentBanners: Permanently delete all trashed banners
-- name: DeleteAllPermanentBanners :exec
DELETE FROM banners WHERE deleted_at IS NOT NULL; 