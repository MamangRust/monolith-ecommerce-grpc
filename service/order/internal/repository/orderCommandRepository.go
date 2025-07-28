package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderRecordMapping
}

func NewOrderCommandRepository(db *db.Queries, mapping recordmapper.OrderRecordMapping) *orderCommandRepository {
	return &orderCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderCommandRepository) CreateOrder(ctx context.Context, request *requests.CreateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.CreateOrderParams{
		MerchantID: int32(request.MerchantID),
		UserID:     int32(request.UserID),
		TotalPrice: int32(request.TotalPrice),
	}

	user, err := r.db.CreateOrder(ctx, req)

	if err != nil {
		return nil, order_errors.ErrCreateOrder
	}

	return r.mapping.ToOrderRecord(user), nil
}

func (r *orderCommandRepository) UpdateOrder(ctx context.Context, request *requests.UpdateOrderRecordRequest) (*record.OrderRecord, error) {
	req := db.UpdateOrderParams{
		OrderID:    int32(request.OrderID),
		TotalPrice: int32(request.TotalPrice),
	}

	res, err := r.db.UpdateOrder(ctx, req)

	if err != nil {
		return nil, order_errors.ErrUpdateOrder
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderCommandRepository) TrashedOrder(ctx context.Context, user_id int) (*record.OrderRecord, error) {
	res, err := r.db.TrashedOrder(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrTrashedOrder
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderCommandRepository) RestoreOrder(ctx context.Context, user_id int) (*record.OrderRecord, error) {
	res, err := r.db.RestoreOrder(ctx, int32(user_id))

	if err != nil {
		return nil, order_errors.ErrRestoreOrder
	}

	return r.mapping.ToOrderRecord(res), nil
}

func (r *orderCommandRepository) DeleteOrderPermanent(ctx context.Context, user_id int) (bool, error) {
	err := r.db.DeleteOrderPermanently(ctx, int32(user_id))

	if err != nil {
		return false, order_errors.ErrDeleteOrderPermanent
	}

	return true, nil
}

func (r *orderCommandRepository) RestoreAllOrder(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrders(ctx)

	if err != nil {
		return false, order_errors.ErrRestoreAllOrder
	}
	return true, nil
}

func (r *orderCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, order_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}
