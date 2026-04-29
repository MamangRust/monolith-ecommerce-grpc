package repository

import (
	"context"

	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	order_item_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type orderItemCommandRepository struct {
	client pb.OrderItemCommandServiceClient
}

func NewOrderItemCommandRepository(client pb.OrderItemCommandServiceClient) *orderItemCommandRepository {
	return &orderItemCommandRepository{
		client: client,
	}
}

func (r *orderItemCommandRepository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error) {
	res, err := r.client.CreateOrderItem(ctx, &pb.CreateOrderItemRecordRequest{
		OrderId:   int32(req.OrderID),
		ProductId: int32(req.ProductID),
		Quantity:  int32(req.Quantity),
		Price:     int32(req.Price),
	})
	if err != nil {
		return nil, order_item_errors.ErrCreateOrderItem.WithInternal(err)
	}

	return &db.CreateOrderItemRow{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error) {
	res, err := r.client.UpdateOrderItem(ctx, &pb.UpdateOrderItemRecordRequest{
		OrderItemId: int32(req.OrderItemID),
		Quantity:    int32(req.Quantity),
		Price:       int32(req.Price),
	})
	if err != nil {
		return nil, order_item_errors.ErrUpdateOrderItem.WithInternal(err)
	}

	return &db.UpdateOrderItemRow{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) Trash(ctx context.Context, order_id int) (*db.OrderItem, error) {
	res, err := r.client.TrashOrderItem(ctx, &pb.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return nil, order_item_errors.ErrTrashedOrderItem.WithInternal(err)
	}

	return &db.OrderItem{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) Restore(ctx context.Context, order_id int) (*db.OrderItem, error) {
	res, err := r.client.RestoreOrderItem(ctx, &pb.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return nil, order_item_errors.ErrRestoreOrderItem.WithInternal(err)
	}

	return &db.OrderItem{
		OrderItemID: res.Data.Id,
		OrderID:     res.Data.OrderId,
		ProductID:   res.Data.ProductId,
		Quantity:    res.Data.Quantity,
		Price:       res.Data.Price,
	}, nil
}

func (r *orderItemCommandRepository) DeletePermanent(ctx context.Context, order_id int) (bool, error) {
	res, err := r.client.DeleteOrderItemPermanent(ctx, &pb.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}

	return res.Status == "success", nil
}

func (r *orderItemCommandRepository) DeleteByOrderIDPermanent(ctx context.Context, order_id int) (bool, error) {
	res, err := r.client.DeleteOrderItemByOrderPermanent(ctx, &pb.FindByIdOrderItemRequest{Id: int32(order_id)})
	if err != nil {
		return false, order_item_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}

	return res.Status == "success", nil
}

func (r *orderItemCommandRepository) RestoreAll(ctx context.Context) (bool, error) {
	res, err := r.client.RestoreAllOrdersItem(ctx, &emptypb.Empty{})
	if err != nil {
		return false, order_item_errors.ErrRestoreAllOrderItem.WithInternal(err)
	}

	return res.Status == "success", nil
}

func (r *orderItemCommandRepository) DeleteAll(ctx context.Context) (bool, error) {
	res, err := r.client.DeleteAllPermanentOrdersItem(ctx, &emptypb.Empty{})
	if err != nil {
		return false, order_item_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}

	return res.Status == "success", nil
}
