package repository

import (
	"context"
	"time"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
)

type StockReservationRepository interface {
	GetByOrder(ctx context.Context, orderID int) ([]*db.OrderStockReservation, error)
	Upsert(ctx context.Context, orderID, productID, quantity int) (*db.OrderStockReservation, error)
	UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*db.OrderStockReservation, error)
	Release(ctx context.Context, orderID, productID int) (*db.OrderStockReservation, error)
	Reserve(ctx context.Context, orderID, productID int) (*db.OrderStockReservation, error)
	GetReservedForTrashedOrders(ctx context.Context) ([]*db.OrderStockReservation, error)
	GetReleasedForTrashedOrders(ctx context.Context) ([]*db.OrderStockReservation, error)
	DeleteByOrder(ctx context.Context, orderID int) error
	DeleteByOrderProduct(ctx context.Context, orderID, productID int) error
	DeleteAllForTrashedOrders(ctx context.Context) error

	// GetReleasedForActiveOrders returns reservations marked released while their
	// order is still active — a drift pattern repaired by durable reconciliation.
	GetReleasedForActiveOrders(ctx context.Context) ([]*db.OrderStockReservation, error)

	// DeleteOldReleasedReservations removes released reservations for trashed orders
	// whose updated_at predates the cutoff. It returns the number of rows removed.
	DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error)

	// DeleteOldProductStockAdjustments purges idempotency ledger rows older than the
	// cutoff. It returns the number of rows removed.
	DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error)
}

type stockReservationRepository struct {
	db *db.Queries
}

func NewStockReservationRepository(queries *db.Queries) StockReservationRepository {
	return &stockReservationRepository{db: queries}
}

func (r *stockReservationRepository) GetByOrder(ctx context.Context, orderID int) ([]*db.OrderStockReservation, error) {
	return r.db.GetOrderStockReservationsByOrder(ctx, int32(orderID))
}

func (r *stockReservationRepository) Upsert(ctx context.Context, orderID, productID, quantity int) (*db.OrderStockReservation, error) {
	return r.db.UpsertOrderStockReservation(ctx, db.UpsertOrderStockReservationParams{
		OrderID: int32(orderID), ProductID: int32(productID), Quantity: int32(quantity),
	})
}

func (r *stockReservationRepository) UpdateQuantity(ctx context.Context, orderID, productID, quantity int) (*db.OrderStockReservation, error) {
	return r.db.UpdateOrderStockReservationQuantity(ctx, db.UpdateOrderStockReservationQuantityParams{
		OrderID: int32(orderID), ProductID: int32(productID), Quantity: int32(quantity),
	})
}

func (r *stockReservationRepository) Release(ctx context.Context, orderID, productID int) (*db.OrderStockReservation, error) {
	return r.db.MarkOrderStockReservationReleased(ctx, db.MarkOrderStockReservationReleasedParams{
		OrderID: int32(orderID), ProductID: int32(productID),
	})
}

func (r *stockReservationRepository) Reserve(ctx context.Context, orderID, productID int) (*db.OrderStockReservation, error) {
	return r.db.MarkOrderStockReservationReserved(ctx, db.MarkOrderStockReservationReservedParams{
		OrderID: int32(orderID), ProductID: int32(productID),
	})
}

func (r *stockReservationRepository) GetReservedForTrashedOrders(ctx context.Context) ([]*db.OrderStockReservation, error) {
	return r.db.GetReservedStockReservationsForTrashedOrders(ctx)
}

func (r *stockReservationRepository) GetReleasedForTrashedOrders(ctx context.Context) ([]*db.OrderStockReservation, error) {
	return r.db.GetReleasedStockReservationsForTrashedOrders(ctx)
}

func (r *stockReservationRepository) DeleteByOrder(ctx context.Context, orderID int) error {
	return r.db.DeleteOrderStockReservations(ctx, int32(orderID))
}

func (r *stockReservationRepository) DeleteByOrderProduct(ctx context.Context, orderID, productID int) error {
	return r.db.DeleteOrderStockReservation(ctx, db.DeleteOrderStockReservationParams{
		OrderID: int32(orderID), ProductID: int32(productID),
	})
}

func (r *stockReservationRepository) DeleteAllForTrashedOrders(ctx context.Context) error {
	return r.db.DeleteAllOrderStockReservations(ctx)
}

func (r *stockReservationRepository) GetReleasedForActiveOrders(ctx context.Context) ([]*db.OrderStockReservation, error) {
	return r.db.GetReleasedReservationsForActiveOrders(ctx)
}

func (r *stockReservationRepository) DeleteOldReleasedReservations(ctx context.Context, cutoff time.Time) (int64, error) {
	return r.db.DeleteOldReleasedReservations(ctx, pgtype.Timestamp{Time: cutoff, Valid: true})
}

func (r *stockReservationRepository) DeleteOldProductStockAdjustments(ctx context.Context, cutoff time.Time) (int64, error) {
	return r.db.DeleteOldProductStockAdjustments(ctx, pgtype.Timestamp{Time: cutoff, Valid: true})
}
