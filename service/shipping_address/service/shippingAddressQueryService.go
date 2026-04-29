package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-shipping-address/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	shipping_address_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/shipping_address_errors"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type shippingAddressQueryService struct {
	observability             observability.TraceLoggerObservability
	cache                     cache.ShippingAddressQueryCache
	shippingAddressRepository repository.ShippingAddressQueryRepository
	logger                    logger.LoggerInterface
}

type ShippingAddressQueryServiceDeps struct {
	Observability             observability.TraceLoggerObservability
	Cache                     cache.ShippingAddressQueryCache
	ShippingAddressRepository repository.ShippingAddressQueryRepository
	Logger                    logger.LoggerInterface
}

func NewShippingAddressQueryService(deps *ShippingAddressQueryServiceDeps) ShippingAddressQueryService {
	return &shippingAddressQueryService{
		observability:             deps.Observability,
		cache:                     deps.Cache,
		shippingAddressRepository: deps.ShippingAddressRepository,
		logger:                    deps.Logger,
	}
}

func (s *shippingAddressQueryService) FindAll(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressRow, *int, error) {
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

	if data, total, found := s.cache.GetShippingAddressAllCache(ctx, req); found {
		logSuccess("Successfully retrieved all shipping address records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	shippingAddresses, err := s.shippingAddressRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetShippingAddressRow](
			s.logger,
			shipping_address_errors.ErrFailedFindAllShippingAddresses,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(shippingAddresses) > 0 {
		totalCount = int(shippingAddresses[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetShippingAddressAllCache(ctx, req, shippingAddresses, &totalCount)

	logSuccess("Successfully fetched all shipping addresses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingAddresses, &totalCount, nil
}

func (s *shippingAddressQueryService) FindByID(ctx context.Context, shipping_id int) (*db.GetShippingByIDRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("shipping_id", shipping_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedShippingAddressCache(ctx, shipping_id); found {
		logSuccess("Successfully retrieved shipping address by ID from cache",
			zap.Int("shipping_id", shipping_id))
		return data, nil
	}

	shippingAddress, err := s.shippingAddressRepository.FindByID(ctx, shipping_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetShippingByIDRow](
			s.logger,
			err,
			method,
			span,

			zap.Int("shipping_id", shipping_id),
		)
	}

	s.cache.SetCachedShippingAddressCache(ctx, shippingAddress)

	logSuccess("Successfully fetched shipping address by ID",
		zap.Int("shipping_id", shipping_id))

	return shippingAddress, nil
}

func (s *shippingAddressQueryService) FindByOrder(ctx context.Context, order_id int) (*db.GetShippingAddressByOrderIDRow, error) {
	const method = "FindShippingAddressByOrder"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("order_id", order_id))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedShippingAddressByOrderCache(ctx, order_id); found {
		logSuccess("Successfully retrieved shipping address by order ID from cache",
			zap.Int("order_id", order_id))
		return data, nil
	}

	shippingAddress, err := s.shippingAddressRepository.FindByOrder(ctx, order_id)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetShippingAddressByOrderIDRow](
			s.logger,
			err,
			method,
			span,

			zap.Int("order_id", order_id),
		)
	}

	s.cache.SetCachedShippingAddressByOrderCache(ctx, shippingAddress)

	logSuccess("Successfully fetched shipping address by order ID",
		zap.Int("order_id", order_id))

	return shippingAddress, nil
}

func (s *shippingAddressQueryService) FindActive(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetShippingAddressActiveCache(ctx, req); found {
		logSuccess("Successfully retrieved active shipping address records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	shippingAddresses, err := s.shippingAddressRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetShippingAddressActiveRow](
			s.logger,
			shipping_address_errors.ErrFailedFindActiveShippingAddresses,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(shippingAddresses) > 0 {
		totalCount = int(shippingAddresses[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetShippingAddressActiveCache(ctx, req, shippingAddresses, &totalCount)

	logSuccess("Successfully fetched active shipping addresses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingAddresses, &totalCount, nil
}

func (s *shippingAddressQueryService) FindTrashed(ctx context.Context, req *requests.FindAllShippingAddress) ([]*db.GetShippingAddressTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetShippingAddressTrashedCache(ctx, req); found {
		logSuccess("Successfully retrieved trashed shipping address records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	shippingAddresses, err := s.shippingAddressRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetShippingAddressTrashedRow](
			s.logger,
			shipping_address_errors.ErrFailedFindTrashedShippingAddresses,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(shippingAddresses) > 0 {
		totalCount = int(shippingAddresses[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetShippingAddressTrashedCache(ctx, req, shippingAddresses, &totalCount)

	logSuccess("Successfully fetched trashed shipping addresses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return shippingAddresses, &totalCount, nil
}
