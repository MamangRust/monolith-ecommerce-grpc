package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order-item/repository"
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

func NewOrderItemCommandService(deps *OrderItemCommandServiceDeps) OrderItemCommandService {
	return &orderItemCommandService{
		observability:       deps.Observability,
		cache:               deps.Cache,
		orderItemRepository: deps.OrderItemRepository,
		logger:              deps.Logger,
	}
}

func (s *orderItemCommandService) Create(ctx context.Context, req *requests.CreateOrderItemRecordRequest) (*db.CreateOrderItemRow, error) {
	const method = "Create"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", req.OrderID),
		attribute.Int("product_id", req.ProductID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.Create(ctx, req)
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

func (s *orderItemCommandService) Update(ctx context.Context, req *requests.UpdateOrderItemRecordRequest) (*db.UpdateOrderItemRow, error) {
	const method = "Update"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", req.OrderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.Update(ctx, req)
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

func (s *orderItemCommandService) Trash(ctx context.Context, orderItemID int) (*db.OrderItem, error) {
	const method = "Trash"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", orderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.Trash(ctx, orderItemID)
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

func (s *orderItemCommandService) Restore(ctx context.Context, orderItemID int) (*db.OrderItem, error) {
	const method = "Restore"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", orderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.Restore(ctx, orderItemID)
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

func (s *orderItemCommandService) DeletePermanent(ctx context.Context, orderItemID int) (bool, error) {
	const method = "DeletePermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_item_id", orderItemID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.DeletePermanent(ctx, orderItemID)
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

func (s *orderItemCommandService) DeleteByOrderPermanent(ctx context.Context, orderID int) (bool, error) {
	const method = "DeleteByOrderPermanent"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.DeleteOrderItemByOrderPermanent(ctx, orderID)
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

	logSuccess("Successfully deleted permanent order items by order", zap.Int("order_id", orderID))

	return res, nil
}

func (s *orderItemCommandService) RestoreAll(ctx context.Context) (bool, error) {
	const method = "RestoreAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.RestoreAll(ctx)
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

func (s *orderItemCommandService) DeleteAll(ctx context.Context) (bool, error) {
	const method = "DeleteAll"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method)

	defer func() {
		end(status)
	}()

	res, err := s.orderItemRepository.DeleteAll(ctx)
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
