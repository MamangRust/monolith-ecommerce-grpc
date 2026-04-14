package service

import (
	"context"

	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/cache"
	"github.com/MamangRust/monolith-ecommerce-grpc-merchant_detail/internal/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/errorhandler"
	merchantdetail_errors "github.com/MamangRust/monolith-ecommerce-shared/errors/merchant_detail"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type merchantDetailQueryService struct {
	observability    observability.TraceLoggerObservability
	cache            cache.MerchantDetailQueryCache
	merchantRepository repository.MerchantDetailQueryRepository
	logger           logger.LoggerInterface
}

type MerchantDetailQueryServiceDeps struct {
	Observability    observability.TraceLoggerObservability
	Cache            cache.MerchantDetailQueryCache
	Repository       repository.MerchantDetailQueryRepository
	Logger           logger.LoggerInterface
}

func NewMerchantDetailQueryService(deps *MerchantDetailQueryServiceDeps) *merchantDetailQueryService {
	return &merchantDetailQueryService{
		observability:    deps.Observability,
		cache:            deps.Cache,
		merchantRepository: deps.Repository,
		logger:           deps.Logger,
	}
}

func (s *merchantDetailQueryService) FindAll(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsRow, *int, error) {
	const method = "FindAllMerchants"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantDetailAll(ctx, req); found {
		logSuccess("Successfully retrieved merchants from cache", zap.Int("total", *total))
		return data, total, nil
	}

	res, err := s.merchantRepository.FindAllMerchants(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDetailsRow](
			s.logger,
			merchantdetail_errors.ErrFindAllMerchantDetails,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.cache.SetCachedMerchantDetailAll(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched merchants", zap.Int("total", totalCount))
	return res, &totalCount, nil
}

func (s *merchantDetailQueryService) FindByActive(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsActiveRow, *int, error) {
	const method = "FindByActiveMerchants"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantDetailActive(ctx, req); found {
		logSuccess("Successfully retrieved active merchants from cache", zap.Int("total", *total))
		return data, total, nil
	}

	res, err := s.merchantRepository.FindByActive(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDetailsActiveRow](
			s.logger,
			merchantdetail_errors.ErrFindActiveMerchantDetails,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.cache.SetCachedMerchantDetailActive(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched active merchants", zap.Int("total", totalCount))
	return res, &totalCount, nil
}

func (s *merchantDetailQueryService) FindByTrashed(ctx context.Context, req *requests.FindAllMerchant) ([]*db.GetMerchantDetailsTrashedRow, *int, error) {
	const method = "FindByTrashedMerchants"

	page, pageSize := s.normalizePagination(req.Page, req.PageSize)

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("page", page),
		attribute.Int("pageSize", pageSize),
		attribute.String("search", req.Search))

	defer func() {
		end(status)
	}()

	if data, total, found := s.cache.GetCachedMerchantDetailTrashed(ctx, req); found {
		logSuccess("Successfully retrieved trashed merchants from cache", zap.Int("total", *total))
		return data, total, nil
	}

	res, err := s.merchantRepository.FindByTrashed(ctx, req)
	if err != nil {
		status = "error"
		return errorhandler.HandlerErrorPagination[[]*db.GetMerchantDetailsTrashedRow](
			s.logger,
			merchantdetail_errors.ErrFindTrashedMerchantDetails,
			method,
			span,
			zap.Int("page", page),
			zap.Int("pageSize", pageSize),
		)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	}

	s.cache.SetCachedMerchantDetailTrashed(ctx, req, res, &totalCount)

	logSuccess("Successfully fetched trashed merchants", zap.Int("total", totalCount))
	return res, &totalCount, nil
}

func (s *merchantDetailQueryService) FindById(ctx context.Context, merchantID int) (*db.GetMerchantDetailRow, error) {
	const method = "FindByIdMerchantDetail"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method,
		attribute.Int("merchantID", merchantID))

	defer func() {
		end(status)
	}()

	if data, found := s.cache.GetCachedMerchantDetail(ctx, merchantID); found {
		logSuccess("Successfully retrieved merchant detail from cache", zap.Int("merchantID", merchantID))
		return data, nil
	}

	res, err := s.merchantRepository.FindById(ctx, merchantID)
	if err != nil {
		status = "error"
		return errorhandler.HandleError[*db.GetMerchantDetailRow](
			s.logger,
			merchantdetail_errors.ErrMerchantDetailNotFound,
			method,
			span,
			zap.Int("merchantID", merchantID),
		)
	}

	s.cache.SetCachedMerchantDetail(ctx, res)

	logSuccess("Successfully fetched merchant detail", zap.Int("merchantID", merchantID))
	return res, nil
}

func (s *merchantDetailQueryService) normalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return page, pageSize
}
