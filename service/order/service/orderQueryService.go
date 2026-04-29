package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-order/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-order/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type orderQueryService struct {
	observability   observability.TraceLoggerObservability
	cache           cache.OrderQueryCache
	orderRepository repository.OrderQueryRepository
	logger          logger.LoggerInterface
}

type OrderQueryServiceDeps struct {
	Observability   observability.TraceLoggerObservability
	Cache           cache.OrderQueryCache
	OrderRepository repository.OrderQueryRepository
	Logger          logger.LoggerInterface
}

func NewOrderQueryService(deps *OrderQueryServiceDeps) OrderQueryService {
	return &orderQueryService{
		observability:   deps.Observability,
		cache:           deps.Cache,
		orderRepository: deps.OrderRepository,
		logger:          deps.Logger,
	}
}

func (s *orderQueryService) FindAll(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersRow, *int, error) {
	const method = "FindAll"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetOrderAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all order records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderAllCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched all orders",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderQueryService) FindActive(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersActiveRow, *int, error) {
	const method = "FindActive"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetOrderActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active order records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersActiveRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderActiveCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched active orders",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderQueryService) FindTrashed(ctx context.Context, req *requests.FindAllOrder) ([]*db.GetOrdersTrashedRow, *int, error) {
	const method = "FindTrashed"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetOrderTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed order records from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersTrashedRow](
			s.logger,
			err,
			method,
			span,
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderTrashedCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched trashed orders",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderQueryService) FindByMerchant(ctx context.Context, req *requests.FindAllOrderByMerchant) ([]*db.GetOrdersByMerchantRow, *int, error) {
	const method = "FindByMerchant"

	page := req.Page
	pageSize := req.PageSize
	search := req.Search

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchant_id", req.MerchantID),
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetOrderByMerchantCache(ctx, req); found {
		logSuccess("Successfully retrieved orders by merchant from cache", zap.Int("totalRecords", *total), zap.Int("page", page), zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	orders, err := s.orderRepository.FindByMerchant(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetOrdersByMerchantRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("merchant_id", req.MerchantID),
			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(orders) > 0 {
		totalCount = int(orders[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetOrderByMerchantCache(ctx, req, orders, &totalCount)

	logSuccess("Successfully fetched orders by merchant",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return orders, &totalCount, nil
}

func (s *orderQueryService) FindByID(ctx context.Context, orderID int) (*db.GetOrderByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("orderID", orderID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedOrderCache(ctx, orderID); found {
		logSuccess("Successfully retrieved order by ID from cache", zap.Int("orderID", orderID))
		return data, nil
	}

	order, err := s.orderRepository.FindByID(ctx, orderID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetOrderByIDRow](
			s.logger,
			err,
			method,
			span,
			zap.Int("order_id", orderID),
		)
	}

	s.cache.SetCachedOrderCache(ctx, order)

	logSuccess("Successfully fetched order by ID", zap.Int("orderID", orderID))
	return order, nil
}
