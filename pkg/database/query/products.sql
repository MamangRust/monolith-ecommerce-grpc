-- GetProducts: Retrieves paginated list of active products with search capability
-- Purpose: List all active (non-deleted) products for display in UI
-- Parameters:
--   $1: search_term - Optional text to filter products (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All product fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Supports partial, case-insensitive search on name, description, brand, slug, and barcode
--   - Orders results by newest first (created_at DESC)
--   - Uses COUNT(*) OVER() to include total matching record count for pagination UI
-- name: GetProducts :many
SELECT
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM products as p
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR p.name ILIKE '%' || $1 || '%'
        OR p.description ILIKE '%' || $1 || '%'
        OR p.brand ILIKE '%' || $1 || '%'
        OR p.slug_product ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetProductsActive: Retrieves paginated list of active products (duplicate of GetProducts)
-- Purpose: Explicitly return active (non-deleted) products with search capability
-- Parameters:
--   $1: search_term - Optional text to filter products (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All product fields plus total_count of matching records
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Supports partial, case-insensitive search on name, description, brand, slug, and barcode
--   - Ordered by newest first (created_at DESC)
--   - Useful if frontend/backend wants clearer distinction in naming
-- name: GetProductsActive :many
SELECT
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM products as p
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR p.name ILIKE '%' || $1 || '%'
        OR p.description ILIKE '%' || $1 || '%'
        OR p.brand ILIKE '%' || $1 || '%'
        OR p.slug_product ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetProductsTrashed: Retrieves paginated list of trashed (soft-deleted) products
-- Purpose: List deleted products for admin to manage recovery or audit
-- Parameters:
--   $1: search_term - Optional text to filter trashed products (NULL for no filter)
--   $2: limit - Maximum number of records to return
--   $3: offset - Number of records to skip for pagination
-- Returns:
--   All product fields plus total_count of matching trashed records
-- Business Logic:
--   - Includes only soft-deleted products (deleted_at IS NOT NULL)
--   - Supports partial, case-insensitive search on name, description, brand, slug, and barcode
--   - Returns by newest first (created_at DESC)
--   - Used for "Trash Bin" UI or soft-delete management
-- name: GetProductsTrashed :many
SELECT
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM products as p
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR p.name ILIKE '%' || $1 || '%'
        OR p.description ILIKE '%' || $1 || '%'
        OR p.brand ILIKE '%' || $1 || '%'
        OR p.slug_product ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetProductsByMerchant: Retrieves paginated and filtered products owned by a specific merchant
-- Purpose: Allow merchants to view and manage their own products with advanced filtering options
-- Parameters:
--   $1: merchant_id - Filter products belonging to this merchant
--   $2: search_term - Optional text to filter by product name or description
--   $3: category_id - Optional category filter (0 or NULL to ignore)
--   $4: min_price - Minimum price filter (0 to ignore)
--   $5: max_price - Maximum price filter (0 to ignore, defaults to very high value)
--   $6: limit - Number of products to return (pagination)
--   $7: offset - Number of products to skip (pagination)
-- Returns:
--   - Filtered list of product fields including category name
--   - total_count of all matching products for pagination UI
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Supports case-insensitive partial search on name and description
--   - Filters by category ID only if provided
--   - Filters by price range only if values provided (>= min_price and <= max_price)
--   - Ordered by newest products first (created_at DESC)
-- name: GetProductsByMerchant :many
WITH
    filtered_products AS (
        SELECT
            p.product_id,
            p.merchant_id,
            p.category_id,
            p.weight,
            p.rating,
            p.slug_product,
            p.name,
            p.description,
            p.price,
            p.count_in_stock,
            p.brand,
            p.image_product,
            p.created_at,
            p.updated_at,
            c.name AS category_name
        FROM products p
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            p.deleted_at IS NULL
            AND p.merchant_id = $1
            AND (
                p.name ILIKE '%' || COALESCE($2, '') || '%'
                OR p.description ILIKE '%' || COALESCE($2, '') || '%'
                OR $2 IS NULL
            )
            AND (
                c.category_id = NULLIF($3, 0)
                OR NULLIF($3, 0) IS NULL
            )
            AND (
                p.price >= COALESCE(NULLIF($4, 0), 0)
                AND p.price <= COALESCE(NULLIF($5, 0), 999999999)
            )
    )
SELECT (
        SELECT COUNT(*)
        FROM filtered_products
    ) AS total_count, fp.*
FROM filtered_products fp
ORDER BY fp.created_at DESC
LIMIT $6
OFFSET
    $7;

-- GetProductsByCategoryName: Retrieves paginated and filtered products under a specific category name
-- Purpose: Display products by category for customers or category-focused pages
-- Parameters:
--   $1: category_name - The name of the category to filter by
--   $2: search_term - Optional text to filter by product name or description
--   $3: min_price - Minimum price filter (0 to ignore)
--   $4: max_price - Maximum price filter (0 to ignore, defaults to very high value)
--   $5: limit - Number of products to return (pagination)
--   $6: offset - Number of products to skip (pagination)
-- Returns:
--   - Filtered list of product fields including category name
--   - total_count of all matching products for pagination UI
-- Business Logic:
--   - Excludes soft-deleted products (deleted_at IS NULL)
--   - Matches category name exactly
--   - Supports case-insensitive partial search on name and description
--   - Filters by category ID only if provided
--   - Filters by price range only if values provided
--   - Ordered by newest products first (created_at DESC)
-- name: GetProductsByCategoryName :many
WITH
    filtered_products AS (
        SELECT
            p.product_id,
            p.merchant_id,
            p.category_id,
            p.weight,
            p.rating,
            p.slug_product,
            p.name,
            p.description,
            p.price,
            p.count_in_stock,
            p.brand,
            p.image_product,
            p.created_at,
            p.updated_at,
            c.name AS category_name
        FROM products p
            JOIN categories c ON p.category_id = c.category_id
        WHERE
            p.deleted_at IS NULL
            AND c.name = $1
            AND (
                $2::TEXT IS NULL
                OR p.name ILIKE '%' || $2::TEXT || '%'
                OR p.description ILIKE '%' || $2::TEXT || '%'
            )
            AND (
                (
                    $3::INTEGER IS NULL
                    OR p.price >= $3::INTEGER
                )
                AND (
                    $4::INTEGER IS NULL
                    OR p.price <= $4::INTEGER
                )
            )
    )
SELECT (
        SELECT COUNT(*)
        FROM filtered_products
    ) AS total_count, fp.*
FROM filtered_products fp
ORDER BY fp.created_at DESC
LIMIT $5
OFFSET
    $6;

-- CreateProduct: Creates a new product entry
-- Purpose: Add new products to merchant's catalog
-- Parameters:
--   $1: merchant_id - ID of the merchant owning the product
--   $2: category_id - Product category ID
--   $3: name - Product display name
--   $4: description - Detailed product description
--   $5: price - Product selling price
--   $6: count_in_stock - Available inventory quantity
--   $7: brand - Product brand name
--   $8: weight - Product weight in grams/kg
--   $9: rating - Initial product rating (0-5)
--   $10: slug_product - URL-friendly product identifier
--   $11: image_product - Main product image URL
-- Returns:
--   The complete created product record
-- Business Logic:
--   - Creates a new active product
--   - Requires all essential product information
--   - Returns full record for immediate use
-- name: CreateProduct :one
INSERT INTO
    products (
        merchant_id,
        category_id,
        name,
        description,
        price,
        count_in_stock,
        brand,
        weight,
        rating,
        slug_product,
        image_product
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11
    )
RETURNING
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at;

-- GetProductByID: Retrieves an active product by ID
-- Purpose: Display product details in storefront/merchant UI
-- Parameters:
--   $1: product_id - ID of the product to retrieve
-- Returns:
--   Complete product record if found and active
-- Business Logic:
--   - Only returns non-deleted (active) products
--   - Used for normal product display operations
-- name: GetProductByID :one
SELECT
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at
FROM products
WHERE
    product_id = $1
    AND deleted_at IS NULL;

-- GetProductByIdTrashed: Retrieves product including soft-deleted ones
-- Purpose: Access products in trash/recycle bin
-- Parameters:
--   $1: product_id - ID of the product to retrieve
-- Returns:
--   Complete product record regardless of deletion status
-- Business Logic:
--   - Bypasses soft-delete filter
--   - Used for admin/recovery operations
-- name: GetProductByIdTrashed :one
SELECT * FROM products WHERE product_id = $1;

-- UpdateProduct: Modifies all product details
-- Purpose: Edit product information in merchant dashboard
-- Parameters:
--   $1: product_id - ID of product to update
--   $2-$11: All product fields (category_id through image_product)
-- Returns:
--   The updated product record
-- Business Logic:
--   - Updates all mutable product fields
--   - Only works on active (non-deleted) products
--   - Automatically sets updated_at timestamp
-- name: UpdateProduct :one
UPDATE products
SET
    category_id = $2,
    name = $3,
    description = $4,
    price = $5,
    count_in_stock = $6,
    brand = $7,
    weight = $8,
    rating = $9,
    slug_product = $10,
    image_product = $11,
    updated_at = CURRENT_TIMESTAMP
WHERE
    product_id = $1
    AND deleted_at IS NULL
RETURNING
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at;

-- UpdateProductCountStock: Adjusts product inventory count
-- Purpose: Update stock levels after purchases/restocking
-- Parameters:
--   $1: product_id - ID of product to update
--   $2: count_in_stock - New inventory quantity
-- Returns:
--   The updated product record
-- Business Logic:
--   - Only modifies stock count
--   - Verifies product is active
--   - Used during order processing
-- name: UpdateProductCountStock :one
UPDATE products
SET
    count_in_stock = $2
WHERE
    product_id = $1
    AND deleted_at IS NULL
RETURNING
    product_id,
    count_in_stock;

-- TrashProduct: Soft-deletes a product
-- Purpose: Remove product from storefront while preserving data
-- Parameters:
--   $1: product_id - ID of product to trash
-- Returns:
--   The trashed product record
-- Business Logic:
--   - Sets deleted_at timestamp
--   - Only works on active products
--   - Allows recovery via RestoreProduct
-- name: TrashProduct :one
UPDATE products
SET
    deleted_at = current_timestamp
WHERE
    product_id = $1
    AND deleted_at IS NULL
RETURNING
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at;

-- RestoreProduct: Recovers a soft-deleted product
-- Purpose: Reactivate previously trashed products
-- Parameters:
--   $1: product_id - ID of product to restore
-- Returns:
--   The restored product record
-- Business Logic:
--   - Clears deleted_at timestamp
--   - Only works on trashed products
--   - Returns product to active status
-- name: RestoreProduct :one
UPDATE products
SET
    deleted_at = NULL
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    product_id,
    merchant_id,
    category_id,
    name,
    description,
    price,
    count_in_stock,
    brand,
    weight,
    rating,
    slug_product,
    image_product,
    created_at,
    updated_at;

-- DeleteProductPermanently: Removes a product from database
-- Purpose: Permanent deletion of trashed products
-- Parameters:
--   $1: product_id - ID of product to delete
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Physical deletion of record
--   - Only works on already-trashed products
--   - Irreversible operation
-- name: DeleteProductPermanently :exec
DELETE FROM products
WHERE
    product_id = $1
    AND deleted_at IS NOT NULL;

-- RestoreAllProducts: Recovers all soft-deleted products
-- Purpose: Bulk restore from trash/recycle bin
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Clears deleted_at for all trashed products
--   - Admin-level operation
--   - Returns all products to active status
-- name: RestoreAllProducts :exec
UPDATE products
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentProducts: Purges all trashed products
-- Purpose: Clean up recycle bin
-- Parameters: None
-- Returns:
--   Nothing (exec-only)
-- Business Logic:
--   - Permanent deletion of all trashed items
--   - Admin-level operation
--   - Irreversible bulk deletion
-- name: DeleteAllPermanentProducts :exec
DELETE FROM products WHERE deleted_at IS NOT NULL;