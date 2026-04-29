package service

import (
	"context"

	mencache "github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_business/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	merchantbusiness_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_business"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantBusinessQueryService struct {
	observability              observability.TraceLoggerObservability
	cache                      mencache.MerchantBusinessQueryCache
	merchantBusinessRepository repository.MerchantBusinessQueryRepository
	logger                     logger.LoggerInterface
}

type MerchantBusinessQueryServiceDeps struct {
	Observability              observability.TraceLoggerObservability
	Cache                      mencache.MerchantBusinessQueryCache
	MerchantBusinessRepository repository.MerchantBusinessQueryRepository
	Logger                     logger.LoggerInterface
}

func NewMerchantBusinessQueryService(deps *MerchantBusinessQueryServiceDeps) MerchantBusinessQueryService {
	return &merchantBusinessQueryService{
		observability:              deps.Observability,
		cache:                      deps.Cache,
		merchantBusinessRepository: deps.MerchantBusinessRepository,
		logger:                     deps.Logger,
	}
}

func (s *merchantBusinessQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantBusinessAll(ctx, req); found {
		logSuccess("Successfully retrieved all merchant business records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantBusinessRepository.FindAll(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindAllMerchantBusiness,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantBusinessAll(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched all merchant businesses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantBusinessQueryService) FindActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationActiveRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantBusinessActive(ctx, req); found {
		logSuccess("Successfully retrieved active merchant business records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantBusinessRepository.FindActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsBusinessInformationActiveRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindActiveMerchantBusiness,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantBusinessActive(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched active merchant businesses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantBusinessQueryService) FindTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantsBusinessInformationTrashedRow, *int, error) {
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

	if data, total, found := s.cache.GetCachedMerchantBusinessTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchant business records from cache",
			zap.Int("totalRecords", *total),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize))
		return data, total, nil
	}

	merchants, err := s.merchantBusinessRepository.FindTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantsBusinessInformationTrashedRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindTrashedMerchantBusiness,
			method,
			span,

			zap.String("search", search),
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int

	if len(merchants) > 0 {
		totalCount = int(merchants[0].TotalCount)
	} else {
		totalCount = 0
	}

	s.cache.SetCachedMerchantBusinessTrashed(ctx, req, merchants, &totalCount)

	logSuccess("Successfully fetched trashed merchant businesses",
		zap.Int("totalRecords", totalCount),
		zap.Int("page", page),
		zap.Int("pageSize", pageSize))

	return merchants, &totalCount, nil
}

func (s *merchantBusinessQueryService) FindByID(ctx context.Context, merchantID int) (*db.GetMerchantBusinessInformationRow, error) {
	const method = "FindByID"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantBusiness(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant business by ID from cache",
			zap.Int("merchantID", merchantID))
		return data, nil
	}

	merchant, err := s.merchantBusinessRepository.FindByID(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantBusinessInformationRow](
			s.logger,
			merchantbusiness_errors.ErrFailedFindMerchantBusinessById,
			method,
			span,

			zap.Int("merchant_id", merchantID),
		)
	}

	s.cache.SetCachedMerchantBusiness(ctx, merchant)

	logSuccess("Successfully fetched merchant business by ID",
		zap.Int("merchantID", merchantID))

	return merchant, nil
}
