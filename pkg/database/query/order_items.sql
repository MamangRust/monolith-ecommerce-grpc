-- GetOrderItems: Retrieves active order items with pagination and search
-- Purpose: Provides paginated listing of non-deleted order items for display or reporting
-- Parameters:
--   $1: Search keyword (matches order_id or product_id, optional)
--   $2: Limit (number of records per page)
--   $3: Offset (starting record index)
-- Returns:
--   All matching order item fields
--   total_count: Total number of results ignoring pagination (for frontend pagination UI)
-- Business Logic:
--   - Filters out soft-deleted items
--   - Supports keyword-based filtering
--   - Includes total result count via window function
-- name: GetOrderItems :many
SELECT
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at,
    COUNT(*) OVER () AS total_count
FROM order_items
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR product_id::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrderItemsActive: Retrieves active order items (duplicate-safe with GetOrderItems)
-- Purpose: Lists active (non-deleted) order items with pagination and optional search
-- Parameters:
--   $1: Search keyword (order_id or product_id, optional)
--   $2: Limit (pagination size)
--   $3: Offset (pagination start)
-- Returns:
--   Order item fields plus total matching count
-- Business Logic:
--   - Behaves similarly to GetOrderItems
--   - Used when clarity between active/trashed context is required
-- name: GetOrderItemsActive :many
SELECT
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM order_items
WHERE
    deleted_at IS NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR product_id::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- GetOrderItemsTrashed: Retrieves soft-deleted order items with pagination
-- Purpose: Allows review and management of trashed order items
-- Parameters:
--   $1: Search keyword (order_id or product_id, optional)
--   $2: Limit (number of rows per page)
--   $3: Offset (starting point for pagination)
-- Returns:
--   All matching deleted order item fields
--   total_count: Total number of trashed results
-- Business Logic:
--   - Only includes records with non-null deleted_at (trashed)
--   - Enables optional keyword search and pagination
--   - Sorted by deletion date for recent trash activity review
-- name: GetOrderItemsTrashed :many
-- name: GetOrderItemsTrashed :many
SELECT
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at,
    deleted_at,
    COUNT(*) OVER () AS total_count
FROM order_items
WHERE
    deleted_at IS NOT NULL
    AND (
        $1::TEXT IS NULL
        OR order_id::TEXT ILIKE '%' || $1 || '%'
        OR product_id::TEXT ILIKE '%' || $1 || '%'
    )
ORDER BY deleted_at DESC
LIMIT $2
OFFSET
    $3;

-- CalculateTotalPrice: Calculates total price of active order items for a specific order
-- Purpose: Provides the aggregated monetary value of an order
-- Parameters:
--   $1: order_id - identifier of the order
-- Returns:
--   total_price: Sum of (quantity * price) for all active items in the order
-- Business Logic:
--   - Ignores soft-deleted items
--   - Ensures result is zero if no items exist
-- name: CalculateTotalPrice :one
SELECT COALESCE(SUM(quantity * price), 0)::int AS total_price
FROM order_items
WHERE
    order_id = $1
    AND deleted_at IS NULL;

-- CreateOrderItem: Inserts a new order item record
-- Purpose: Adds a product to a specific order
-- Parameters:
--   $1: order_id
--   $2: product_id
--   $3: quantity
--   $4: price
-- Returns:
--   The newly created order item
-- Business Logic:
--   - Assumes quantity and price are validated in application layer
-- name: CreateOrderItem :one
INSERT INTO
    order_items (
        order_id,
        product_id,
        quantity,
        price
    )
VALUES ($1, $2, $3, $4)
RETURNING
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at;

-- GetOrderItemsByOrder: Retrieves active order items for a specific order
-- Purpose: Fetches all non-deleted order items under one order
-- Parameters:
--   $1: order_id
-- Returns:
--   List of active order items
-- Business Logic:
--   - Excludes soft-deleted entries
-- name: GetOrderItemsByOrder :many
SELECT
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at
FROM order_items
WHERE
    order_id = $1
    AND deleted_at IS NULL;

-- UpdateOrderItem: Updates quantity and price of an existing order item
-- Purpose: Allows modification of product details in an order
-- Parameters:
--   $1: order_item_id
--   $2: new quantity
--   $3: new price
-- Returns:
--   The updated order item
-- Business Logic:
--   - Applies changes only to active items
--   - Automatically updates `updated_at` timestamp
-- name: UpdateOrderItem :one
UPDATE order_items
SET
    quantity = $2,
    price = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    order_item_id = $1
    AND deleted_at IS NULL
RETURNING
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at;

-- TrashOrderItem: Soft-deletes a specific order item
-- Purpose: Marks an item as deleted without removing it from DB
-- Parameters:
--   $1: order_item_id
-- Returns:
--   The soft-deleted order item
-- Business Logic:
--   - Preserves record for potential restoration or audit
-- name: TrashOrderItem :one
UPDATE order_items
SET
    deleted_at = current_timestamp
WHERE
    order_item_id = $1
    AND deleted_at IS NULL
RETURNING
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at,
    deleted_at;

-- RestoreOrderItem: Restores a previously trashed order item
-- Purpose: Undoes a soft-delete action
-- Parameters:
--   $1: order_item_id
-- Returns:
--   The restored order item
-- Business Logic:
--   - Only restores items currently soft-deleted
-- name: RestoreOrderItem :one
UPDATE order_items
SET
    deleted_at = NULL
WHERE
    order_item_id = $1
    AND deleted_at IS NOT NULL
RETURNING
    order_item_id,
    order_id,
    product_id,
    quantity,
    price,
    created_at,
    updated_at,
    deleted_at;

-- DeleteOrderItemPermanently: Permanently deletes a trashed order item
-- Purpose: Removes the record entirely from the database
-- Parameters:
--   $1: order_item_id
-- Returns: None
-- Business Logic:
--   - Only deletes if already soft-deleted
--   - Irreversible action
-- name: DeleteOrderItemPermanently :exec
DELETE FROM order_items
WHERE
    order_item_id = $1
    AND deleted_at IS NOT NULL;

-- DeleteOrderItemsByOrderPermanent: Permanently deletes all trashed order items for a specific order
-- name: DeleteOrderItemsByOrderPermanent :exec
DELETE FROM order_items
WHERE
    order_id = $1;

-- RestoreAllOrdersItem: Restores all soft-deleted order items
-- Purpose: Mass recovery of trashed items
-- Parameters: None
-- Returns: None
-- Business Logic:
--   - Resets deleted_at to NULL for all trashed items
-- name: RestoreAllOrdersItem :exec
UPDATE order_items
SET
    deleted_at = NULL
WHERE
    deleted_at IS NOT NULL;

-- DeleteAllPermanentOrdersItem: Permanently deletes all trashed order items
-- Purpose: Performs hard delete of all soft-deleted items
-- Parameters: None
-- Returns: None
-- Business Logic:
--   - Used for data cleanup or archival enforcement
-- name: DeleteAllPermanentOrdersItem :exec
DELETE FROM order_items
WHERE
    deleted_at IS NOT NULL
    OR order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL);