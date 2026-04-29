-- GetReviews: Retrieves paginated list of all reviews (both active and trashed) with search capability
-- Purpose: Display reviews in admin dashboard
-- Parameters:
--   $1: search_term - Optional text to filter reviews by ID or name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All review fields plus total_count of matching records
-- Business Logic:
--   - Includes both active and trashed reviews
--   - Supports partial text matching on review_id (cast as text) and name fields (case-insensitive)
--   - Returns newest reviews first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetReviews :many
SELECT
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM reviews
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR review_id::TEXT ILIKE '%' || $1 || '%'
        OR name ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetReviewsActive: Retrieves paginated list of active reviews with search capability
-- Purpose: Display active reviews in storefront/admin UI
-- Parameters:
--   $1: search_term - Optional text to filter reviews by ID or name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All review fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted reviews (deleted_at IS NULL)
--   - Supports partial text matching on review_id (cast as text) and name fields (case-insensitive)
--   - Returns newest reviews first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetReviewsActive :many
SELECT
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM reviews
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR review_id::TEXT ILIKE '%' || $1 || '%'
        OR name ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetReviewsTrashed: Retrieves paginated list of trashed reviews with search capability
-- Purpose: Display deleted reviews in admin recycle bin
-- Parameters:
--   $1: search_term - Optional text to filter reviews by ID or name (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All review fields plus total_count of matching records
-- Business Logic:
--   - Only includes soft-deleted reviews (deleted_at IS NOT NULL)
--   - Supports partial text matching on review_id (cast as text) and name fields (case-insensitive)
--   - Returns newest reviews first (created_at DESC)
--   - Provides total_count for pagination calculations
-- name: GetReviewsTrashed :many
SELECT
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM reviews
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR review_id::TEXT ILIKE '%' || $1 || '%'
        OR name ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetReviewByProductId :many
-- Retrieves paginated list of product reviews with details and optional rating filter
--
-- Purpose: Display reviews for a specific product in storefront
--
-- Parameters:
--   $1: product_id - ID of the product to get reviews for
--   $2: rating_filter - Optional rating value to filter by (NULL for all ratings)
--   $3: limit - Maximum number of records to return
--   $4: offset - Number of records to skip for pagination
--
-- Returns:
--   - Review fields with aggregated review details (images/videos) as JSON array
--   - Includes total_count for pagination
--
-- Business Logic:
--   - Only shows active (non-deleted) reviews
--   - Filters by specific product only
--   - Optional rating filter (1-5 stars)
--   - Aggregates review details (media attachments) as JSON array
--   - Returns newest reviews first
--   - Uses LEFT JOIN to include reviews without attachments
SELECT
    r.review_id,
    r.user_id,
    r.product_id,
    r.name,
    r.comment,
    r.rating,
    r.created_at,
    r.updated_at,
    r.deleted_at,
    COUNT(*) OVER () AS total_count,
    COALESCE(
        (
            SELECT json_agg(
                    jsonb_build_object(
                        'detail_id', rd.review_detail_id, 'type', rd.type, 'url', rd.url, 'caption', rd.caption, 'created_at', rd.created_at
                    )
                )
            FROM review_details rd
            WHERE
                rd.review_id = r.review_id
        ),
        '[]'
    ) AS review_details
FROM reviews r
WHERE
    r.deleted_at IS NULL
    AND r.product_id = $1
    AND (
        $2::INT IS NULL
        OR r.rating = $2
    )
ORDER BY r.created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: GetReviewByMerchantId :many
-- Retrieves paginated reviews for all products belonging to a merchant
--
-- Purpose: Display all reviews for merchant's products in dashboard
--
-- Parameters:
--   $1: merchant_id - ID of the merchant whose products' reviews to fetch
--   $2: rating_filter - Optional rating value to filter by (NULL for all ratings)
--   $3: limit - Maximum number of records to return
--   $4: offset - Number of records to skip for pagination
--
-- Returns:
--   - Review fields with aggregated review details (images/videos) as JSON array
--   - Includes total_count for pagination
--
-- Business Logic:
--   - Only shows active (non-deleted) reviews
--   - Filters by merchant's products through JOIN
--   - Optional rating filter (1-5 stars)
--   - Aggregates review details (media attachments) as JSON array
--   - Returns newest reviews first
--   - Uses JOIN with products table to ensure merchant ownership
SELECT
    r.review_id,
    r.user_id,
    r.product_id,
    r.name,
    r.comment,
    r.rating,
    r.created_at,
    r.updated_at,
    r.deleted_at,
    COUNT(*) OVER () AS total_count,
    COALESCE(
        (
            SELECT json_agg(
                    jsonb_build_object(
                        'detail_id', rd.review_detail_id, 'type', rd.type, 'url', rd.url, 'caption', rd.caption, 'created_at', rd.created_at
                    )
                )
            FROM review_details rd
            WHERE
                rd.review_id = r.review_id
        ),
        '[]'
    ) AS review_details
FROM reviews r
    JOIN products p ON r.product_id = p.product_id
WHERE
    r.deleted_at IS NULL
    AND p.merchant_id = $1
    AND (
        $2::INT IS NULL
        OR r.rating = $2
    )
ORDER BY r.created_at DESC
LIMIT $3
OFFSET
    $4;

-- GetReviewByID: Retrieves a single active review by ID
-- Purpose: Display review details in UI
-- Parameters:
--   $1: review_id - ID of the review to retrieve
-- Returns:
--   Complete review record if found and active
-- Business Logic:
--   - Only returns non-deleted (active) reviews
--   - Used for displaying individual review details
-- name: GetReviewByID :one
SELECT
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at
FROM reviews
WHERE
    review_id = $1
    AND deleted_at IS NULL;

-- CreateReview: Creates a new product review
-- Purpose: Allow users to submit product reviews
-- Parameters:
--   $1: user_id - ID of the user submitting review
--   $2: product_id - ID of the product being reviewed
--   $3: name - Display name for the review
--   $4: comment - Review text content
--   $5: rating - Numeric rating (typically 1-5)
-- Returns:
--   The complete created review record
-- Business Logic:
--   - Creates a new active review
--   - Requires user_id and product_id for reference
--   - Returns full record for immediate display
-- name: CreateReview :one
INSERT INTO
    reviews (
        user_id,
        product_id,
        name,
        comment,
        rating
    )
VALUES ($1, $2, $3, $4, $5)
RETURNING
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at;

-- UpdateReview: Modifies an existing review
-- Purpose: Allow users to edit their reviews
-- Parameters:
--   $1: review_id - ID of review to update
--   $2: name - Updated display name
--   $3: comment - Updated review text
--   $4: rating - Updated numeric rating
-- Returns:
--   The updated review record
-- Business Logic:
--   - Only updates mutable review fields
--   - Requires review to be active (not deleted)
--   - Automatically updates timestamp
-- name: UpdateReview :one
UPDATE reviews
SET
    name = $2,
    comment = $3,
    rating = $4,
    updated_at = CURRENT_TIMESTAMP
WHERE
    review_id = $1
    AND deleted_at IS NULL
RETURNING
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at;

-- TrashReview: Soft-deletes a review
-- Purpose: Remove review from public view while preserving data
-- Parameters:
--   $1: review_id - ID of review to trash
-- Returns:
--   The trashed review record
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only works on active reviews
--   - Allows recovery via RestoreReview
-- name: TrashReview :one
UPDATE reviews
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    review_id = $1
    AND deleted_at IS NULL
RETURNING
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at,
    deleted_at;

-- RestoreReview: Recovers a soft-deleted review
-- Purpose: Reactivate previously trashed reviews
-- Parameters:
--   $1: review_id - ID of review to restore
-- Returns:
--   The restored review record
-- Business Logic:
--   - Clears deleted_at timestamp
--   - Only works on trashed reviews
--   - Returns review to active status
-- name: RestoreReview :one
UPDATE reviews
SET
    deleted_at = NULL
WHERE
    review_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    review_id,
    user_id,
    product_id,
    name,
    comment,
    rating,
    created_at,
    updated_at,
    deleted_at;

-- DeleteReviewPermanently: Removes a review from database
-- Purpose: Permanent deletion of trashed reviews
-- Parameters:
--   $1: review_id - ID of review to delete
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Physical deletion of record
--   - Only works on already-trashed reviews
--   - Irreversible operation
-- name: DeleteReviewPermanently :exec
DELETE FROM reviews WHERE review_id = $1 AND deleted_at IS NOT NULL;

-- RestoreAllReviews: Recovers all soft-deleted reviews
-- Purpose: Bulk restore from trash/recycle bin
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Clears deleted_at for all trashed reviews
--   - Admin-level operation
--   - Returns all reviews to active status
-- name: RestoreAllReviews :exec
UPDATE reviews SET deleted_at = NULL WHERE deleted_at IS NOT NULL;

-- DeleteAllPermanentReviews: Purges all trashed reviews
-- Purpose: Clean up review recycle bin
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Permanent deletion of all trashed reviews
--   - Admin-level operation
--   - Irreversible bulk deletion
-- name: DeleteAllPermanentReviews :exec
DELETE FROM reviews WHERE deleted_at IS NOT NULL;