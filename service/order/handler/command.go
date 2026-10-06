package handler

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-order/service"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errors"
	"github.com/MamangRust/monolith-ecommerce-shared/errors/order_errors"
	"google.golang.org/protobuf/types/known/emptypb"

	pborder "github.com/MamangRust/monolith-ecommerce-pb/order"
)

type orderCommandHandler struct {
	pborder.UnimplementedOrderCommandServiceServer
	orderCommand service.OrderCommandService
	logger       logger.LoggerInterface
}

func NewOrderCommandHandler(orderCommand service.OrderCommandService, logger logger.LoggerInterface) pborder.OrderCommandServiceServer {
	return &orderCommandHandler{
		orderCommand: orderCommand,
		logger:       logger,
	}
}

func (s *orderCommandHandler) Create(ctx context.Context, request *pborder.CreateOrderRequest) (*pborder.ApiResponseOrder, error) {
	var items []requests.CreateOrderItemRequest
	for _, item := range request.GetItems() {
		items = append(items, requests.CreateOrderItemRequest{
			ProductID: int(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
			Price:     int(item.GetPrice()),
		})
	}

	req := &requests.CreateOrderRequest{
		MerchantID: int(request.GetMerchantId()),
		UserID:     int(request.GetUserId()),
		TotalPrice: int(request.GetTotalPrice()),
		Items:      items,
		ShippingAddress: requests.CreateShippingAddressRequest{
			Alamat:         request.GetShipping().GetAlamat(),
			Provinsi:       request.GetShipping().GetProvinsi(),
			Kota:           request.GetShipping().GetKota(),
			Courier:        request.GetShipping().GetCourier(),
			ShippingMethod: request.GetShipping().GetShippingMethod(),
			ShippingCost:   int(request.GetShipping().GetShippingCost()),
			Negara:         request.GetShipping().GetNegara(),
		},
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateCreateOrder
	}

	order, err := s.orderCommand.Create(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully created order",
		Data:    mapToProtoOrderResponse(order),
	}, nil
}

func (s *orderCommandHandler) Update(ctx context.Context, request *pborder.UpdateOrderRequest) (*pborder.ApiResponseOrder, error) {
	var items []requests.UpdateOrderItemRequest
	for _, item := range request.GetItems() {
		items = append(items, requests.UpdateOrderItemRequest{
			OrderItemID: int(item.GetOrderItemId()),
			ProductID:   int(item.GetProductId()),
			Quantity:    int(item.GetQuantity()),
			Price:       int(item.GetPrice()),
		})
	}

	orderID := int(request.GetOrderId())

	var shippingID *int
	if request.GetShipping() != nil && request.GetShipping().GetShippingId() > 0 {
		id := int(request.GetShipping().GetShippingId())
		shippingID = &id
	}

	req := &requests.UpdateOrderRequest{
		OrderID:    &orderID,
		UserID:     int(request.GetUserId()),
		TotalPrice: int(request.GetTotalPrice()),
		Items:      items,
	}
	if shipping := request.GetShipping(); shipping != nil {
		req.ShippingAddress = &requests.UpdateShippingAddressRequest{
			ShippingID:     shippingID,
			Alamat:         shipping.GetAlamat(),
			Provinsi:       shipping.GetProvinsi(),
			Kota:           shipping.GetKota(),
			Courier:        shipping.GetCourier(),
			ShippingMethod: shipping.GetShippingMethod(),
			ShippingCost:   int(shipping.GetShippingCost()),
			Negara:         shipping.GetNegara(),
		}
	}

	if err := req.Validate(); err != nil {
		return nil, order_errors.ErrGrpcValidateUpdateOrder
	}

	order, err := s.orderCommand.Update(ctx, req)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrder{
		Status:  "success",
		Message: "Successfully updated order",
		Data:    mapToProtoOrderResponse(order),
	}, nil
}

func (s *orderCommandHandler) TrashedOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderCommand.Trash(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully trashed order",
		Data:    mapToProtoOrderResponseDeleteAt(order),
	}, nil
}

func (s *orderCommandHandler) RestoreOrder(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDeleteAt, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	order, err := s.orderCommand.Restore(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrderDeleteAt{
		Status:  "success",
		Message: "Successfully restored order",
		Data:    mapToProtoOrderResponseDeleteAt(order),
	}, nil
}

func (s *orderCommandHandler) DeleteOrderPermanent(ctx context.Context, request *pborder.FindByIdOrderRequest) (*pborder.ApiResponseOrderDelete, error) {
	id := int(request.GetId())
	if id == 0 {
		return nil, order_errors.ErrGrpcFailedInvalidId
	}

	_, err := s.orderCommand.DeletePermanent(ctx, id)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrderDelete{
		Status:  "success",
		Message: "Successfully deleted order permanently",
	}, nil
}

func (s *orderCommandHandler) RestoreAllOrder(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error) {
	_, err := s.orderCommand.RestoreAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully restored all orders",
	}, nil
}

func (s *orderCommandHandler) DeleteAllOrderPermanent(ctx context.Context, _ *emptypb.Empty) (*pborder.ApiResponseOrderAll, error) {
	_, err := s.orderCommand.DeleteAll(ctx)
	if err != nil {
		return nil, errors.ToGrpcError(err)
	}

	return &pborder.ApiResponseOrderAll{
		Status:  "success",
		Message: "Successfully deleted all orders permanently",
	}, nil
}
