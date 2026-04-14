package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	orderitem_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/order_item_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderItemCommandService struct {
	observability       observability.TraceLoggerObservability
	cache               cache.OrderItemCommandCache
	orderItemRepository repository.OrderItemCommandRepository
	logger              logger.LoggerInterface
}

type OrderItemCommandServiceDeps struct {
	Observability       observability.TraceLoggerObservability
	Cache               cache.OrderItemCommandCache
	OrderItemRepository repository.OrderItemCommandRepository
	Logger              logger.LoggerInterface
}

func NewOrderItemCommandService(deps *OrderItemCommandServiceDeps) *orderItemCommandService {
	return &orderItemCommandService{
		observability:       deps.Observability,
		cache:               deps.Cache,
		orderItemRepository: deps.OrderItemRepository,
		logger:              deps.Logger,
	}
}

func (s *orderItemCommandService) CreateOrderItem(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error) {
	const method = "CreateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", req.OrderID),
		attribute.Int("product_id", req.ProductID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.CreateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.CreateOrderItemRow](
			s.logger,
			orderitem_errors.ErrFailedCreateOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully created order item", zap.Int("order_item_id", int(res.OrderItemID)))

	return res, nil
}

func (s *orderItemCommandService) UpdateOrderItem(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error) {
	const method = "UpdateOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", req.OrderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.UpdateOrderItem(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.UpdateOrderItemRow](
			s.logger,
			orderitem_errors.ErrFailedUpdateOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully updated order item", zap.Int("order_item_id", req.OrderItemID))

	return res, nil
}

func (s *orderItemCommandService) TrashOrderItem(ctx context.Context, orderItemID int) (*db.OrderItem, error) {
	const method = "TrashOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", orderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.TrashOrderItem(ctx, orderItemID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedTrashedOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully trashed order item", zap.Int("order_item_id", orderItemID))

	return res, nil
}

func (s *orderItemCommandService) RestoreOrderItem(ctx context.Context, orderItemID int) (*db.OrderItem, error) {
	const method = "RestoreOrderItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", orderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.RestoreOrderItem(ctx, orderItemID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.OrderItem](
			s.logger,
			orderitem_errors.ErrFailedRestoreOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully restored order item", zap.Int("order_item_id", orderItemID))

	return res, nil
}

func (s *orderItemCommandService) DeleteOrderItemPermanent(ctx context.Context, orderItemID int) (bool, error) {
	const method = "DeleteOrderItemPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", orderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.DeleteOrderItemPermanent(ctx, orderItemID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully deleted permanent order item", zap.Int("order_item_id", orderItemID))

	return res, nil
}

func (s *orderItemCommandService) RestoreAllOrdersItem(ctx context.Context) (bool, error) {
	const method = "RestoreAllOrdersItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.RestoreAllOrdersItem(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedRestoreAllOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully restored all order items")

	return res, nil
}

func (s *orderItemCommandService) DeleteAllPermanentOrdersItem(ctx context.Context) (bool, error) {
	const method = "DeleteAllPermanentOrdersItem"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.DeleteAllPermanentOrdersItem(ctx)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[bool](
			s.logger,
			orderitem_errors.ErrFailedDeleteAllOrderItem,
			method,
			span,
		)
	}

	_ = s.cache.InvalidateOrderItemCache(ctx)

	logSuccess("Successfully deleted all permanent order items")

	return res, nil
}

func (s *orderItemCommandService) CalculateTotalPrice(ctx context.Context, orderID int) (int, error) {
	const method = "CalculateTotalPrice"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.CalculateTotalPrice(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[int](
			s.logger,
			orderitem_errors.ErrFailedCalculateTotal,
			method,
			span,
		)
	}

	logSuccess("Successfully calculated total price", zap.Int("order_id", orderID), zap.Int("total_price", res))

	return res, nil
}
