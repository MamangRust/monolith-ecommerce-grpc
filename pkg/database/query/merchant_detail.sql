-- GetMerchantDetails: Retrieves all merchant details regardless of merchant status
-- Parameters:
--   $1: search - Keyword to filter merchant name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns: Merchant detail records matching the search, with total_count and social media links
-- name: GetMerchantDetails :many
SELECT
    md.merchant_detail_id,
    md.merchant_id,
    md.display_name,
    md.cover_image_url,
    md.logo_url,
    md.short_description,
    md.website_url,
    md.created_at,
    md.updated_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count,
    json_agg(
        json_build_object(
            'id',
            sml.merchant_social_id,
            'platform',
            sml.platform,
            'url',
            sml.url
        )
    ) AS social_media_links
FROM
    merchant_details md
    JOIN merchants m ON md.merchant_id = m.merchant_id
    LEFT JOIN merchant_social_media_links sml ON sml.merchant_detail_id = md.merchant_detail_id
WHERE
    md.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
GROUP BY
    md.merchant_detail_id,
    m.merchant_id
LIMIT $2
OFFSET
    $3;

-- GetMerchantDetailsActive: Retrieves merchant details for active merchants
-- Parameters:
--   $1: search - Keyword to filter merchant name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns: Active merchant detail records matching the search, with total_count and social media links
-- name: GetMerchantDetailsActive :many
SELECT
    md.merchant_detail_id,
    md.merchant_id,
    md.display_name,
    md.cover_image_url,
    md.logo_url,
    md.short_description,
    md.website_url,
    md.created_at,
    md.updated_at,
    md.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count,
    json_agg(
        json_build_object(
            'id',
            sml.merchant_social_id,
            'platform',
            sml.platform,
            'url',
            sml.url
        )
    ) AS social_media_links
FROM
    merchant_details md
    JOIN merchants m ON md.merchant_id = m.merchant_id
    LEFT JOIN merchant_social_media_links sml ON sml.merchant_detail_id = md.merchant_detail_id
WHERE
    md.deleted_at IS NULL
    AND m.deleted_at IS NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
GROUP BY
    md.merchant_detail_id,
    m.merchant_id
LIMIT $2
OFFSET
    $3;

-- GetMerchantDetailsTrashed: Retrieves merchant details for soft-deleted merchants
-- Parameters:
--   $1: search - Keyword to filter merchant name (case-insensitive, partial match)
--   $2: limit - Pagination limit
--   $3: offset - Pagination offset
-- Returns: Trashed merchant detail records matching the search, with total_count and social media links
-- name: GetMerchantDetailsTrashed :many
SELECT
    md.merchant_detail_id,
    md.merchant_id,
    md.display_name,
    md.cover_image_url,
    md.logo_url,
    md.short_description,
    md.website_url,
    md.created_at,
    md.updated_at,
    md.deleted_at,
    m.name AS merchant_name,
    COUNT(*) OVER () AS total_count,
    json_agg(
        json_build_object(
            'id',
            sml.merchant_social_id,
            'platform',
            sml.platform,
            'url',
            sml.url
        )
    ) AS social_media_links
FROM
    merchant_details md
    JOIN merchants m ON md.merchant_id = m.merchant_id
    LEFT JOIN merchant_social_media_links sml ON sml.merchant_detail_id = md.merchant_detail_id
WHERE
    md.deleted_at IS NOT NULL
    AND LOWER(m.name) LIKE LOWER(CONCAT('%', $1::text, '%'))
GROUP BY
    md.merchant_detail_id,
    m.merchant_id
LIMIT $2
OFFSET
    $3;

-- GetMerchantDetail: Retrieves a single merchant detail that is not soft-deleted
-- Parameters:
--   $1: merchant_detail_id - ID of the merchant detail
-- Business Logic:
--   - Returns the merchant detail where deleted_at IS NULL
-- Returns: A single merchant_details record with social media links
-- name: GetMerchantDetail :one
SELECT
    md.merchant_detail_id,
    md.merchant_id,
    md.display_name,
    md.cover_image_url,
    md.logo_url,
    md.short_description,
    md.website_url,
    md.created_at,
    md.updated_at,
    m.name AS merchant_name,
    json_agg(
        json_build_object(
            'id',
            sml.merchant_social_id,
            'platform',
            sml.platform,
            'url',
            sml.url
        )
    ) AS social_media_links
FROM
    merchant_details md
    JOIN merchants m ON md.merchant_id = m.merchant_id
    LEFT JOIN merchant_social_media_links sml ON sml.merchant_detail_id = md.merchant_detail_id
WHERE
    md.merchant_detail_id = $1
    AND md.deleted_at IS NULL
GROUP BY
    md.merchant_detail_id,
    m.merchant_id;

-- GetMerchantDetailTrashed: Fetches a single category by its ID
-- Purpose: Retrieve details of an active (non-deleted) category
-- Parameters:
--   $1: Category ID to search for
-- Returns:
--   Full category record if found and not deleted
-- Business Logic:
--   - Excludes soft-deleted categories
-- name: GetMerchantDetailTrashed :one
SELECT md.merchant_detail_id, md.merchant_id, md.display_name, md.cover_image_url, md.logo_url, md.short_description, md.website_url, md.created_at, md.updated_at
FROM merchant_details md
WHERE
    merchant_detail_id = $1
    AND deleted_at IS NOT NULL;

-- CreateMerchantDetail: Inserts a new merchant detail record
-- Purpose: Register profile details for a merchant
-- Parameters:
--   $1: merchant_id
--   $2: display_name
--   $3: cover_image_url
--   $4: logo_url
--   $5: short_description
--   $6: website_url
-- Returns: The newly created merchant detail record
-- name: CreateMerchantDetail :one
INSERT INTO
    merchant_details (
        merchant_id,
        display_name,
        cover_image_url,
        logo_url,
        short_description,
        website_url
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    merchant_detail_id,
    merchant_id,
    display_name,
    cover_image_url,
    logo_url,
    short_description,
    website_url,
    created_at,
    updated_at;

-- UpdateMerchantDetail: Updates an existing merchant detail record
-- Purpose: Modify merchant profile
-- Parameters:
--   $1: merchant_detail_id
--   $2: display_name
--   $3: cover_image_url
--   $4: logo_url
--   $5: short_description
--   $6: website_url
-- Returns: The updated merchant detail
-- name: UpdateMerchantDetail :one
UPDATE merchant_details
SET
    display_name = $2,
    cover_image_url = $3,
    logo_url = $4,
    short_description = $5,
    website_url = $6,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_detail_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_detail_id,
    merchant_id,
    display_name,
    cover_image_url,
    logo_url,
    short_description,
    website_url,
    created_at,
    updated_at;

-- TrashMerchantDetail: Soft-deletes a merchant detail
-- Purpose: Temporarily hide the merchant profile
-- Parameters:
--   $1: merchant_detail_id
-- Returns: The soft-deleted detail record
-- name: TrashMerchantDetail :one
UPDATE merchant_details
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    merchant_detail_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_detail_id,
    merchant_id,
    display_name,
    cover_image_url,
    logo_url,
    short_description,
    website_url,
    created_at,
    updated_at,
    deleted_at;

-- RestoreMerchantDetail: Restores a soft-deleted merchant detail
-- Purpose: Reactivate previously hidden profile
-- Parameters:
--   $1: merchant_detail_id
-- Returns: The restored record
-- name: RestoreMerchantDetail :one
UPDATE merchant_details
SET
    deleted_at = NULL
WHERE
    merchant_detail_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    merchant_detail_id,
    merchant_id,
    display_name,
    cover_image_url,
    logo_url,
    short_description,
    website_url,
    created_at,
    updated_at,
    deleted_at;

-- DeleteMerchantDetailPermanently: Hard delete a merchant detail
-- Parameters:
--   $1: merchant_detail_id
-- name: DeleteMerchantDetailPermanently :exec
DELETE FROM merchant_details
WHERE
    merchant_detail_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllMerchantDetails: Restores all soft-deleted merchant details
-- name: RestoreAllMerchantDetails :exec
UPDATE merchant_details
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentMerchantDetails: Permanently delete all soft-deleted records
-- name: DeleteAllPermanentMerchantDetails :exec
DELETE FROM merchant_details WHERE deleted_at IS NOT NULL;