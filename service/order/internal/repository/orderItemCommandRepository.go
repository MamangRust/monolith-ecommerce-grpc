package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/record"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	recordmapper "github.com/MamangRust/monolith-ecommerce-shared/mapper/record"
)

type orderItemCommandRepository struct {
	db      *db.Queries
	mapping recordmapper.OrderItemRecordMapping
}

func NewOrderItemCommandRepository(db *db.Queries, mapping recordmapper.OrderItemRecordMapping) *orderItemCommandRepository {
	return &orderItemCommandRepository{
		db:      db,
		mapping: mapping,
	}
}

func (r *orderItemCommandRepository) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*record.OrderItemRecord, error) {
	res, err := r.db.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})

	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*record.OrderItemRecord, error) {
	res, err := r.db.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		OrderItemID: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})

	if err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) TrashedOrderItem(ctx context.Context, order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.TrashOrderItem(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) RestoreOrderItem(ctx context.Context, order_id int) (*record.OrderItemRecord, error) {
	res, err := r.db.RestoreOrderItem(ctx, int32(order_id))

	if err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}

	return r.mapping.ToOrderItemRecord(res), nil
}

func (r *orderItemCommandRepository) DeleteOrderItemPermanent(ctx context.Context, order_id int) (bool, error) {
	err := r.db.DeleteOrderItemPermanently(ctx, int32(order_id))

	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}

	return true, nil
}

func (r *orderItemCommandRepository) RestoreAllOrderItem(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllUsers(ctx)

	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteAllOrderPermanent(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrders(ctx)

	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent
	}

	return true, nil
}
