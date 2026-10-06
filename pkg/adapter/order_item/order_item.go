// Package order_item adapts the Order Item service gRPC API into the shared
// pkg/database sqlc row types.
package order_item

import (
	"context"

	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
	"github.com/MamangRust/monolith-ecommerce-pkg/adapter"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/resilience"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"google.golang.org/protobuf/types/known/emptypb"
)

// QueryRepository is the contract consumers depend on for order-item reads.
type QueryRepository interface {
	FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.GetOrderItemsByOrderRow, error)
	// CalculateTotalPrice is served by the command service (it is the RPC the
	// owning service exposes), but callers treat it as a read.
	CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error)
}

// CommandRepository is the contract consumers depend on for order-item writes.
type CommandRepository interface {
	Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error)
	Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error)
	Trash(ctx context.Context, orderID int) (*db.OrderItem, error)
	Restore(ctx context.Context, orderID int) (*db.OrderItem, error)
	DeletePermanent(ctx context.Context, orderID int) (bool, error)
	DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error)
	RestoreAll(ctx context.Context) (bool, error)
	DeleteAll(ctx context.Context) (bool, error)
}

// Repository implements both QueryRepository and CommandRepository over the
// generated order-item query/command clients.
type Repository struct {
	query   pborder_item.OrderItemQueryServiceClient
	command pborder_item.OrderItemCommandServiceClient
	guard   *resilience.DependencyGuard
}

// New builds an order-item adapter. Passing zero options leaves the guard nil,
// which makes DependencyGuard.Call a plain passthrough.
func New(query pborder_item.OrderItemQueryServiceClient, command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) *Repository {
	r := &Repository{query: query, command: command}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// NewQueryAdapter builds a read-only adapter for consumers that never write.
func NewQueryAdapter(query pborder_item.OrderItemQueryServiceClient, opts ...adapter.GuardOption) QueryRepository {
	return New(query, nil, opts...)
}

// NewCommandAdapter builds a write-only adapter for consumers that never read.
func NewCommandAdapter(command pborder_item.OrderItemCommandServiceClient, opts ...adapter.GuardOption) CommandRepository {
	return New(nil, command, opts...)
}

func (r *Repository) SetGuard(g *resilience.DependencyGuard) {
	r.guard = g
}

func (r *Repository) FindOrderItemByOrder(ctx context.Context, orderID int) ([]*db.GetOrderItemsByOrderRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponsesOrderItem, error) {
		return r.query.FindOrderItemByOrder(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder.WithInternal(err)
	}
	if resp == nil {
		return nil, orderitem_errors.ErrFindOrderItemByOrder
	}

	items := make([]*db.GetOrderItemsByOrderRow, 0, len(resp.Data))
	for _, item := range resp.Data {
		items = append(items, &db.GetOrderItemsByOrderRow{
			OrderItemID: item.Id,
			OrderID:     item.OrderId,
			ProductID:   item.ProductId,
			Quantity:    item.Quantity,
			Price:       item.Price,
		})
	}

	return items, nil
}

func (r *Repository) CalculateTotalPrice(ctx context.Context, orderID int) (*int32, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.CalculateTotalPriceResponse, error) {
		return r.command.CalculateTotalPrice(callCtx, &pborder_item.CalculateTotalPriceRequest{OrderId: int32(orderID)})
	})
	if err != nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice.WithInternal(err)
	}
	if resp == nil {
		return nil, orderitem_errors.ErrCalculateTotalPrice
	}

	total := int32(resp.TotalPrice)
	return &total, nil
}

func (r *Repository) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.CreateOrderItem(callCtx, &pborder_item.CreateOrderItemRecordRequest{
			OrderId:   int32(req.OrderID),
			ProductId: int32(req.ProductID),
			Quantity:  int32(req.Quantity),
			Price:     int32(req.Price),
		})
	})
	if err != nil {
		return nil, orderitem_errors.ErrCreateOrderItem.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrCreateOrderItem
	}

	return &db.CreateOrderItemRow{
		OrderItemID: resp.Data.Id,
		OrderID:     resp.Data.OrderId,
		ProductID:   resp.Data.ProductId,
		Quantity:    resp.Data.Quantity,
		Price:       resp.Data.Price,
	}, nil
}

func (r *Repository) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.UpdateOrderItem(callCtx, &pborder_item.UpdateOrderItemRecordRequest{
			OrderItemId: int32(req.OrderItemID),
			Quantity:    int32(req.Quantity),
			Price:       int32(req.Price),
		})
	})
	if err != nil {
		return nil, orderitem_errors.ErrUpdateOrderItem.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrUpdateOrderItem
	}

	return &db.UpdateOrderItemRow{
		OrderItemID: resp.Data.Id,
		OrderID:     resp.Data.OrderId,
		ProductID:   resp.Data.ProductId,
		Quantity:    resp.Data.Quantity,
		Price:       resp.Data.Price,
	}, nil
}

func (r *Repository) Trash(ctx context.Context, orderID int) (*db.OrderItem, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.TrashOrderItem(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, orderitem_errors.ErrTrashedOrderItem.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrTrashedOrderItem
	}

	return &db.OrderItem{
		OrderItemID: resp.Data.Id,
		OrderID:     resp.Data.OrderId,
		ProductID:   resp.Data.ProductId,
		Quantity:    resp.Data.Quantity,
		Price:       resp.Data.Price,
	}, nil
}

func (r *Repository) Restore(ctx context.Context, orderID int) (*db.OrderItem, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItem, error) {
		return r.command.RestoreOrderItem(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return nil, orderitem_errors.ErrRestoreOrderItem.WithInternal(err)
	}
	if resp == nil || resp.Data == nil {
		return nil, orderitem_errors.ErrRestoreOrderItem
	}

	return &db.OrderItem{
		OrderItemID: resp.Data.Id,
		OrderID:     resp.Data.OrderId,
		ProductID:   resp.Data.ProductId,
		Quantity:    resp.Data.Quantity,
		Price:       resp.Data.Price,
	}, nil
}

func (r *Repository) DeletePermanent(ctx context.Context, orderID int) (bool, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemDelete, error) {
		return r.command.DeleteOrderItemPermanent(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	return resp != nil && resp.Status == "success", nil
}

func (r *Repository) DeleteByOrderIDPermanent(ctx context.Context, orderID int) (bool, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemDelete, error) {
		return r.command.DeleteOrderItemByOrderPermanent(callCtx, &pborder_item.FindByIdOrderItemRequest{Id: int32(orderID)})
	})
	if err != nil {
		return false, orderitem_errors.ErrDeleteOrderItemPermanent.WithInternal(err)
	}
	return resp != nil && resp.Status == "success", nil
}

func (r *Repository) RestoreAll(ctx context.Context) (bool, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemAll, error) {
		return r.command.RestoreAllOrdersItem(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, orderitem_errors.ErrRestoreAllOrderItem.WithInternal(err)
	}
	return resp != nil && resp.Status == "success", nil
}

func (r *Repository) DeleteAll(ctx context.Context) (bool, error) {
	resp, err := adapter.Call(r.guard, ctx, func(callCtx context.Context) (*pborder_item.ApiResponseOrderItemAll, error) {
		return r.command.DeleteAllPermanentOrdersItem(callCtx, &emptypb.Empty{})
	})
	if err != nil {
		return false, orderitem_errors.ErrDeleteAllOrderPermanent.WithInternal(err)
	}
	return resp != nil && resp.Status == "success", nil
}
