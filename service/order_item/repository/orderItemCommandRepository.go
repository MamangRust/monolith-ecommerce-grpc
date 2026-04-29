package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
)

type orderItemCommandRepository struct {
	db *db.Queries
}

func NewOrderItemCommandRepository(db *db.Queries) *orderItemCommandRepository {
	return &orderItemCommandRepository{
		db: db,
	}
}

func (r *orderItemCommandRepository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error) {
	return r.db.CreateOrderItem(ctx, db.CreateOrderItemParams{
		OrderID:   int32(req.OrderID),
		ProductID: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})
}

func (r *orderItemCommandRepository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error) {
	return r.db.UpdateOrderItem(ctx, db.UpdateOrderItemParams{
		OrderItemID: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})
}

func (r *orderItemCommandRepository) Trash(ctx context.Context, orderItemID int) (*db.OrderItem, error) {
	res, err := r.db.TrashOrderItem(ctx, int32(orderItemID))
	if err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}
	return res, nil
}

func (r *orderItemCommandRepository) Restore(ctx context.Context, orderItemID int) (*db.OrderItem, error) {
	res, err := r.db.RestoreOrderItem(ctx, int32(orderItemID))
	if err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}
	return res, nil
}

func (r *orderItemCommandRepository) DeletePermanent(ctx context.Context, orderItemID int) (bool, error) {
	err := r.db.DeleteOrderItemPermanently(ctx, int32(orderItemID))
	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteOrderItemByOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	err := r.db.DeleteOrderItemsByOrderPermanent(ctx, int32(orderID))
	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	err := r.db.RestoreAllOrdersItem(ctx)
	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem
	}
	return true, nil
}

func (r *orderItemCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	err := r.db.DeleteAllPermanentOrdersItem(ctx)
	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent
	}
	return true, nil
}

func (r *orderItemCommandRepository) CalculateTotalPrice(ctx context.Context, orderID int) (int, error) {
	res, err := r.db.CalculateTotalPrice(ctx, int32(orderID))
	if err != nil {
		return 0, orderitem_errors.ErrCalculateTotalPrice
	}
	return int(res), nil
}
