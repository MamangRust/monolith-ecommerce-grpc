-- GetOrderStockReservationsByOrder: Retrieves reservations for an order.
-- name: GetOrderStockReservationsByOrder :many
SELECT reservation_id, order_id, product_id, quantity, status, created_at, updated_at
FROM order_stock_reservations
WHERE order_id = $1
ORDER BY product_id;

-- GetReservedStockReservationsForTrashedOrders: Finds active reservations belonging to trashed orders.
-- name: GetReservedStockReservationsForTrashedOrders :many
SELECT r.reservation_id, r.order_id, r.product_id, r.quantity, r.status, r.created_at, r.updated_at
FROM order_stock_reservations r
JOIN orders o ON o.order_id = r.order_id
WHERE o.deleted_at IS NOT NULL AND r.status = 'reserved'
ORDER BY r.order_id, r.product_id;

-- GetReleasedStockReservationsForTrashedOrders: Finds released reservations belonging to trashed orders.
-- name: GetReleasedStockReservationsForTrashedOrders :many
SELECT r.reservation_id, r.order_id, r.product_id, r.quantity, r.status, r.created_at, r.updated_at
FROM order_stock_reservations r
JOIN orders o ON o.order_id = r.order_id
WHERE o.deleted_at IS NOT NULL AND r.status = 'released'
ORDER BY r.order_id, r.product_id;

-- UpsertOrderStockReservation: Adds quantity to an order/product reservation.
-- name: UpsertOrderStockReservation :one
INSERT INTO order_stock_reservations (order_id, product_id, quantity, status)
VALUES ($1, $2, $3, 'reserved')
ON CONFLICT (order_id, product_id)
DO UPDATE SET
    quantity = order_stock_reservations.quantity + EXCLUDED.quantity,
    status = 'reserved',
    updated_at = CURRENT_TIMESTAMP
RETURNING reservation_id, order_id, product_id, quantity, status, created_at, updated_at;

-- UpdateOrderStockReservationQuantity: Sets the active reservation quantity.
-- name: UpdateOrderStockReservationQuantity :one
UPDATE order_stock_reservations
SET quantity = $3, status = 'reserved', updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1 AND product_id = $2
RETURNING reservation_id, order_id, product_id, quantity, status, created_at, updated_at;

-- MarkOrderStockReservationReleased: Marks one active reservation as released.
-- name: MarkOrderStockReservationReleased :one
UPDATE order_stock_reservations
SET status = 'released', updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1 AND product_id = $2 AND status = 'reserved'
RETURNING reservation_id, order_id, product_id, quantity, status, created_at, updated_at;

-- MarkOrderStockReservationReserved: Marks one released reservation as reserved.
-- name: MarkOrderStockReservationReserved :one
UPDATE order_stock_reservations
SET status = 'reserved', updated_at = CURRENT_TIMESTAMP
WHERE order_id = $1 AND product_id = $2 AND status = 'released'
RETURNING reservation_id, order_id, product_id, quantity, status, created_at, updated_at;

-- DeleteOrderStockReservation: Removes one reservation row for an order/product.
-- name: DeleteOrderStockReservation :exec
DELETE FROM order_stock_reservations WHERE order_id = $1 AND product_id = $2;

-- DeleteOrderStockReservations: Removes reservation rows for a permanently deleted order.
-- name: DeleteOrderStockReservations :exec
DELETE FROM order_stock_reservations WHERE order_id = $1;

-- DeleteAllOrderStockReservations: Removes reservations for permanently deleted orders.
-- name: DeleteAllOrderStockReservations :exec
DELETE FROM order_stock_reservations
WHERE order_id IN (SELECT order_id FROM orders WHERE deleted_at IS NOT NULL);

-- GetReleasedReservationsForActiveOrders: Finds reservations marked released that belong to
-- still-active orders. This is a drift pattern that durable reconciliation must repair by
-- re-reserving stock (an order should only release stock when it is trashed or completed).
-- name: GetReleasedReservationsForActiveOrders :many
SELECT r.reservation_id, r.order_id, r.product_id, r.quantity, r.status, r.created_at, r.updated_at
FROM order_stock_reservations r
JOIN orders o ON o.order_id = r.order_id
WHERE o.deleted_at IS NULL AND r.status = 'released'
ORDER BY r.order_id, r.product_id;

-- DeleteOldReleasedReservations: Removes released reservations older than a cutoff for
-- trashed orders. Retention keeps the idempotency ledger bounded while preserving active
-- reservations; only trashed orders may have their released rows cleaned up.
-- name: DeleteOldReleasedReservations :execrows
DELETE FROM order_stock_reservations r
USING orders o
WHERE r.order_id = o.order_id
  AND r.status = 'released'
  AND o.deleted_at IS NOT NULL
  AND r.updated_at < $1;

-- DeleteOldProductStockAdjustments: Purges idempotency ledger rows older than a cutoff.
-- Retention keeps the product_stock_adjustments table bounded without touching fresh rows
-- that may still be referenced by in-flight retries.
-- name: DeleteOldProductStockAdjustments :execrows
DELETE FROM product_stock_adjustments
WHERE created_at < $1;
