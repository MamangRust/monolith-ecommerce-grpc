-- CreateMerchantSocialMediaLink: Inserts a new merchant social media link
-- Purpose: Register a new social media link for a merchant
-- Parameters:
--   $1: merchant_detail_id
--   $2: platform
--   $3: url
-- Returns: The newly created social media link
-- Business Logic:
--   - Sets created_at, updated_at automatically
-- name: CreateMerchantSocialMediaLink :one
INSERT INTO
    merchant_social_media_links (
        merchant_detail_id,
        platform,
        url
    )
VALUES ($1, $2, $3)
RETURNING
    merchant_social_id,
    merchant_detail_id,
    platform,
    url,
    created_at,
    updated_at;

-- UpdateMerchantSocialMediaLink: Updates an existing merchant social media link
-- Purpose: Modify platform or URL
-- Parameters:
--   $1: merchant_social_id
--   $2: platform
--   $3: url
-- Returns: The updated social media link
-- Business Logic:
--   - Automatically updates updated_at timestamps
--   - Only affects non-deleted records
-- name: UpdateMerchantSocialMediaLink :one
UPDATE merchant_social_media_links
SET
    platform = $2,
    url = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    merchant_social_id = $1
RETURNING
    merchant_social_id,
    merchant_detail_id,
    platform,
    url,
    created_at,
    updated_at;

-- TrashMerchantSocialMediaLink: Soft-deletes a merchant social media link
-- Purpose: Temporarily deactivate a link without permanent deletion
-- Parameters:
--   $1: merchant_social_id
-- Returns: The soft-deleted social media link
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only affects active records
-- name: TrashMerchantSocialMediaLink :one
UPDATE merchant_social_media_links
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    merchant_social_id = $1
    AND deleted_at IS NULL
RETURNING
    merchant_social_id,
    merchant_detail_id,
    platform,
    url,
    created_at,
    updated_at,
    deleted_at;

-- RestoreMerchantSocialMediaLink: Restores a soft-deleted social media link
-- Purpose: Reactivate a previously deleted link
-- Parameters:
--   $1: merchant_social_id
-- Returns: The restored social media link
-- Business Logic:
--   - Clears deleted_at
--   - Only works on previously soft-deleted entries
-- name: RestoreMerchantSocialMediaLink :one
UPDATE merchant_social_media_links
SET
    deleted_at = NULL
WHERE
    merchant_social_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    merchant_social_id,
    merchant_detail_id,
    platform,
    url,
    created_at,
    updated_at,
    deleted_at;

-- DeleteMerchantSocialMediaLinkPermanently: Hard-deletes a social media link
-- Purpose: Permanently remove a link
-- Parameters:
--   $1: merchant_social_id
-- name: DeleteMerchantSocialMediaLinkPermanently :exec
DELETE FROM merchant_social_media_links
WHERE
    merchant_social_id = $1;

-- RestoreAllMerchantSocialMediaLinks: Restores all soft-deleted social media links
-- Purpose: Bulk recovery operation
-- Business Logic:
--   - Clears deleted_at on all trashed records
-- name: RestoreAllMerchantSocialMediaLinks :exec
UPDATE merchant_social_media_links
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllMerchantSocialMediaLinksPermanently: Hard-deletes all soft-deleted social media links
-- Purpose: Permanently remove all trashed links
-- Business Logic:
--   - Irreversible delete
--   - Only affects already soft-deleted records
-- name: DeleteAllMerchantSocialMediaLinksPermanently :exec
DELETE FROM merchant_social_media_links WHERE deleted_at IS NOT NULL;