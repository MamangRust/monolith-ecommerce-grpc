package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
)

type orderCommandRepository struct {
	db *db.Queries
}

func NewOrderCommandRepository(db *db.Queries) OrderCommandRepository {
	return &orderCommandRepository{
		db: db,
	}
}

func (r *orderCommandRepository) Create(ctx context.Context, request *requests.CreateOrderRecordRequest) (*db.CreateOrderRow, error) {
	req := db.CreateOrderParams{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.CreateOrder(ctx, req)

	if err != nil {
		return nil, order_errors.ErrCreateOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderCommandRepository) Update(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*db.UpdateOrderRow, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(ctx, req)

	if err != nil {
		return nil, order_errors.ErrUpdateOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderCommandRepository) Trash(ctx context.Context, order_id int) (*db.Order, error) {
	res, err := r.db.TrashedOrder(ctx, int32(order_id))

	if err != nil {
		return nil, order_errors.ErrTrashedOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderCommandRepository) Restore(ctx context.Context, order_id int) (*db.Order, error) {
	res, err := r.db.RestoreOrder(ctx, int32(order_id))

	if err != nil {
		return nil, order_errors.ErrRestoreOrder.WithInternal(err)
	}

	return res, nil
}

func (r *orderCommandRepository) DeletePermanent(ctx context.Context, order_id int) (bool, error) {
	err := r.db.DeleteOrderPermanently(ctx, int32(order_id))

	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent.WithInternal(err)
	}

	return true, nil
}

func (r *orderCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrders(ctx)

	if err != nil {
		return false, order_errors.ErrRestoreAllOrder.WithInternal(err)
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return true, nil
}

