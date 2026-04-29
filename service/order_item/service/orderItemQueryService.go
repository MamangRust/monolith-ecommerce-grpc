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

type orderItemQueryService struct {
	observability       observability.TraceLoggerObservability
	cache               cache.OrderItemQueryCache
	orderItemRepository repository.OrderItemQueryRepository
	logger              logger.LoggerInterface
}

type OrderItemQueryServiceDeps struct {
	Observability       observability.TraceLoggerObservability
	Cache               cache.OrderItemQueryCache
	OrderItemRepository repository.OrderItemQueryRepository
	Logger              logger.LoggerInterface
}

func NewOrderItemQueryService(deps *OrderItemQueryServiceDeps) OrderItemQueryService {
	return &orderItemQueryService{
		observability:       deps.Observability,
		cache:               deps.Cache,
		orderItemRepository: deps.OrderItemRepository,
		logger:              deps.Logger,
	}
}

func (s *orderItemQueryService) FindAll(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsRow, *int, error) {
	const method = "FindAll"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedOrderItemsAll(ctx, req); found {
		logSuccess("Successfully retrieved all order item records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, err := s.orderItemRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrderItemsRow](
			s.logger,
			orderitem_errors.ErrFailedFindAllOrderItems,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedOrderItemsAll(ctx, req, orderItems, &totalCount)

	logSuccess("Successfully fetched all order items",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderItems, &totalCount, nil
}

func (s *orderItemQueryService) FindActive(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsActiveRow, *int, error) {
	const method = "FindActive"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedOrderItemActive(ctx, req); found {
		logSuccess("Successfully retrieved active order item records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, err := s.orderItemRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrderItemsActiveRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemsByActive,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedOrderItemActive(ctx, req, orderItems, &totalCount)

	logSuccess("Successfully fetched active order items",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderItems, &totalCount, nil
}

func (s *orderItemQueryService) FindTrashed(ctx context.Context, req *requests.FindAllOrderItems) ([]*db.GetOrderItemsTrashedRow, *int, error) {
	const method = "FindTrashed"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)
	search := req.Search

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedOrderItemTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed order item records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orderItems, err := s.orderItemRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrderItemsTrashedRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemsByTrashed,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(orderItems) > 0 {
		totalCount = int(orderItems[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedOrderItemTrashed(ctx, req, orderItems, &totalCount)

	logSuccess("Successfully fetched trashed order items",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orderItems, &totalCount, nil
}

func (s *orderItemQueryService) FindByOrder(ctx context.Context, orderID int) ([]*db.GetOrderItemsByOrderRow, error) {
	const method = "FindByOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedOrderItems(ctx, orderID); found {
		logSuccess("Successfully retrieved order items by order ID from cache",
			zap.Int("order_id", orderID))
		return data, nil
	}

	orderItems, err := s.orderItemRepository.FindOrderItemByOrder(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[[]*db.GetOrderItemsByOrderRow](
			s.logger,
			orderitem_errors.ErrFailedFindOrderItemByOrder,
			method,
			span,

			zap.Int("order_id", orderID),
		)
	}

	s.cache.SetCachedOrderItems(ctx, orderItems)

	logSuccess("Successfully fetched order items by order ID",
		zap.Int("order_id", orderID))

	return orderItems, nil
}

func (s *orderItemQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
