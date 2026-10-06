package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-order-item/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pborder_item "github.com/MamangRust/monolith-ecommerce-pb/order_item"
)

type orderItemCommandHandler struct {
	pborder_item.UnimplementedOrderItemCommandServiceServer
	orderItemService service.OrderItemCommandService
	logger           logger.LoggerInterface
}

func NewOrderItemCommandHandler(orderItemService service.OrderItemCommandService, logger logger.LoggerInterface) *orderItemCommandHandler {
	return &orderItemCommandHandler{
		orderItemService: orderItemService,
		logger:           logger,
	}
}

func (h *orderItemCommandHandler) CreateOrderItem(ctx context.Context, request *pborder_item.CreateOrderItemRecordRequest) (*pborder_item.ApiResponseOrderItem, error) {
	req := &requests.CreateOrderItemRecordRequest{
		OrderID:   int(request.GetOrderId()),
		ProductID: int(request.GetProductId()),
		Quantity:  int(request.GetQuantity()),
		Price:     int(request.GetPrice()),
	}

	if err := req.Validate(); err != nil {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItem, err := h.orderItemService.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully created order item",
		Data:    mapToProtoOrderItemResponse(orderItem),
	}, nil
}

func (h *orderItemCommandHandler) UpdateOrderItem(ctx context.Context, request *pborder_item.UpdateOrderItemRecordRequest) (*pborder_item.ApiResponseOrderItem, error) {
	id := int(request.GetOrderItemId())
	req := &requests.UpdateOrderItemRecordRequest{
		OrderItemID: id,
		Quantity:    int(request.GetQuantity()),
		Price:       int(request.GetPrice()),
	}

	if err := req.Validate(); err != nil {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItem, err := h.orderItemService.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully updated order item",
		Data:    mapToProtoOrderItemResponse(orderItem),
	}, nil
}

func (h *orderItemCommandHandler) TrashOrderItem(ctx context.Context, request *pborder_item.FindByIdOrderItemRequest) (*pborder_item.ApiResponseOrderItem, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItem, err := h.orderItemService.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully trashed order item",
		Data:    mapToProtoOrderItemResponse(orderItem),
	}, nil
}

func (h *orderItemCommandHandler) RestoreOrderItem(ctx context.Context, request *pborder_item.FindByIdOrderItemRequest) (*pborder_item.ApiResponseOrderItem, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	orderItem, err := h.orderItemService.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItem{
		Status:  "success",
		Message: "Successfully restored order item",
		Data:    mapToProtoOrderItemResponse(orderItem),
	}, nil
}

func (h *orderItemCommandHandler) DeleteOrderItemPermanent(ctx context.Context, request *pborder_item.FindByIdOrderItemRequest) (*pborder_item.ApiResponseOrderItemDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	_, err := h.orderItemService.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItemDelete{
		Status:  "success",
		Message: "Successfully deleted order item permanently",
	}, nil
}

func (h *orderItemCommandHandler) RestoreAllOrdersItem(ctx context.Context, _ *emptypb.Empty) (*pborder_item.ApiResponseOrderItemAll, error) {
	_, err := h.orderItemService.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItemAll{
		Status:  "success",
		Message: "Successfully restored all order items",
	}, nil
}

func (h *orderItemCommandHandler) DeleteAllPermanentOrdersItem(ctx context.Context, _ *emptypb.Empty) (*pborder_item.ApiResponseOrderItemAll, error) {
	_, err := h.orderItemService.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItemAll{
		Status:  "success",
		Message: "Successfully deleted all order items permanently",
	}, nil
}

func (h *orderItemCommandHandler) DeleteOrderItemByOrderPermanent(ctx context.Context, request *pborder_item.FindByIdOrderItemRequest) (*pborder_item.ApiResponseOrderItemDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	_, err := h.orderItemService.DeleteByOrderPermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.ApiResponseOrderItemDelete{
		Status:  "success",
		Message: "Successfully deleted order items by order permanently",
	}, nil
}

func (h *orderItemCommandHandler) CalculateTotalPrice(ctx context.Context, request *pborder_item.CalculateTotalPriceRequest) (*pborder_item.CalculateTotalPriceResponse, error) {
	id := int(request.GetOrderId())
	if id == 0 {
		return nil, orderitem_errors.ErrGrpcInvalidID
	}

	totalPrice, err := h.orderItemService.CalculateTotalPrice(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder_item.CalculateTotalPriceResponse{
		Status:     "success",
		Message:    "Successfully calculated total price",
		TotalPrice: int32(totalPrice),
	}, nil
}
